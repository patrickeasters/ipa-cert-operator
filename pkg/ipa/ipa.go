package ipa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"net/http"
	"time"

	ipa "github.com/patrickeasters/goipa"
	certv1alpha1 "github.com/patrickeasters/ipa-cert-operator/pkg/apis/cert/v1alpha1"
	"github.com/patrickeasters/ipa-cert-operator/pkg/settings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CertData struct {
	Serial         string
	Cn             string
	AlternateNames []string
	Issued         time.Time
	Expiry         time.Time
}

func GenerateCsr(cn string, sans []string) (string, string) {
	keyBytes, _ := rsa.GenerateKey(rand.Reader, 4096)
	pkey := string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(keyBytes)}))

	subj := pkix.Name{
		CommonName: cn,
	}

	template := x509.CertificateRequest{
		Subject:            subj,
		DNSNames:           sans,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csrBytes, _ := x509.CreateCertificateRequest(rand.Reader, &template, keyBytes)
	csr := string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrBytes}))

	return csr, pkey
}

func RequestCert(principalType, principal, csr string) (string, error) {
	httpClient := &http.Client{
		Timeout: settings.Instance.IpaTimeout,
	}
	client := ipa.NewClientCustomHttp(settings.Instance.IpaHost, settings.Instance.IpaRealm, httpClient)

	err := client.RemoteLogin(settings.Instance.IpaUser, settings.Instance.IpaPassword)
	if err != nil {
		return "", err
	}

	var profile string
	switch principalType {
	case "host":
		profile = settings.Instance.CertProfileHost
	case "user":
		profile = settings.Instance.CertProfileUser
	}

	// Request the cert from IPA
	ipaCert, err := client.CertRequest(principal, csr, profile)
	if err != nil {
		return "", err
	}

	// We want the PEM format, so let's convert this
	cert, _ := ipaCert.CertPem()
	if err != nil {
		return "", err
	}

	return cert, nil
}

// Gets status from a PEM-encoded cert
func CertStatus(cert string) (*certv1alpha1.IpaCertData, error) {
	block, _ := pem.Decode([]byte(cert))
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode PEM block containing certificate")
	}

	pub, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	data := certv1alpha1.IpaCertData{
		Serial:   fmt.Sprintf("%X", pub.SerialNumber),
		Subject:  pub.Subject.String(),
		Cn:       pub.Subject.CommonName,
		DnsNames: pub.DNSNames,
		Issued:   metav1.Time{Time: pub.NotBefore},
		Expiry:   metav1.Time{Time: pub.NotAfter},
	}
	return &data, nil
}
