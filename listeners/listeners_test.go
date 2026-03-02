package listeners

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultRetryPolicy(t *testing.T) {
	p := DefaultRetryPolicy()
	assert.Equal(t, 5, p.MaxAttempts)
	assert.Equal(t, 100*time.Millisecond, p.InitialBackoff)
	assert.Equal(t, 2*time.Second, p.MaxBackoff)
}

func TestWithRecv(t *testing.T) {
	opt := WithRecv(func(ctx context.Context) (Event, error) { return nil, nil })
	err := opt(nil)
	assert.NoError(t, err)
}

func TestWithRecv_Nil(t *testing.T) {
	opt := WithRecv(nil)
	err := opt(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestWithHandler(t *testing.T) {
	opt := WithHandler(func(ctx context.Context, event Event) error { return nil })
	err := opt(nil)
	assert.NoError(t, err)
}

func TestWithHandler_Nil(t *testing.T) {
	opt := WithHandler(nil)
	err := opt(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestWithCheckpointStore(t *testing.T) {
	opt := WithCheckpointStore(
		func(ctx context.Context) (int64, error) { return 0, nil },
		func(ctx context.Context, version int64) (bool, error) { return true, nil },
	)
	err := opt(nil)
	assert.NoError(t, err)
}

func TestWithCheckpointStore_NilLoad(t *testing.T) {
	opt := WithCheckpointStore(nil, func(ctx context.Context, version int64) (bool, error) { return true, nil })
	err := opt(nil)
	assert.Error(t, err)
}

func TestWithCheckpointStore_NilSave(t *testing.T) {
	opt := WithCheckpointStore(func(ctx context.Context) (int64, error) { return 0, nil }, nil)
	err := opt(nil)
	assert.Error(t, err)
}

func TestWithRetryPolicy(t *testing.T) {
	opt := WithRetryPolicy(RetryPolicy{MaxAttempts: 3})
	err := opt(nil)
	assert.NoError(t, err)
}
