// Package config provides configuration structures for the security package.
package config

import "time"

// Config represents the unified security configuration.
type Config struct {
	Enabled      bool
	AuthProvider string
	JWT          JWTConfig
	OAuth2       OAuth2Config
	Hydra        HydraConfig
	OAuthkeeper  OAuthkeeperConfig
	APIKey       APIKeyConfig
	Tenant       TenantConfig
	RateLimiting RateLimitConfig
	CORS         CORSConfig
	TLS          TLSConfig
	Endpoints    []EndpointConfig
}

// JWTConfig contains JWT authentication configuration.
type JWTConfig struct {
	Enabled     bool
	Secret      string //nolint:gosec // G117 config field name for JWT secret
	Issuer      string
	Audience    string
	Expiration  int
	RefreshTime int
}

// OAuth2Config contains OAuth2 configuration.
type OAuth2Config struct {
	Enabled      bool
	ClientID     string
	ClientSecret string //nolint:gosec // G117 config field name for OAuth2 client secret
	RedirectURL  string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Scopes       string
}

// HydraConfig contains Ory Hydra configuration.
type HydraConfig struct {
	Enabled      bool
	AdminURL     string
	PublicURL    string
	ClientID     string
	ClientSecret string //nolint:gosec // G117 config field name for Hydra client secret
}

// OAuthkeeperConfig contains OAuthkeeper (Zanzibar) configuration.
type OAuthkeeperConfig struct {
	Enabled     bool
	Endpoint    string
	Namespace   string
	Timeout     time.Duration
	MaxRetries  int
	TenantAware bool
}

// TenantConfig contains multi-tenant configuration.
type TenantConfig struct {
	Enabled       bool
	Mode          string
	HeaderName    string
	QueryParam    string
	DefaultTenant string
	Required      bool
	Validation    bool
}

// APIKeyConfig contains API key configuration.
type APIKeyConfig struct {
	Enabled    bool
	HeaderName string
	ValidKeys  []string
	Rotation   bool
}

// RateLimitConfig contains rate limiting configuration.
type RateLimitConfig struct {
	Enabled           bool
	RequestsPerSecond int
	BurstSize         int
	Window            time.Duration
	PerIP             bool
	PerUser           bool
	PerEndpoint       bool
	PerTenant         bool
}

// CORSConfig contains CORS configuration.
type CORSConfig struct {
	Enabled          bool
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

// TLSConfig contains TLS configuration.
type TLSConfig struct {
	Enabled    bool
	CertFile   string
	KeyFile    string
	CAFile     string
	MinVersion string
	MaxVersion string
}

// EndpointConfig contains endpoint-specific security configuration.
type EndpointConfig struct {
	Path         string
	Methods      []string
	AuthRequired bool
	Roles        []string
	RateLimit    *int
}

// TenantInfo represents tenant information.
type TenantInfo struct {
	ID   string
	Name string
	Mode string
}

// DefaultConfig returns a default security configuration.
func DefaultConfig() Config {
	return Config{
		Enabled:      false,
		AuthProvider: "none",
		OAuthkeeper: OAuthkeeperConfig{
			Enabled:     false,
			Timeout:     5 * time.Second,
			MaxRetries:  3,
			TenantAware: false,
		},
		Tenant: TenantConfig{
			Enabled:       false,
			Mode:          "single",
			HeaderName:    "X-Tenant-ID",
			QueryParam:    "tenant",
			DefaultTenant: "default",
			Required:      false,
			Validation:    false,
		},
		APIKey: APIKeyConfig{
			Enabled:    false,
			HeaderName: "X-API-Key",
		},
		RateLimiting: RateLimitConfig{
			Enabled:           false,
			RequestsPerSecond: 100,
			BurstSize:         200,
			Window:            1 * time.Minute,
		},
		CORS: CORSConfig{
			Enabled:        false,
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"*"},
		},
		TLS:       TLSConfig{Enabled: false},
		Endpoints: []EndpointConfig{},
	}
}
