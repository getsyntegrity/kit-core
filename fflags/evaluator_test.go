package fflags

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvalFailClosed_NilEvaluator(t *testing.T) {
	ctx := context.Background()
	assert.False(t, EvalFailClosed(ctx, nil, "key"))
}

func TestEvalFailClosed_Error(t *testing.T) {
	ctx := context.Background()
	eval := &stubEvaluator{err: errors.New("eval error")}
	assert.False(t, EvalFailClosed(ctx, eval, "key"))
}

func TestEvalFailClosed_Disabled(t *testing.T) {
	ctx := context.Background()
	eval := &stubEvaluator{enabled: false}
	assert.False(t, EvalFailClosed(ctx, eval, "key"))
}

func TestEvalFailClosed_Enabled(t *testing.T) {
	ctx := context.Background()
	eval := &stubEvaluator{enabled: true}
	assert.True(t, EvalFailClosed(ctx, eval, "key"))
}

type stubEvaluator struct {
	enabled bool
	err     error
}

func (s *stubEvaluator) IsEnabled(ctx context.Context, key string) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.enabled, nil
}
