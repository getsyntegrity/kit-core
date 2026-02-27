package fflags

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalFailClosed_NilEvaluator_ReturnsFalse(t *testing.T) {
	got := EvalFailClosed(context.Background(), nil, "any-key")
	assert.False(t, got)
}

func TestEvalFailClosed_EvaluatorReturnsError_ReturnsFalse(t *testing.T) {
	eval := &mockEvaluator{enabled: true, err: errors.New("provider error")}
	got := EvalFailClosed(context.Background(), eval, "key")
	assert.False(t, got)
}

// TestEvalFailClosed_EvaluatorReturnsTimeout_ReturnsFalse asserts timeout (context.DeadlineExceeded) is treated as deny (fail-closed).
func TestEvalFailClosed_EvaluatorReturnsTimeout_ReturnsFalse(t *testing.T) {
	eval := &mockEvaluator{enabled: true, err: context.DeadlineExceeded}
	got := EvalFailClosed(context.Background(), eval, "key")
	assert.False(t, got)
}

func TestEvalFailClosed_EvaluatorReturnsFalse_ReturnsFalse(t *testing.T) {
	eval := &mockEvaluator{enabled: false, err: nil}
	got := EvalFailClosed(context.Background(), eval, "key")
	assert.False(t, got)
}

func TestEvalFailClosed_EvaluatorReturnsTrue_ReturnsTrue(t *testing.T) {
	eval := &mockEvaluator{enabled: true, err: nil}
	got := EvalFailClosed(context.Background(), eval, "key")
	assert.True(t, got)
}

type mockEvaluator struct {
	enabled bool
	err     error
}

func (m *mockEvaluator) IsEnabled(ctx context.Context, key string) (bool, error) {
	if m.err != nil {
		return false, m.err
	}
	return m.enabled, nil
}

var _ RuntimeFlagEvaluator = (*mockEvaluator)(nil)
