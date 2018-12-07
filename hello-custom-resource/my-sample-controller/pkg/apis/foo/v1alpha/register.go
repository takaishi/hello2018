package v1alpha

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/apis/foo"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var SchemeGroupVersion = schema.GroupVersion{Group: foo.GroupName, Version: "v1alpha"}

func Kind(kind string) schema.GroupKind {
	return SchemeGroupVersion.WithKind(kind).GroupKind()
}

func Resource(resource string) schema.GroupResource {
	return SchemeGroupVersion.WithResource(resource).GroupResource()
}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToSCheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion, &Foo{}, &FooList{})
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
