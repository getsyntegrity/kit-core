package tenant

import "context"

// TenantContextServiceInterface represents the interface for tenant context operations.
//
//nolint:revive // exported: tenant-prefixed names are intentional for clarity at call sites
type TenantContextServiceInterface interface {
	WithTenantContext(ctx context.Context, tenantCtx TenantContext) context.Context
	GetTenantID(ctx context.Context) (string, error)
	GetUserID(ctx context.Context) (string, error)
	GetCorrelationID(ctx context.Context) (string, error)
	GetCausationID(ctx context.Context) (string, error)
	GetRequestID(ctx context.Context) (string, error)
	GetSessionID(ctx context.Context) (string, error)
	GetTenantContext(ctx context.Context) (TenantContext, error)
	ValidateTenantContext(tenantCtx TenantContext) error
	SetDefaultTenantID(tenantID string) error
	GetDefaultTenantID() string
	IsTenantContextPresent(ctx context.Context) bool
	ClearTenantContext(ctx context.Context) context.Context
	CopyTenantContext(from, to context.Context) context.Context
	MergeTenantContext(from, to context.Context) context.Context
}

// TenantValidator represents a tenant validator interface.
//
//nolint:revive // exported: tenant-prefixed names are intentional for clarity at call sites
type TenantValidator interface {
	ValidateTenantID(tenantID string) error
	ValidateUserID(userID string) error
}
