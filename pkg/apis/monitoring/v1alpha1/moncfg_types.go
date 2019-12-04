package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MonCfgSpec defines the desired state of MonCfg
// +k8s:openapi-gen=true
type MonCfgSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Global       Global       `json:"global"`
	Prometheus   Prometheus   `json:"prometheus"`
	Alertmanager Alertmanager `json:"alertmanager"`
}

// Global defines global variables
// +k8s:openapi-gen=true
type Global struct {
	NameSpace    string `json:"namespace"`
	SvcAccntName string `json:"serviceAccountName"`
}

// Selector defines the desired state of Selectors
// +k8s:openapi-gen=true
type Selector struct {
	Key    string   `json:"key"`
	Values []string `json:"values,omitempty"`
}

// Alertmanager defines alertmanager config
// +k8s:openapi-gen=true
type Alertmanager struct {
	Name      string      `json:"name"`
	Replicas  int32       `json:"replicas"`
	Resources Resources   `json:"resources"`
	Receivers []Receivers `json:"receivers"`
}

// Receivers defines the desired state of Receivers
// +k8s:openapi-gen=true
type Receivers struct {
	Type   string  `json:"type"`
	Params []Param `json:"params,omitempty"`
}

// Param is a list of alerting receivers.
// +k8s:openapi-gen=true
type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Prometheus defines prometheus config
// +k8s:openapi-gen=true
type Prometheus struct {
	Name              string     `json:"name"`
	Replicas          int32      `json:"replicas"`
	Retention         string     `json:"retention"`
	Resources         Resources  `json:"resources"`
	Port              string     `json:"port"`
	NameSpaceSelector []string   `json:"namespaceselector"`
	Selector          []Selector `json:"selector"`
}

// Resources struct
// +k8s:openapi-gen=true
type Resources struct {
	Requests Requests `json:"requests"`
}

// Requests struct
// +k8s:openapi-gen=true
type Requests struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

// MonCfgStatus defines the observed state of MonCfg
// +k8s:openapi-gen=true
type MonCfgStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MonCfg is the Schema for the moncfgs API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=moncfgs,scope=Namespaced
type MonCfg struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MonCfgSpec   `json:"spec,omitempty"`
	Status MonCfgStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MonCfgList contains a list of MonCfg
type MonCfgList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MonCfg `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MonCfg{}, &MonCfgList{})
}
