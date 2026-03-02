package tenant

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultTenantValidator(t *testing.T) {
	v := NewDefaultTenantValidator()
	assert.NotNil(t, v)
	assert.NotNil(t, v.TenantIDPattern)
	assert.NotNil(t, v.UserIDPattern)
}

func TestDefaultTenantValidator_ValidateTenantID(t *testing.T) {
	v := NewDefaultTenantValidator()
	assert.NoError(t, v.ValidateTenantID("tenant-1"))
	assert.NoError(t, v.ValidateTenantID("a"))
	assert.Error(t, v.ValidateTenantID(""))
	assert.Error(t, v.ValidateTenantID(string(make([]byte, 65))))
	assert.Error(t, v.ValidateTenantID("invalid!"))
}

func TestDefaultTenantValidator_ValidateUserID(t *testing.T) {
	v := NewDefaultTenantValidator()
	assert.NoError(t, v.ValidateUserID("user-1"))
	assert.Error(t, v.ValidateUserID(""))
	assert.Error(t, v.ValidateUserID("invalid!"))
}

func TestNewTenantContextService(t *testing.T) {
	s := NewTenantContextService()
	assert.NotNil(t, s)
	assert.Equal(t, DefaultTenantID, s.defaultTenantID)
}

func TestNewTenantContextServiceWithDefaults(t *testing.T) {
	v := NewDefaultTenantValidator()
	s := NewTenantContextServiceWithDefaults("custom", v)
	assert.NotNil(t, s)
	assert.Equal(t, "custom", s.defaultTenantID)
}

func TestNewTenantContextServiceWithDefaults_NilValidator(t *testing.T) {
	s := NewTenantContextServiceWithDefaults("x", nil)
	assert.NotNil(t, s)
	assert.NotNil(t, s.tenantValidator)
}

func TestTenantContextService_WithTenantContext_GetTenantID(t *testing.T) {
	s := NewTenantContextService()
	ctx := context.Background()
	tc := TenantContext{TenantID: "t1", UserID: "u1"}
	ctx = s.WithTenantContext(ctx, tc)
	id, err := s.GetTenantID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "t1", id)
}

func TestTenantContextService_GetTenantID_DefaultWhenEmpty(t *testing.T) {
	s := NewTenantContextService()
	ctx := context.Background()
	id, err := s.GetTenantID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, DefaultTenantID, id)
}

func TestTenantContextService_RequireTenantID(t *testing.T) {
	s := NewTenantContextService()
	ctx := context.Background()
	_, err := s.RequireTenantID(ctx)
	assert.Error(t, err)
	ctx = s.WithTenantContext(ctx, TenantContext{TenantID: "t1"})
	id, err := s.RequireTenantID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "t1", id)
}

func TestTenantContextService_GetUserID(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{UserID: "u1"})
	id, err := s.GetUserID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "u1", id)
}

func TestTenantContextService_GetCorrelationID(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{CorrelationID: "corr-1"})
	id, err := s.GetCorrelationID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "corr-1", id)
}

func TestTenantContextService_GetCausationID(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{CausationID: "caus-1"})
	id, err := s.GetCausationID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "caus-1", id)
}

func TestTenantContextService_GetRequestID(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{RequestID: "req-1"})
	id, err := s.GetRequestID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "req-1", id)
}

func TestTenantContextService_GetSessionID(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{SessionID: "sess-1"})
	id, err := s.GetSessionID(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "sess-1", id)
}

func TestTenantContextService_GetTenantContext(t *testing.T) {
	s := NewTenantContextService()
	tc := TenantContext{TenantID: "t1", UserID: "u1"}
	ctx := s.WithTenantContext(context.Background(), tc)
	got, err := s.GetTenantContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, tc.TenantID, got.TenantID)
	assert.Equal(t, tc.UserID, got.UserID)
}

func TestTenantContextService_ValidateTenantContext(t *testing.T) {
	s := NewTenantContextService()
	assert.NoError(t, s.ValidateTenantContext(TenantContext{TenantID: "t1"}))
	assert.Error(t, s.ValidateTenantContext(TenantContext{}))
}

func TestTenantContextService_SetDefaultTenantID_GetDefaultTenantID(t *testing.T) {
	s := NewTenantContextService()
	s.SetDefaultTenantID("new-default")
	assert.Equal(t, "new-default", s.GetDefaultTenantID())
}

func TestTenantContextService_ClearTenantContext(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{TenantID: "t1"})
	ctx = s.ClearTenantContext(ctx)
	id, _ := s.GetTenantID(ctx)
	assert.Equal(t, DefaultTenantID, id)
}

