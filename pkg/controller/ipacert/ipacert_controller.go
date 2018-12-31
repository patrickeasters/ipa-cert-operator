package ipacert

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-cmp/cmp"
	certv1alpha1 "github.com/patrickeasters/ipa-cert-operator/pkg/apis/cert/v1alpha1"
	"github.com/patrickeasters/ipa-cert-operator/pkg/ipa"
	"github.com/patrickeasters/ipa-cert-operator/pkg/settings"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_ipacert")

// Add creates a new IpaCert Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileIpaCert{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("ipacert-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IpaCert
	err = c.Watch(&source.Kind{Type: &certv1alpha1.IpaCert{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary secret resources and requeue the owner IpaCert
	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &certv1alpha1.IpaCert{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileIpaCert{}

// ReconcileIpaCert reconciles a IpaCert object
type ReconcileIpaCert struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a IpaCert object and makes changes based on the state read
// and what is in the IpaCert.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIpaCert) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("IpaCert.Namespace", request.Namespace, "IpaCert.Name", request.Name)
	reqLogger.Info("Reconciling IpaCert")

	// Fetch the IpaCert instance
	instance := &certv1alpha1.IpaCert{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Check if a secret already exists and create if needed
	found := &corev1.Secret{}
	secretName := instance.Name + "-tls"
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: secretName, Namespace: instance.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Secret", "Secret.Namespace", instance.Namespace, "Secret.Name", secretName)

		// Define a new Secret object
		secret, err := newSecretForCR(secretName, instance)
		if err != nil {
			// Update status on CR and return an error
			return r.returnError(instance, err)
		}

		// Set IpaCert instance as the owner and controller
		if err := controllerutil.SetControllerReference(instance, secret, r.scheme); err != nil {
			return reconcile.Result{}, err
		}

		// Finally, create the secret
		err = r.client.Create(context.TODO(), secret)
		if err != nil {
			// Update status on CR and return an error
			return r.returnError(instance, err)
		}
		// Secret created successfully - return and requeue
		return reconcile.Result{Requeue: true}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Get current status info from cert secret
	secretBytes, ok := found.Data["tls.crt"]
	if !ok {
		// This secret appears to be malformed, so let's delete it and start fresh
		r.client.Delete(context.TODO(), found)
		return reconcile.Result{}, fmt.Errorf("Secret %s deleted due to missing tls.crt data", found.Name)
	}
	secretPem := string(secretBytes)
	secretStatus, err := ipa.CertStatus(secretPem)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("Secret %s deleted due to unreadable cert: %s", found.Name, err)
	}

	// Let's check and see if the cert is outdated or expiring soon
	switch {
	case instance.Spec.Cn != secretStatus.Cn:
		reqLogger.Info("Re-issuing cert with updated CN")
		return r.renewCert(instance, found)
	case !containsNames(instance.Spec.AdditionalNames, secretStatus.DnsNames):
		reqLogger.Info("Re-issuing cert with updated SANs")
		return r.renewCert(instance, found)
	case time.Now().After(secretStatus.Expiry.Add(-1 * settings.Instance.RenewalPeriod)):
		reqLogger.Info("Renewing cert because it expires soon.")
		return r.renewCert(instance, found)
	}

	// If nothing else, let's update status and check back later
	newStatus := certv1alpha1.IpaCertStatus{
		Status:   "ok",
		CertData: *secretStatus,
	}
	if !cmp.Equal(newStatus, instance.Status) {
		reqLogger.Info("Updating IpaCert status to match reality")
		instance.Status = newStatus
		r.client.Status().Update(context.TODO(), instance)
	}
	reqLogger.Info("Cert looks good. Nothing else to do here")
	return reconcile.Result{RequeueAfter: settings.Instance.RequeuePeriod}, nil
}

// newSecretForCR requests a cert and generates a secret based on it
func newSecretForCR(name string, cr *certv1alpha1.IpaCert) (*corev1.Secret, error) {
	cert, key, err := issueCert(cr)
	if err != nil {
		return nil, err
	}

	// Make our secret object
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-tls",
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"ipa-cert": cr.Name,
			},
		},
		Data: map[string][]byte{
			"tls.crt": []byte(cert),
			"tls.key": []byte(key),
		},
	}
	return secret, nil
}

// Get new cert from IPA
func issueCert(cr *certv1alpha1.IpaCert) (string, string, error) {
	// We support either a host or user principal
	principalType := cr.Spec.PrincipalType
	if principalType != "host" && principalType != "user" {
		// We can probably guess this pretty easily... Worst case we error out later
		if strings.Contains(cr.Spec.Cn, ".") {
			principalType = "host"
		} else {
			principalType = "user"
		}
	}

	// Generate a CSR and request a cert from IPA
	sans := append([]string{cr.Spec.Cn}, cr.Spec.AdditionalNames...)
	csr, key := ipa.GenerateCsr(cr.Spec.Cn, sans)
	principal := principalType + "/" + cr.Spec.Cn
	cert, err := ipa.RequestCert(principalType, principal, csr)
	if err != nil {
		return "", "", err
	}

	// Add the CA chain if present unless the CR tells us otherwise
	if len(settings.Instance.CaChain) > 0 && !cr.Spec.ExcludeChain {
		cert += settings.Instance.CaChain
	}

	return cert, key, nil
}

// Renew existing certificate and end reconciliation
func (r *ReconcileIpaCert) renewCert(cr *certv1alpha1.IpaCert, secret *corev1.Secret) (reconcile.Result, error) {
	cert, key, err := issueCert(cr)
	if err != nil {
		// Update status on CR and return an error
		return r.returnError(cr, err)
	}

	// Update existing secret
	secret.Data["tls.crt"] = []byte(cert)
	secret.Data["tls.key"] = []byte(key)
	err = r.client.Update(context.TODO(), secret)
	if err != nil {
		// Update status on CR and return an error
		return r.returnError(cr, err)
	}

	// Secret updated successfully
	return reconcile.Result{}, nil
}

// Return an error and update CR status with details
func (r *ReconcileIpaCert) returnError(cr *certv1alpha1.IpaCert, err error) (reconcile.Result, error) {
	newStatus := certv1alpha1.IpaCertStatus{
		Status:       "error",
		StatusReason: err.Error(),
		CertData:     cr.Status.CertData,
	}
	// Only update status if it's different to avoid a reconciliation loop
	if !cmp.Equal(newStatus, cr) {
		cr.Status = newStatus
		r.client.Status().Update(context.TODO(), cr)
	}
	return reconcile.Result{}, err
}

// Ensures all additional SANs are present in SAN slice
func containsNames(sub, super []string) bool {
	for _, i := range sub {
		for _, j := range super {
			if i == j {
				break
			}
		}
		return false
	}
	return true
}
