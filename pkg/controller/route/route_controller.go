package route

import (
	"context"
	"encoding/pem"
	"fmt"
	"time"

	routev1 "github.com/openshift/api/route/v1"
	"github.com/patrickeasters/ipa-cert-operator/pkg/ipa"
	"github.com/patrickeasters/ipa-cert-operator/pkg/settings"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_route")

const ROUTE_IPA_MANAGED = "cert.patrickeasters.com/ipa-managed"
const ROUTE_IPA_STATUS = "cert.patrickeasters.com/ipa-status"
const ROUTE_IPA_SERIAL = "cert.patrickeasters.com/ipa-serial"

// Add creates a new Route Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileRoute{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("route-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Route
	err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileRoute{}

// ReconcileRoute reconciles a Route object
type ReconcileRoute struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Route object and makes changes based on the state read
// and what is in the Route.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileRoute) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Route.Namespace", request.Namespace, "Route.Name", request.Name)
	reqLogger.Info("Reconciling Route")

	// Fetch the Route instance
	instance := &routev1.Route{}
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

	if instance.ObjectMeta.Annotations[ROUTE_IPA_MANAGED] != "true" {
		// Route is not associated with an IPA cert, so we don't need to do anything else
		reqLogger.Info("Skip reconcile: this is not managed by us")
		return reconcile.Result{}, nil
	}

	if instance.Spec.TLS.Termination == routev1.TLSTerminationPassthrough {
		err := fmt.Errorf("Route certs are not configurable for passthrough routes.")
		instance.ObjectMeta.Annotations[ROUTE_IPA_STATUS] = fmt.Sprintf("Error: %s", err)
		r.client.Status().Update(context.TODO(), instance)
		// Not returning erorr since this is a user config issue that won't be fixed by a retry
		return reconcile.Result{}, nil
	}

	// Check status of route cert
	routeCert, err := ipa.CertStatus(instance.Spec.TLS.Certificate)
	keyBlock, _ := pem.Decode([]byte(instance.Spec.TLS.Key))

	// Let's renew/re-issue the cert if needed
	switch {
	case err != nil:
		reqLogger.Info("Unable to parse route certificate. Proceeding with renewal.")
		return r.renewCert(instance)
	case keyBlock == nil || keyBlock.Type != "PRIVATE KEY":
		reqLogger.Info("Route is missing private key. Proceeding with renewal.")
		return r.renewCert(instance)
	case time.Now().After(routeCert.Expiry.Add(-1 * settings.Instance.RenewalPeriod)):
		reqLogger.Info("Renewing cert because it expires soon.")
		return r.renewCert(instance)
	}

	// All looks good. Check back later
	reqLogger.Info("Skip reconcile: everything looks good")
	instance.ObjectMeta.Annotations[ROUTE_IPA_STATUS] = "ok"
	r.client.Status().Update(context.TODO(), instance)
	return reconcile.Result{RequeueAfter: settings.Instance.RequeuePeriod}, nil
}

// Renew existing certificate and end reconciliation
func (r *ReconcileRoute) renewCert(route *routev1.Route) (reconcile.Result, error) {
	// Generate a CSR and request a cert from IPA
	csr, key := ipa.GenerateCsr(route.Spec.Host, []string{route.Spec.Host})
	principal := "host/" + route.Spec.Host
	cert, err := ipa.RequestCert("host", principal, csr)
	if err != nil {
		return r.returnError(route, err)
	}

	// Update route
	route.Spec.TLS.Certificate = cert
	route.Spec.TLS.Key = key
	route.Spec.TLS.CACertificate = settings.Instance.CaChain
	err = r.client.Update(context.TODO(), route)
	if err != nil {
		// Update status on CR and return an error
		return r.returnError(route, err)
	}

	// Secret updated successfully
	return reconcile.Result{}, nil
}

// Return an error and update CR status with details
func (r *ReconcileRoute) returnError(route *routev1.Route, err error) (reconcile.Result, error) {
	// Only update status if it's different to avoid a reconciliation loop
	newStatus := fmt.Sprintf("Error: %s", err)
	if route.Annotations[ROUTE_IPA_STATUS] != newStatus {
		route.Annotations[ROUTE_IPA_STATUS] = newStatus
		r.client.Status().Update(context.TODO(), route)
	}
	return reconcile.Result{}, err
}
