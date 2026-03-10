package config

import (
	"testing"
	"time"

	"github.com/pablogore/go-specs/specs"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	specs.Describe(t, "config", func(s *specs.Spec) {
		s.When("DefaultConfig", func(s *specs.Spec) {
			s.It("returns disabled config with expected defaults", func(ctx *specs.Context) {
				c := DefaultConfig()
				assert.False(ctx.T, c.Enabled)
				assert.Equal(ctx.T, "none", c.AuthProvider)
				assert.False(ctx.T, c.OAuthkeeper.Enabled)
				assert.Equal(ctx.T, 5*time.Second, c.OAuthkeeper.Timeout)
				assert.Equal(ctx.T, 3, c.OAuthkeeper.MaxRetries)
				assert.False(ctx.T, c.Tenant.Enabled)
				assert.Equal(ctx.T, "X-Tenant-ID", c.Tenant.HeaderName)
				assert.False(ctx.T, c.APIKey.Enabled)
				assert.False(ctx.T, c.RateLimiting.Enabled)
				assert.Equal(ctx.T, 100, c.RateLimiting.RequestsPerSecond)
				assert.False(ctx.T, c.CORS.Enabled)
				assert.False(ctx.T, c.TLS.Enabled)
			})
		})
	})
}
