package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IpaCertSpec defines the desired state of IpaCert
type IpaCertSpec struct {
	Cn              string   `json:"cn"`
	PrincipalType   string   `json:"type,omitempty"`
	AdditionalNames []string `json:"additionalNames,omitempty"`
	ExcludeChain    bool     `json:"excludeChain,omitempty"`
}

// IpaCertStatus defines the observed state of IpaCert
type IpaCertStatus struct {
	Status       string      `json:"status,omitempty"`
	StatusReason string      `json:"statusReason,omitempty"`
	CertData     IpaCertData `json:",inline"`
}

type IpaCertData struct {
	Serial   string      `json:"serial,omitempty"`
	Issued   metav1.Time `json:"issued,omitempty"`
	Expiry   metav1.Time `json:"expiry,omitempty"`
	Subject  string      `json:"subject,omitempty"`
	Cn       string      `json:"cn,omitempty"`
	DnsNames []string    `json:"dnsNames,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaCert is the Schema for the ipacerts API
// +k8s:openapi-gen=true
type IpaCert struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IpaCertSpec   `json:"spec,omitempty"`
	Status IpaCertStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// IpaCertList contains a list of IpaCert
type IpaCertList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IpaCert `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IpaCert{}, &IpaCertList{})
}
