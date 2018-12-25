package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IpaClientCertSpec defines the desired state of IpaClientCert
type IpaClientCertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// IpaClientCertStatus defines the observed state of IpaClientCert
type IpaClientCertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaClientCert is the Schema for the ipaclientcerts API
// +k8s:openapi-gen=true
type IpaClientCert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpaClientCertSpec   `json:"spec,omitempty"`
	Status IpaClientCertStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaClientCertList contains a list of IpaClientCert
type IpaClientCertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpaClientCert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpaClientCert{}, &IpaClientCertList{})
}
