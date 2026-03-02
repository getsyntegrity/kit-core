package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	c := DefaultConfig()
	assert.False(t, c.Enabled)
	assert.Equal(t, "none", c.AuthProvider)
	assert.False(t, c.OAuthkeeper.Enabled)
	assert.Equal(t, 5*time.Second, c.OAuthkeeper.Timeout)
	assert.Equal(t, 3, c.OAuthkeeper.MaxRetries)
	assert.False(t, c.Tenant.Enabled)
	assert.Equal(t, "X-Tenant-ID", c.Tenant.HeaderName)
	assert.False(t, c.APIKey.Enabled)
	assert.False(t, c.RateLimiting.Enabled)
	assert.Equal(t, 100, c.RateLimiting.RequestsPerSecond)
	assert.False(t, c.CORS.Enabled)
	assert.False(t, c.TLS.Enabled)
}
