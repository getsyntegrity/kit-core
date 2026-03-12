package tenant

import (
	"context"
	"strings"
	"testing"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestTenant(t *testing.T) {
	specs.Describe(t, "tenant", func(s *specs.Spec) {
		s.When("NewDefaultTenantValidator", func(s *specs.Spec) {
			s.It("returns non-nil validator with patterns", func(ctx *specs.Context) {
				v := NewDefaultTenantValidator()
				assert.NotNil(ctx.T, v)
				assert.NotNil(ctx.T, v.TenantIDPattern)
				assert.NotNil(ctx.T, v.UserIDPattern)
			})
		})

		s.When("DefaultTenantValidator_ValidateTenantID", func(s *specs.Spec) {
			s.It("accepts valid rejects invalid", func(ctx *specs.Context) {
				v := NewDefaultTenantValidator()
				assert.NoError(ctx.T, v.ValidateTenantID("tenant-1"))
				assert.NoError(ctx.T, v.ValidateTenantID("a"))
				assert.Error(ctx.T, v.ValidateTenantID(""))
				assert.Error(ctx.T, v.ValidateTenantID(string(make([]byte, 65))))
				assert.Error(ctx.T, v.ValidateTenantID("invalid!"))
			})
		})

		s.When("DefaultTenantValidator_ValidateUserID", func(s *specs.Spec) {
			s.It("accepts valid rejects empty and invalid chars", func(ctx *specs.Context) {
				v := NewDefaultTenantValidator()
				assert.NoError(ctx.T, v.ValidateUserID("user-1"))
				assert.Error(ctx.T, v.ValidateUserID(""))
				assert.Error(ctx.T, v.ValidateUserID("invalid!"))
			})
			s.It("rejects user ID longer than 64 characters", func(ctx *specs.Context) {
				v := NewDefaultTenantValidator()
				long := strings.Repeat("a", 65)
				err := v.ValidateUserID(long)
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "64")
			})
		})

		s.When("NewTenantContextService", func(s *specs.Spec) {
			s.It("returns service with default tenant ID", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				assert.NotNil(ctx.T, svc)
				assert.Equal(ctx.T, DefaultTenantID, svc.defaultTenantID)
			})
		})

		s.When("NewTenantContextServiceWithDefaults", func(s *specs.Spec) {
			s.It("sets custom default tenant ID", func(ctx *specs.Context) {
				v := NewDefaultTenantValidator()
				svc := NewTenantContextServiceWithDefaults("custom", v)
				assert.NotNil(ctx.T, svc)
				assert.Equal(ctx.T, "custom", svc.defaultTenantID)
			})
			s.It("accepts nil validator", func(ctx *specs.Context) {
				svc := NewTenantContextServiceWithDefaults("x", nil)
				assert.NotNil(ctx.T, svc)
				assert.NotNil(ctx.T, svc.tenantValidator)
			})
		})

		s.When("TenantContextService_WithTenantContext_GetTenantID", func(s *specs.Spec) {
			s.It("returns tenant ID from context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := context.Background()
				tc := TenantContext{TenantID: "t1", UserID: "u1"}
				bg = svc.WithTenantContext(bg, tc)
				id, err := svc.GetTenantID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "t1", id)
			})
		})

		s.When("TenantContextService_GetTenantID", func(s *specs.Spec) {
			s.It("returns default when context empty", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := context.Background()
				id, err := svc.GetTenantID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, DefaultTenantID, id)
			})
		})

		s.When("TenantContextService_RequireTenantID", func(s *specs.Spec) {
			s.It("errors when empty", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := context.Background()
				_, err := svc.RequireTenantID(bg)
				assert.Error(ctx.T, err)
			})
			s.It("errors when tenant ID is default", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: DefaultTenantID})
				_, err := svc.RequireTenantID(bg)
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "cannot be default")
			})
			s.It("returns ID when set and not default", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "t1"})
				id, err := svc.RequireTenantID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "t1", id)
			})
			s.It("returns error when tenant ID fails validation", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "invalid!"})
				_, err := svc.RequireTenantID(bg)
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "invalid tenant ID")
			})
			s.It("errors when tenant ID equals custom default", func(ctx *specs.Context) {
				svc := NewTenantContextServiceWithDefaults("custom-default", nil)
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "custom-default"})
				_, err := svc.RequireTenantID(bg)
				assert.Error(ctx.T, err)
				assert.Contains(ctx.T, err.Error(), "cannot be default")
			})
		})

		s.When("TenantContextService_GetUserID", func(s *specs.Spec) {
			s.It("returns user ID from context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{UserID: "u1"})
				id, err := svc.GetUserID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "u1", id)
			})
		})

		s.When("TenantContextService_GetCorrelationID", func(s *specs.Spec) {
			s.It("returns correlation ID from context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CorrelationID: "corr-1"})
				id, err := svc.GetCorrelationID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "corr-1", id)
			})
		})

		s.When("TenantContextService_GetCausationID", func(s *specs.Spec) {
			s.It("returns causation ID from context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CausationID: "caus-1"})
				id, err := svc.GetCausationID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "caus-1", id)
			})
		})

		s.When("TenantContextService_GetRequestID", func(s *specs.Spec) {
			s.It("returns request ID from context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{RequestID: "req-1"})
				id, err := svc.GetRequestID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "req-1", id)
			})
		})

		s.When("TenantContextService_GetSessionID", func(s *specs.Spec) {
			s.It("returns session ID from context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{SessionID: "sess-1"})
				id, err := svc.GetSessionID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "sess-1", id)
			})
		})

		s.When("TenantContextService_GetTenantContext", func(s *specs.Spec) {
			s.It("returns full tenant context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				tc := TenantContext{TenantID: "t1", UserID: "u1"}
				bg := svc.WithTenantContext(context.Background(), tc)
				got, err := svc.GetTenantContext(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, tc.TenantID, got.TenantID)
				assert.Equal(ctx.T, tc.UserID, got.UserID)
			})
		})

		s.When("TenantContextService_ValidateTenantContext", func(s *specs.Spec) {
			s.It("passes when tenant ID valid and user ID empty", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				assert.NoError(ctx.T, svc.ValidateTenantContext(TenantContext{TenantID: "t1"}))
			})
			s.It("fails when tenant ID empty", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				assert.Error(ctx.T, svc.ValidateTenantContext(TenantContext{}))
			})
			s.It("fails when user ID invalid", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				assert.Error(ctx.T, svc.ValidateTenantContext(TenantContext{TenantID: "t1", UserID: "invalid!"}))
			})
		})

		s.When("TenantContextService_SetDefaultTenantID_GetDefaultTenantID", func(s *specs.Spec) {
			s.It("sets and returns default tenant ID", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				assert.NoError(ctx.T, svc.SetDefaultTenantID("new-default"))
				assert.Equal(ctx.T, "new-default", svc.GetDefaultTenantID())
			})
			s.It("returns error for invalid default tenant ID", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				err := svc.SetDefaultTenantID("")
				assert.Error(ctx.T, err)
			})
		})

		s.When("TenantContextService_ClearTenantContext", func(s *specs.Spec) {
			s.It("clears context to default", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "t1"})
				bg = svc.ClearTenantContext(bg)
				id, _ := svc.GetTenantID(bg)
				assert.Equal(ctx.T, DefaultTenantID, id)
			})
		})

		s.When("TenantContextService_CopyTenantContext", func(s *specs.Spec) {
			s.It("copies tenant context to another context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "t1", UserID: "u1"})
				bg2 := context.Background()
				bg2 = svc.CopyTenantContext(bg, bg2)
				id, _ := svc.GetTenantID(bg2)
				assert.Equal(ctx.T, "t1", id)
			})
			s.It("copies correlation and causation when present in from", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "t1", CorrelationID: "c1", CausationID: "c2"})
				bg2 := context.Background()
				bg2 = svc.CopyTenantContext(bg, bg2)
				cid, _ := svc.GetCorrelationID(bg2)
				caid, _ := svc.GetCausationID(bg2)
				assert.Equal(ctx.T, "c1", cid)
				assert.Equal(ctx.T, "c2", caid)
			})
			s.It("copies all fields from source to empty context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				from := svc.WithTenantContext(context.Background(), TenantContext{
					TenantID: "t1", UserID: "u1", CorrelationID: "c1",
					CausationID: "c2", RequestID: "r1", SessionID: "s1",
				})
				to := context.Background()
				out := svc.CopyTenantContext(from, to)
				tc, _ := svc.GetTenantContext(out)
				assert.Equal(ctx.T, "t1", tc.TenantID)
				assert.Equal(ctx.T, "u1", tc.UserID)
				assert.Equal(ctx.T, "c1", tc.CorrelationID)
				assert.Equal(ctx.T, "c2", tc.CausationID)
				assert.Equal(ctx.T, "r1", tc.RequestID)
				assert.Equal(ctx.T, "s1", tc.SessionID)
			})
		})

		s.When("TenantContextService_MergeTenantContext", func(s *specs.Spec) {
			s.It("merges from and to contexts with precedence", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				from := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "t1"})
				to := svc.WithTenantContext(context.Background(), TenantContext{UserID: "u2"})
				merged := svc.MergeTenantContext(from, to)
				tc, err := svc.GetTenantContext(merged)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "t1", tc.TenantID)
				assert.Equal(ctx.T, "u2", tc.UserID)
			})
			s.It("overwrites to fields with from when from has values", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				from := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "from-t", CorrelationID: "from-corr"})
				to := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "to-t", CausationID: "to-caus"})
				merged := svc.MergeTenantContext(from, to)
				tc, _ := svc.GetTenantContext(merged)
				assert.Equal(ctx.T, "from-t", tc.TenantID)
				assert.Equal(ctx.T, "from-corr", tc.CorrelationID)
				assert.Equal(ctx.T, "to-caus", tc.CausationID)
			})
			s.It("merges all six context fields when from has all set", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				from := svc.WithTenantContext(context.Background(), TenantContext{
					TenantID: "t1", UserID: "u1", CorrelationID: "c1",
					CausationID: "c2", RequestID: "r1", SessionID: "s1",
				})
				to := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "t0"})
				merged := svc.MergeTenantContext(from, to)
				tc, _ := svc.GetTenantContext(merged)
				assert.Equal(ctx.T, "t1", tc.TenantID)
				assert.Equal(ctx.T, "u1", tc.UserID)
				assert.Equal(ctx.T, "c1", tc.CorrelationID)
				assert.Equal(ctx.T, "c2", tc.CausationID)
				assert.Equal(ctx.T, "r1", tc.RequestID)
				assert.Equal(ctx.T, "s1", tc.SessionID)
			})
			s.It("when to is empty context, merged result has only from's fields", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				from := svc.WithTenantContext(context.Background(), TenantContext{
					TenantID: "f1", UserID: "f2", CorrelationID: "f3",
					CausationID: "f4", RequestID: "f5", SessionID: "f6",
				})
				to := context.Background()
				merged := svc.MergeTenantContext(from, to)
				tc, _ := svc.GetTenantContext(merged)
				assert.Equal(ctx.T, "f1", tc.TenantID)
				assert.Equal(ctx.T, "f2", tc.UserID)
				assert.Equal(ctx.T, "f3", tc.CorrelationID)
				assert.Equal(ctx.T, "f4", tc.CausationID)
				assert.Equal(ctx.T, "f5", tc.RequestID)
				assert.Equal(ctx.T, "f6", tc.SessionID)
			})
		})

		s.When("IsTenantContextPresent", func(s *specs.Spec) {
			s.It("returns false when absent true when present", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := context.Background()
				assert.False(ctx.T, svc.IsTenantContextPresent(bg))
				bg = svc.WithTenantContext(bg, TenantContext{TenantID: "t1"})
				assert.True(ctx.T, svc.IsTenantContextPresent(bg))
			})
		})

		s.When("WithCorrelationID", func(s *specs.Spec) {
			s.It("stores and retrieves correlation ID", func(ctx *specs.Context) {
				bg := context.Background()
				bg = WithCorrelationID(bg, "cid")
				cid, ok := CorrelationIDFromContext(bg)
				assert.True(ctx.T, ok)
				assert.Equal(ctx.T, "cid", cid)
			})
		})

		s.When("WithCausationID", func(s *specs.Spec) {
			s.It("stores and retrieves causation ID", func(ctx *specs.Context) {
				bg := context.Background()
				bg = WithCausationID(bg, "caid")
				caid, ok := CausationIDFromContext(bg)
				assert.True(ctx.T, ok)
				assert.Equal(ctx.T, "caid", caid)
			})
		})

		s.When("CorrelationIDFromContext", func(s *specs.Spec) {
			s.It("returns false when absent and value when from tenant context", func(ctx *specs.Context) {
				bg := context.Background()
				_, ok := CorrelationIDFromContext(bg)
				assert.False(ctx.T, ok)
				svc := NewTenantContextService()
				bg = svc.WithTenantContext(bg, TenantContext{CorrelationID: "x"})
				id, ok := CorrelationIDFromContext(bg)
				assert.True(ctx.T, ok)
				assert.Equal(ctx.T, "x", id)
			})
		})

		s.When("CausationIDFromContext", func(s *specs.Spec) {
			s.It("returns false when absent and value when from tenant context", func(ctx *specs.Context) {
				bg := context.Background()
				_, ok := CausationIDFromContext(bg)
				assert.False(ctx.T, ok)
				svc := NewTenantContextService()
				bg = svc.WithTenantContext(bg, TenantContext{CausationID: "y"})
				id, ok := CausationIDFromContext(bg)
				assert.True(ctx.T, ok)
				assert.Equal(ctx.T, "y", id)
			})
		})

		s.When("EnsureCorrelationID", func(s *specs.Spec) {
			s.It("returns empty when absent and existing when in context", func(ctx *specs.Context) {
				bg := context.Background()
				_, id := EnsureCorrelationID(bg)
				assert.Equal(ctx.T, "", id)
				svc := NewTenantContextService()
				bg = svc.WithTenantContext(bg, TenantContext{CorrelationID: "existing"})
				_, id = EnsureCorrelationID(bg)
				assert.Equal(ctx.T, "existing", id)
			})
		})

		s.When("EnsureCausationID", func(s *specs.Spec) {
			s.It("returns correlation when provided and existing when in context and fallback to correlation", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CorrelationID: "corr"})
				_, causID := EnsureCausationID(bg, "corr")
				assert.Equal(ctx.T, "corr", causID)
				bg = svc.WithTenantContext(context.Background(), TenantContext{CausationID: "existing-caus"})
				_, causID = EnsureCausationID(bg, "")
				assert.Equal(ctx.T, "existing-caus", causID)
				bg = svc.WithTenantContext(context.Background(), TenantContext{CorrelationID: "fallback-corr"})
				_, causID = EnsureCausationID(bg, "")
				assert.Equal(ctx.T, "fallback-corr", causID)
			})
		})

		s.When("GetTenantID package-level", func(s *specs.Spec) {
			s.It("returns tenant ID when present in context", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{TenantID: "pkg-t"})
				id, err := GetTenantID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "pkg-t", id)
			})
			s.It("returns default when not in context", func(ctx *specs.Context) {
				id, err := GetTenantID(context.Background())
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, DefaultTenantID, id)
			})
		})

		s.When("GetCorrelationID package-level", func(s *specs.Spec) {
			s.It("returns correlation ID when present", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CorrelationID: "pkg-corr"})
				id, err := GetCorrelationID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "pkg-corr", id)
			})
			s.It("returns error when missing", func(ctx *specs.Context) {
				_, err := GetCorrelationID(context.Background())
				assert.Error(ctx.T, err)
			})
		})

		s.When("GetCausationID package-level", func(s *specs.Spec) {
			s.It("returns causation ID when present", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CausationID: "pkg-caus"})
				id, err := GetCausationID(bg)
				assert.NoError(ctx.T, err)
				assert.Equal(ctx.T, "pkg-caus", id)
			})
			s.It("returns error when missing", func(ctx *specs.Context) {
				_, err := GetCausationID(context.Background())
				assert.Error(ctx.T, err)
			})
		})

		s.When("ConfigureLoggerCorrelationExtractor", func(s *specs.Spec) {
			s.It("runs without panic and can be called multiple times", func(_ *specs.Context) {
				ConfigureLoggerCorrelationExtractor()
				ConfigureLoggerCorrelationExtractor()
			})
		})

		s.When("EnsureTraceIDs", func(s *specs.Spec) {
			s.It("sets correlation and causation from pointers", func(ctx *specs.Context) {
				corr, caus := "c1", "c2"
				bg := context.Background()
				bg = EnsureTraceIDs(bg, &corr, &caus)
				tc, _ := NewTenantContextService().GetTenantContext(bg)
				assert.Equal(ctx.T, "c1", tc.CorrelationID)
				assert.Equal(ctx.T, "c2", tc.CausationID)
			})
			s.It("handles nil pointers without panic", func(_ *specs.Context) {
				_ = EnsureTraceIDs(context.Background(), nil, nil)
			})
			s.It("fills pointers from context when empty", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CorrelationID: "from-ctx", CausationID: "caus-ctx"})
				var c, c2 string
				_ = EnsureTraceIDs(bg, &c, &c2)
				assert.Equal(ctx.T, "from-ctx", c)
				assert.Equal(ctx.T, "caus-ctx", c2)
			})
			s.It("fills correlation pointer from context when pointer empty", func(ctx *specs.Context) {
				svc := NewTenantContextService()
				bg := svc.WithTenantContext(context.Background(), TenantContext{CorrelationID: "ctx-corr"})
				var corr, caus string
				_ = EnsureTraceIDs(bg, &corr, &caus)
				assert.Equal(ctx.T, "ctx-corr", corr)
			})
			s.It("falls back causation to correlation when both empty in context", func(ctx *specs.Context) {
				corr := "only-corr"
				var caus string
				bg := EnsureTraceIDs(context.Background(), &corr, &caus)
				tc, _ := NewTenantContextService().GetTenantContext(bg)
				assert.Equal(ctx.T, "only-corr", tc.CausationID)
			})
		})
	})
}
