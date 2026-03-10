package fflags

import (
	"context"
	"errors"
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestEvalFailClosed(t *testing.T) {
	specs.Describe(t, "EvalFailClosed", func(s *specs.Spec) {
		s.When("evaluator is nil", func(s *specs.Spec) {
			s.It("returns false", func(ctx *specs.Context) {
				bg := context.Background()
				assert.False(ctx.T, EvalFailClosed(bg, nil, "key"))
			})
		})
		s.When("evaluator returns error", func(s *specs.Spec) {
			s.It("returns false", func(ctx *specs.Context) {
				bg := context.Background()
				eval := &stubEvaluator{err: errors.New("eval error")}
				assert.False(ctx.T, EvalFailClosed(bg, eval, "key"))
			})
		})
		s.When("evaluator returns disabled", func(s *specs.Spec) {
			s.It("returns false", func(ctx *specs.Context) {
				bg := context.Background()
				eval := &stubEvaluator{enabled: false}
				assert.False(ctx.T, EvalFailClosed(bg, eval, "key"))
			})
		})
		s.When("evaluator returns enabled", func(s *specs.Spec) {
			s.It("returns true", func(ctx *specs.Context) {
				bg := context.Background()
				eval := &stubEvaluator{enabled: true}
				assert.True(ctx.T, EvalFailClosed(bg, eval, "key"))
			})
		})
	})
}

type stubEvaluator struct {
	enabled bool
	err     error
}

func (s *stubEvaluator) IsEnabled(_ context.Context, _ string) (bool, error) {
	if s.err != nil {
		return false, s.err
	}
	return s.enabled, nil
}
