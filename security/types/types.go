// Package types provides shared type definitions for the security package.
package types

// Tenant mode constants.
const (
	TenantModeSingle   = "single"
	TenantModeMulti    = "multi"
	TenantModeHybrid   = "hybrid"
	TenantModeDisabled = "disabled"
	DefaultTenantID    = "default"
)

// TokenIntrospection represents Hydra token introspection response.
type TokenIntrospection struct {
	Active    bool     `json:"active"`
	Scope     string   `json:"scope"`
	ClientID  string   `json:"client_id"`
	Username  string   `json:"username"`
	TokenType string   `json:"token_type"`
	Exp       int64    `json:"exp"`
	Iat       int64    `json:"iat"`
	Sub       string   `json:"sub"`
	Aud       []string `json:"aud"`
	Iss       string   `json:"iss"`
}