func TestTenantContextService_CopyTenantContext(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{TenantID: "t1", UserID: "u1"})
	ctx2 := context.Background()
	ctx2 = s.CopyTenantContext(ctx, ctx2)
	id, _ := s.GetTenantID(ctx2)
	assert.Equal(t, "t1", id)
}

func TestTenantContextService_MergeTenantContext(t *testing.T) {
	s := NewTenantContextService()
	from := s.WithTenantContext(context.Background(), TenantContext{TenantID: "t1"})
	to := s.WithTenantContext(context.Background(), TenantContext{UserID: "u2"})
	ctx := s.MergeTenantContext(from, to)
	tc, err := s.GetTenantContext(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "t1", tc.TenantID)
	assert.Equal(t, "u2", tc.UserID)
}

func TestIsTenantContextPresent(t *testing.T) {
	s := NewTenantContextService()
	ctx := context.Background()
	assert.False(t, s.IsTenantContextPresent(ctx))
	ctx = s.WithTenantContext(ctx, TenantContext{TenantID: "t1"})
	assert.True(t, s.IsTenantContextPresent(ctx))
}

func TestWithCorrelationID(t *testing.T) {
	ctx := context.Background()
	ctx = WithCorrelationID(ctx, "cid")
	cid, ok := CorrelationIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "cid", cid)
}

func TestWithCausationID(t *testing.T) {
	ctx := context.Background()
	ctx = WithCausationID(ctx, "caid")
	caid, ok := CausationIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "caid", caid)
}

func TestCorrelationIDFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := CorrelationIDFromContext(ctx)
	assert.False(t, ok)
	s := NewTenantContextService()
	ctx = s.WithTenantContext(ctx, TenantContext{CorrelationID: "x"})
	id, ok := CorrelationIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "x", id)
}

func TestCausationIDFromContext(t *testing.T) {
	ctx := context.Background()
	_, ok := CausationIDFromContext(ctx)
	assert.False(t, ok)
	s := NewTenantContextService()
	ctx = s.WithTenantContext(ctx, TenantContext{CausationID: "y"})
	id, ok := CausationIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "y", id)
}

func TestEnsureCorrelationID(t *testing.T) {
	ctx := context.Background()
	_, id := EnsureCorrelationID(ctx)
	assert.Equal(t, "", id)
	// when context has correlation ID
	s := NewTenantContextService()
	ctx = s.WithTenantContext(ctx, TenantContext{CorrelationID: "existing"})
	ctx, id = EnsureCorrelationID(ctx)
	assert.Equal(t, "existing", id)
}

func TestEnsureCausationID(t *testing.T) {
	s := NewTenantContextService()
	ctx := s.WithTenantContext(context.Background(), TenantContext{CorrelationID: "corr"})
	ctx, causID := EnsureCausationID(ctx, "corr")
	assert.Equal(t, "corr", causID)
	// existing causation in context
	ctx = s.WithTenantContext(context.Background(), TenantContext{CausationID: "existing-caus"})
	_, causID = EnsureCausationID(ctx, "")
	assert.Equal(t, "existing-caus", causID)
	// fallback to correlation when causation empty
	ctx = s.WithTenantContext(context.Background(), TenantContext{CorrelationID: "fallback-corr"})
	_, causID = EnsureCausationID(ctx, "")
	assert.Equal(t, "fallback-corr", causID)
}

func TestConfigureLoggerCorrelationExtractor(t *testing.T) {
	ConfigureLoggerCorrelationExtractor() // no-op, just ensure no panic
}

func TestEnsureTraceIDs(t *testing.T) {
	corr, caus := "c1", "c2"
	ctx := context.Background()
	ctx = EnsureTraceIDs(ctx, &corr, &caus)
	tc, _ := NewTenantContextService().GetTenantContext(ctx)
	assert.Equal(t, "c1", tc.CorrelationID)
	assert.Equal(t, "c2", tc.CausationID)
	// nil pointers - no panic
	ctx = EnsureTraceIDs(context.Background(), nil, nil)
	// empty pointers get filled from context
	s := NewTenantContextService()
	ctx = s.WithTenantContext(context.Background(), TenantContext{CorrelationID: "from-ctx", CausationID: "caus-ctx"})
	var c, c2 string
	ctx = EnsureTraceIDs(ctx, &c, &c2)
	assert.Equal(t, "from-ctx", c)
	assert.Equal(t, "caus-ctx", c2)
}
