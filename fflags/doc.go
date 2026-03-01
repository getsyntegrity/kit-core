// Package fflags holds shared feature-flag semantics for the platform.
//
// Responsibilities:
//   - Define minimal runtime interfaces: RuntimeFlagEvaluator (IsEnabled(ctx, key)).
//   - Provide fail-closed helpers: EvalFailClosed.
//
// This package must NOT import Flagsmith or any provider SDK (providers live in service adapters).
package fflags
