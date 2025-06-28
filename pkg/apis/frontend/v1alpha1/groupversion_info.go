// +k8s:deepcopy-gen=package
// +groupName=frontend.jraver.io
package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

var (
	// SchemeGroupVersion is the group version used to register these objects
	SchemeGroupVersion = schema.GroupVersion{Group: "frontend.jraver.io", Version: "v1alpha1"}
	// SchemeBuilder is used to add go types to the scheme
	SchemeBuilder      = &scheme.Builder{GroupVersion: SchemeGroupVersion}
	// AddToScheme is used to add the types to the scheme pkg/client/...
	AddToScheme        = SchemeBuilder.AddToScheme
)
