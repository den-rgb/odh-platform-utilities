package types_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"

	ctrlTypes "github.com/opendatahub-io/odh-platform-utilities/pkg/controller/types"
)

// stubController is a minimal implementation proving the interface can
// be satisfied without importing the reconciler builder.
type stubController struct {
	cli        client.Client
	disc       discovery.DiscoveryInterface
	dyn        dynamic.Interface
	owned      map[schema.GroupVersionKind]struct{}
	excluded   map[schema.GroupVersionKind]struct{}
	dynOwned   sync.Map
	dynEnabled bool
}

func (s *stubController) Owns(gvk schema.GroupVersionKind) bool {
	if _, ok := s.owned[gvk]; ok {
		return true
	}

	_, ok := s.dynOwned.Load(gvk)

	return ok
}

func (s *stubController) AddDynamicOwnedType(gvk schema.GroupVersionKind) {
	s.dynOwned.Store(gvk, struct{}{})
}

func (s *stubController) GetClient() client.Client                         { return s.cli }
func (s *stubController) GetDiscoveryClient() discovery.DiscoveryInterface { return s.disc }
func (s *stubController) GetDynamicClient() dynamic.Interface              { return s.dyn }
func (s *stubController) IsDynamicOwnershipEnabled() bool                  { return s.dynEnabled }

func (s *stubController) IsExcludedFromDynamicOwnership(gvk schema.GroupVersionKind) bool {
	_, ok := s.excluded[gvk]
	return ok
}

// Compile-time verification that stubController satisfies Controller.
var _ ctrlTypes.Controller = (*stubController)(nil)

func TestController_MockSatisfiesInterface(t *testing.T) {
	t.Parallel()

	var ctrl ctrlTypes.Controller = &stubController{
		owned: map[schema.GroupVersionKind]struct{}{
			{Group: "apps", Version: "v1", Kind: "Deployment"}: {},
		},
		excluded: map[schema.GroupVersionKind]struct{}{
			{Group: "", Version: "v1", Kind: "Namespace"}: {},
		},
		dynEnabled: true,
	}

	require.NotNil(t, ctrl)
}

func TestController_Owns_StaticAndDynamic(t *testing.T) {
	t.Parallel()

	deploymentGVK := schema.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	serviceGVK := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}
	configMapGVK := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}

	ctrl := &stubController{
		owned: map[schema.GroupVersionKind]struct{}{
			deploymentGVK: {},
		},
	}

	assert.True(t, ctrl.Owns(deploymentGVK), "static owned GVK should return true")
	assert.False(t, ctrl.Owns(serviceGVK), "unregistered GVK should return false")

	ctrl.AddDynamicOwnedType(serviceGVK)
	assert.True(t, ctrl.Owns(serviceGVK), "dynamically added GVK should return true")
	assert.False(t, ctrl.Owns(configMapGVK), "other GVKs remain unowned")
}

func TestController_DynamicOwnership(t *testing.T) {
	t.Parallel()

	namespaceGVK := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}
	serviceGVK := schema.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}

	ctrl := &stubController{
		owned: map[schema.GroupVersionKind]struct{}{},
		excluded: map[schema.GroupVersionKind]struct{}{
			namespaceGVK: {},
		},
		dynEnabled: true,
	}

	assert.True(t, ctrl.IsDynamicOwnershipEnabled())
	assert.True(t, ctrl.IsExcludedFromDynamicOwnership(namespaceGVK))
	assert.False(t, ctrl.IsExcludedFromDynamicOwnership(serviceGVK))
}

func TestController_ClientAccessors_NilSafe(t *testing.T) {
	t.Parallel()

	ctrl := &stubController{}

	assert.Nil(t, ctrl.GetClient())
	assert.Nil(t, ctrl.GetDiscoveryClient())
	assert.Nil(t, ctrl.GetDynamicClient())
}
