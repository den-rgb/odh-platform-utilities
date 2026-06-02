// Package types defines the [Controller] interface — the ownership-query
// contract shared between the action pipeline, deploy, and GC packages.
//
// The interface mirrors the Controller contract in the ODH operator so
// that module teams can satisfy it with their existing reconciler without
// code changes.
package types
