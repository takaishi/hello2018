package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Status            FooStatus `json:"status"`
	Spec              FooSpec   `json:"spec"`
}

type FooStatus struct {
	Name string `json:"name"`
}

type FooSpec struct {
	Name string `json:"name"`
}

type FooLIst struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items           []Foo `json:"items"`
}
