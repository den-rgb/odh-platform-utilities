package types

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Controller defines the ownership-query contract that deploy and GC
// actions use to handle OwnerReferences and resource tracking.
// Any struct satisfying this interface can drive the action pipeline.
type Controller interface {
	// Owns reports whether the controller manages resources of the given GVK.
	// This includes both static ownership and dynamic ownership registered
	// via AddDynamicOwnedType.
	Owns(gvk schema.GroupVersionKind) bool

	// AddDynamicOwnedType registers a GVK as dynamically owned.
	// Once added, Owns returns true for this GVK.
	AddDynamicOwnedType(gvk schema.GroupVersionKind)

	GetClient() client.Client
	GetDiscoveryClient() discovery.DiscoveryInterface
	GetDynamicClient() dynamic.Interface

	// IsDynamicOwnershipEnabled reports whether dynamic ownership is active.
	// When true, deploy actions may call AddDynamicOwnedType for newly
	// encountered GVKs.
	IsDynamicOwnershipEnabled() bool

	// IsExcludedFromDynamicOwnership reports whether the given GVK should
	// not have owner references set, even when dynamic ownership is enabled.
	IsExcludedFromDynamicOwnership(gvk schema.GroupVersionKind) bool
}
