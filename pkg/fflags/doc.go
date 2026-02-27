// Package fflags holds shared feature-flag semantics for the platform.
//
// Responsibilities:
//   - Define canonical flag keys (constants used across multiple services; see keys.go).
//   - Define minimal runtime interfaces: RuntimeFlagEvaluator (IsEnabled(ctx, key)).
//   - Provide fail-closed helpers: EvalFailClosed.
//
// This package must NOT:
//   - Import Flagsmith or any provider SDK (providers live in service adapters).
//   - Depend on Flagsmith or its HTTP client in go.mod (supply-chain boundary: only services that need feature flags depend on Flagsmith).
//   - Contain service-specific logic (e.g. identity-based evaluation stays in services).
//   - Mention feature names that belong to other domains (e.g. profiling or observability).
//
// No HTTP routing or infra/kit dependencies. Providers implement RuntimeFlagEvaluator
// in each service's internal/adapters/fflags; service-specific keys and guards live in
// each service's internal/fflags.
package fflags
