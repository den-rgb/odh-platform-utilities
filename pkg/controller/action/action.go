package action

import (
	"context"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/opendatahub-io/odh-platform-utilities/api/common"
	"github.com/opendatahub-io/odh-platform-utilities/pkg/controller/conditions"
	ctrlTypes "github.com/opendatahub-io/odh-platform-utilities/pkg/controller/types"
	"github.com/opendatahub-io/odh-platform-utilities/pkg/deploy"
)

// Fn is the action-pipeline function signature. Deploy, GC, and render
// packages produce standalone functions; wrap them in an Fn closure to
// use in a pipeline (see package doc).
type Fn func(ctx context.Context, rr *ReconciliationRequest) error

// ReconciliationRequest carries the shared state for an action pipeline.
type ReconciliationRequest struct {
	Client client.Client

	// Controller may be nil when the pipeline does not use
	// ownership-aware actions.
	Controller ctrlTypes.Controller

	Instance common.PlatformObject

	// Deployer may be nil for render-only pipelines.
	Deployer *deploy.Deployer

	Conditions *conditions.Manager

	// Resources accumulates rendered resources. Render actions append;
	// deploy and GC actions read.
	Resources []unstructured.Unstructured
}
