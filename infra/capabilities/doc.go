// Package capabilities provides a placeholder for capability-based infra policy.
//
// Infra packages define ONLY hard platform rules:
//   - Infra MAY define: defaults (safe-by-default), non-bypassable capabilities.
//   - Infra MUST NOT: call feature flags, depend on fflags, or depend on runtime providers.
//
// The kit does not define debug-surface capabilities. Services implement
// their own policy and route gating. Allowed(cfg, capability) is fail-closed (returns false).
package capabilities
