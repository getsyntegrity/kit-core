package fflags

import "context"

// RuntimeFlagEvaluator evaluates a feature flag at runtime by key.
type RuntimeFlagEvaluator interface {
	// IsEnabled returns whether the flag for the given key is enabled.
	IsEnabled(ctx context.Context, key string) (bool, error)
}

// EvalFailClosed evaluates the flag for key and returns false on any error or nil evaluator (fail-closed).
func EvalFailClosed(ctx context.Context, eval RuntimeFlagEvaluator, key string) bool {
	if eval == nil {
		return false
	}
	enabled, err := eval.IsEnabled(ctx, key)
	if err != nil {
		return false
	}
	return enabled
}
