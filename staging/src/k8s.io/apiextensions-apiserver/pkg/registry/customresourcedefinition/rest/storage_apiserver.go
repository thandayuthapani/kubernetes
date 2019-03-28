package rest

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/registry/rest"
	//"k8s.io/apimachinery/pkg/runtime/schema"
	genericapiserver "k8s.io/apiserver/pkg/server"
	serverstorage "k8s.io/apiserver/pkg/server/storage"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/install"
	"k8s.io/apiextensions-apiserver/pkg/registry/customresourcedefinition"
	customresource "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var(
	Scheme = runtime.NewScheme()
)
type RESTStorageProvider struct{}

func (p RESTStorageProvider) NewRESTStorage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) (genericapiserver.APIGroupInfo, bool) {
	apiGroupInfo := genericapiserver.NewDefaultAPIGroupInfo(customresource.GroupName, legacyscheme.Scheme, legacyscheme.ParameterCodec, legacyscheme.Codecs)
	// If you add a version here, be sure to add an entry in `k8s.io/kubernetes/cmd/kube-apiserver/app/aggregator.go with specific priorities.
	// TODO refactor the plumbing to provide the information in the APIGroupInfo

	if apiResourceConfigSource.VersionEnabled(customresource.SchemeGroupVersion) {
		apiGroupInfo.VersionedResourcesStorageMap[customresource.SchemeGroupVersion.Version] = p.v1beta1Storage(apiResourceConfigSource, restOptionsGetter)
	}
	return apiGroupInfo, true
}

func (p RESTStorageProvider) v1beta1Storage(apiResourceConfigSource serverstorage.APIResourceConfigSource, restOptionsGetter generic.RESTOptionsGetter) map[string]rest.Storage {
	storage := map[string]rest.Storage{}
	// validatingwebhookconfigurations
	//Scheme.AddKnownTypes(schema.GroupVersion{Group:"apiextensions.k8s.io",Version:"v1beta1"})
	//Scheme.Print()
	install.Install(Scheme)
	customResourceDefintionStorage := customresourcedefinition.NewREST(Scheme, restOptionsGetter)
	storage["customresourcedefinitions"] = customResourceDefintionStorage
	storage["customresourcedefinitions/status"] = customresourcedefinition.NewStatusREST(Scheme, customResourceDefintionStorage)

	return storage
}

func (p RESTStorageProvider) GroupName() string {
	return customresource.GroupName
}

