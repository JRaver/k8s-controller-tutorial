package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:generate=true
type FrontendPageSpec struct {
	Content  string `json:"content"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
	Port     int    `json:"port"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=frontendpages,shortName=fp,scope=Namespaced
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
type FrontendPage struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrontendPageSpec   `json:"spec,omitempty"`
	Status FrontendPageStatus `json:"status,omitempty"`
}

// FrontendPageStatus defines the observed state of FrontendPage
type FrontendPageStatus struct {
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
type FrontendPageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []FrontendPage `json:"items"`
}

func init() {
	SchemeBuilder.Register(&FrontendPage{}, &FrontendPageList{})
}
