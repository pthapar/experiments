package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// TestSpec defines specifications of a test
type TestSpec struct {
	// Image defines the test image along with the tag
	Image string `json:"image"`

	// RunCmd defines the cmd used for running the test inside given image
	RunCmd string `json:"runCmd"`

	// AssertCmd defines the assertion cmd for cheking success or pass of the test
	AssertCmd string `json:"assertCmd"`

	// Local defines whether or not this test is run locally on target infra or not
	// If the test is run locally, then it assumes a super user/container role
	Local bool `json:"local"`
}

// TestStatus is the test status
type TestStatus struct {
	// State of the test like Running, Failed, Successful
	State string `json:"state"`
}

// InfraSpec is the spec for infra
type InfraSpec struct {
	// BootImage defines the image w/ tag/version used to boot up the cluster
	BootImage string `json:"boolImage"`

	// Cloud defines the cloud where this infra exists
	Cloud string `json:"cloud"`
}

// OutPutSpec specifies the details for output of tests
type OutPutSpec struct {
	// Service is the  service used to put output to
	Service string `json:"service"`

	// Provider is the provider of the given provider
	Provider string `json:"provider"`

	// Retention is the period is number of days for which the output is available
	Retention int `json:"retention"`

	// Name is the identifier of the given output endpoint. Example s3 bucket name, kafka topic name
	Name string `json:"name"`
}

// TestSetSpec defines the desired state of TestSet
type TestSetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Tests []TestSpec `json:"tests"`

	// InfraSpec is the infra spec for running the test
	// +optional
	Infra *InfraSpec `json:"infra"`

	// Output is the output specification for the test set
	Output OutPutSpec `json:"output"`

	// Custom labels to be added to the tests associated with this test set.
	// +optional
	Labels map[string]string `json:"labels,omitempty"`
}

// TestSetStatus defines the observed state of TestSet
type TestSetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
	Status map[string]TestSetStatus `json:status`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TestSet is the Schema for the testsets API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=testsets,scope=Namespaced
type TestSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestSetSpec   `json:"spec,omitempty"`
	Status TestSetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TestSetList contains a list of TestSet
type TestSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TestSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&TestSet{}, &TestSetList{})
}
