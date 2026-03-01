// Package securityinterfaces provides security-related interfaces.
package securityinterfaces

import (
	"context"

	"github.com/getsyntegrity/kit-core/security/config"
	"github.com/getsyntegrity/kit-core/security/types"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
)

// SecurityManagerInterface defines the interface for the security manager.
type SecurityManagerInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	GetHTTPMiddleware() HTTPMiddlewareInterface
	GetGRPCInterceptors() GRPCInterceptorsInterface
	GetTenantManager() TenantManagerInterface
	GetHydraClient() HydraClientInterface
	GetOAuthkeeperClient() OAuthkeeperClientInterface
	HTTPMiddleware() echo.MiddlewareFunc
	GRPCUnaryInterceptor() grpc.UnaryServerInterceptor
	GRPCStreamInterceptor() grpc.StreamServerInterceptor
	IsStarted() bool
	IsStopped() bool
	GetConfig() config.Config
	UpdateConfig(config config.Config)
}

// HydraClientInterface defines the interface for Hydra client operations.
type HydraClientInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	IntrospectToken(ctx context.Context, token string) (*types.TokenIntrospection, error)
	ValidateToken(ctx context.Context, token string) (map[string]interface{}, error)
	GetUserInfo(ctx context.Context, token string) (map[string]interface{}, error)
	CheckScope(token *types.TokenIntrospection, requiredScope string) bool
	CheckScopes(token *types.TokenIntrospection, requiredScopes []string) bool
}

// OAuthkeeperClientInterface defines the interface for OAuthkeeper client operations.
type OAuthkeeperClientInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
	CheckPermission(ctx context.Context, namespace, object, relation, subject string) (bool, error)
	CheckUserPermission(ctx context.Context, userID, action, resource string) (bool, error)
	CheckTenantPermission(ctx context.Context, userID, tenantID, namespace, object, relation string) (bool, error)
	CheckPermissionWithMetadata(ctx context.Context, namespace, object, relation, subject string, metadata map[string]string) (bool, error)
	CreateRelationship(ctx context.Context, namespace, object, relation, subject string) error
	DeleteRelationship(ctx context.Context, namespace, object, relation, subject string) error
	GrantUserPermission(ctx context.Context, userID, relation, object string) error
	RevokeUserPermission(ctx context.Context, userID, relation, object string) error
	IsEnabled() bool
	GetConfig() config.OAuthkeeperConfig
}

// TenantManagerInterface defines the interface for tenant management operations.
type TenantManagerInterface interface {
	IsSingleTenant() bool
	IsMultiTenant() bool
	IsHybridTenant() bool
	GetDefaultTenant() string
	ResolveTenant(c echo.Context) (*config.TenantInfo, error)
	ResolveTenantGRPC(ctx context.Context) (*config.TenantInfo, error)
	AddTenantToContext(ctx context.Context, tenantInfo *config.TenantInfo) context.Context
	AddTenantToContextWithMetadata(ctx context.Context, tenantInfo *config.TenantInfo, metadata map[string]string) context.Context
	GetTenantFromContext(ctx context.Context) (*config.TenantInfo, error)
	AddTenantToGRPCMetadata(ctx context.Context, tenantInfo *config.TenantInfo) context.Context
	AddTenantToResponse(c echo.Context, tenantInfo *config.TenantInfo)
}

// HTTPMiddlewareInterface defines the interface for HTTP middleware operations.
type HTTPMiddlewareInterface interface {
	SecurityMiddleware() echo.MiddlewareFunc
	CORSMiddleware() echo.MiddlewareFunc
	AuthMiddleware() echo.MiddlewareFunc
	RateLimitMiddleware() echo.MiddlewareFunc
	TenantMiddleware() echo.MiddlewareFunc
	Middleware() echo.MiddlewareFunc
}

// GRPCInterceptorsInterface defines the interface for gRPC interceptors operations.
type GRPCInterceptorsInterface interface {
	UnaryAuthInterceptor() grpc.UnaryServerInterceptor
	StreamAuthInterceptor() grpc.StreamServerInterceptor
	UnaryClientInterceptor() grpc.UnaryClientInterceptor
	StreamClientInterceptor() grpc.StreamClientInterceptor
	UnaryInterceptor() grpc.UnaryServerInterceptor
	StreamInterceptor() grpc.StreamServerInterceptor
}

// RateLimiterInterface defines the interface for rate limiting operations.
type RateLimiterInterface interface {
	Check() error
}
