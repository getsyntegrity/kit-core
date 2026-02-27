package capabilities

// Config holds per-capability enablement for infra policy.
// Callers pass resolved values; the package does not read config or environment.
// The kit does not define debug-surface capabilities; services implement their own policy.
type Config struct{}
