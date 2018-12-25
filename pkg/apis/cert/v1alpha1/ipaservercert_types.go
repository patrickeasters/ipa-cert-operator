package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IpaServerCertSpec defines the desired state of IpaServerCert
type IpaServerCertSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// IpaServerCertStatus defines the observed state of IpaServerCert
type IpaServerCertStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaServerCert is the Schema for the ipaservercerts API
// +k8s:openapi-gen=true
type IpaServerCert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpaServerCertSpec   `json:"spec,omitempty"`
	Status IpaServerCertStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaServerCertList contains a list of IpaServerCert
type IpaServerCertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpaServerCert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpaServerCert{}, &IpaServerCertList{})
}
