package listeners

import (
	"context"
	"testing"
	"time"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestListeners(t *testing.T) {
	specs.Describe(t, "listeners", func(s *specs.Spec) {
		s.When("DefaultRetryPolicy", func(s *specs.Spec) {
			s.It("returns expected defaults", func(ctx *specs.Context) {
				p := DefaultRetryPolicy()
				assert.Equal(ctx.T, 5, p.MaxAttempts)
				assert.Equal(ctx.T, 100*time.Millisecond, p.InitialBackoff)
				assert.Equal(ctx.T, 2*time.Second, p.MaxBackoff)
			})
		})

		s.When("WithRecv", func(s *specs.Spec) {
			s.It("accepts non-nil recv function", func(ctx *specs.Context) {
				opt := WithRecv(func(_ context.Context) (Event, error) { return nil, nil })
				err := opt(nil)
				assert.NoError(ctx.T, err)
			})
			s.It("returns error when recv is nil", func(ctx *specs.Context) {
				opt := WithRecv(nil)
				err := opt(nil)
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "cannot be nil")
			})
		})

		s.When("WithHandler", func(s *specs.Spec) {
			s.It("accepts non-nil handler", func(ctx *specs.Context) {
				opt := WithHandler(func(_ context.Context, _ Event) error { return nil })
				err := opt(nil)
				assert.NoError(ctx.T, err)
			})
			s.It("returns error when handler is nil", func(ctx *specs.Context) {
				opt := WithHandler(nil)
				err := opt(nil)
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "cannot be nil")
			})
		})

		s.When("WithCheckpointStore", func(s *specs.Spec) {
			s.It("accepts non-nil load and save", func(ctx *specs.Context) {
				opt := WithCheckpointStore(
					func(_ context.Context) (int64, error) { return 0, nil },
					func(_ context.Context, _ int64) (bool, error) { return true, nil },
				)
				err := opt(nil)
				assert.NoError(ctx.T, err)
			})
			s.It("returns error when load is nil", func(ctx *specs.Context) {
				opt := WithCheckpointStore(nil, func(_ context.Context, _ int64) (bool, error) { return true, nil })
				err := opt(nil)
				assert.Error(ctx.T, err)
			})
			s.It("returns error when save is nil", func(ctx *specs.Context) {
				opt := WithCheckpointStore(func(_ context.Context) (int64, error) { return 0, nil }, nil)
				err := opt(nil)
				assert.Error(ctx.T, err)
			})
		})

		s.When("WithRetryPolicy", func(s *specs.Spec) {
			s.It("accepts retry policy", func(ctx *specs.Context) {
				opt := WithRetryPolicy(RetryPolicy{MaxAttempts: 3})
				err := opt(nil)
				assert.NoError(ctx.T, err)
			})
		})
	})
}
