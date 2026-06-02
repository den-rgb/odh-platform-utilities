// Package action defines the core types for the optional action pipeline
// pattern used in ODH module controllers.
//
// The pipeline composes reconciliation steps as a sequence of [Fn]
// functions sharing state through a [ReconciliationRequest].
//
// # Pipeline Usage
//
// Build a pipeline as a slice of [Fn] and execute sequentially.
// The render package has its own [render.Fn] / [render.ReconciliationRequest];
// wrap render calls in an [Fn] closure to bridge the two:
//
//	rr := &action.ReconciliationRequest{
//	    Client:     cli,
//	    Controller: myController,
//	    Instance:   myCR,
//	    Deployer:   deployer,
//	    Conditions: conditions.NewManager(myCR,
//	        string(common.ConditionTypeReady),
//	        string(common.ConditionTypeProvisioningSucceeded),
//	    ),
//	}
//
//	renderStep := action.Fn(func(ctx context.Context, rr *action.ReconciliationRequest) error {
//	    resources, err := kustomize.Render(manifestPath, engineOpts)
//	    if err != nil {
//	        return err
//	    }
//	    rr.Resources = append(rr.Resources, resources...)
//	    return nil
//	})
//
//	deployStep := action.Fn(func(ctx context.Context, rr *action.ReconciliationRequest) error {
//	    return rr.Deployer.Deploy(ctx, deploy.DeployInput{
//	        Client:    rr.Client,
//	        Owner:     rr.Instance,
//	        Resources: rr.Resources,
//	    })
//	})
//
//	pipeline := []action.Fn{renderStep, deployStep}
//
//	for _, step := range pipeline {
//	    if err := step(ctx, rr); err != nil {
//	        return err
//	    }
//	}
//
// # Standalone Usage
//
// Teams that prefer not to use the pipeline can call render, deploy, and
// GC functions directly without constructing a [ReconciliationRequest].
//
// # Design Note
//
// [ReconciliationRequest] is intentionally minimal. Render-specific fields
// (Manifests, Templates, HelmCharts) live in [render.ReconciliationRequest],
// keeping the action pipeline decoupled from rendering concerns.
package action
