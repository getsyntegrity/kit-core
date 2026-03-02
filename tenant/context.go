// Package tenant provides tenant context management.
package tenant

import (
	"context"
	"fmt"
	"regexp"
)

// TenantContextKey is the key used to store tenant information in context.
type TenantContextKey string

const (
	TenantIDKey      TenantContextKey = "tenant_id"
	UserIDKey        TenantContextKey = "user_id"
	CorrelationIDKey TenantContextKey = "correlation_id"
	CausationIDKey   TenantContextKey = "causation_id"
	RequestIDKey     TenantContextKey = "request_id"
	SessionIDKey     TenantContextKey = "session_id"
	DefaultTenantID                   = "default"
)

// TenantContext represents tenant context information.
type TenantContext struct {
	TenantID      string `json:"tenant_id"`
	UserID        string `json:"user_id"`
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"`
	RequestID     string `json:"request_id"`
	SessionID     string `json:"session_id"`
}

// TenantContextService provides methods for managing tenant context.
type TenantContextService struct {
	defaultTenantID string
	tenantValidator TenantValidator
}

// DefaultTenantValidator provides default tenant validation.
type DefaultTenantValidator struct {
	TenantIDPattern *regexp.Regexp
	UserIDPattern   *regexp.Regexp
}

// NewDefaultTenantValidator creates a new default tenant validator.
func NewDefaultTenantValidator() *DefaultTenantValidator {
	return &DefaultTenantValidator{
		TenantIDPattern: regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`),
		UserIDPattern:   regexp.MustCompile(`^[a-zA-Z0-9_-]{1,64}$`),
	}
}

// ValidateTenantID validates a tenant ID.
func (v *DefaultTenantValidator) ValidateTenantID(tenantID string) error {
	if tenantID == "" {
		return fmt.Errorf("tenant ID cannot be empty")
	}
	if len(tenantID) > 64 {
		return fmt.Errorf("tenant ID cannot be longer than 64 characters")
	}
	if !v.TenantIDPattern.MatchString(tenantID) {
		return fmt.Errorf("tenant ID contains invalid characters")
	}
	return nil
}

// ValidateUserID validates a user ID.
func (v *DefaultTenantValidator) ValidateUserID(userID string) error {
	if userID == "" {
		return fmt.Errorf("user ID cannot be empty")
	}
	if len(userID) > 64 {
		return fmt.Errorf("user ID cannot be longer than 64 characters")
	}
	if !v.UserIDPattern.MatchString(userID) {
		return fmt.Errorf("user ID contains invalid characters")
	}
	return nil
}

// NewTenantContextService creates a new tenant context service.
func NewTenantContextService() *TenantContextService {
	return &TenantContextService{
		defaultTenantID: DefaultTenantID,
		tenantValidator: NewDefaultTenantValidator(),
	}
}

// NewTenantContextServiceWithDefaults creates a new tenant context service with custom defaults.
func NewTenantContextServiceWithDefaults(defaultTenantID string, validator TenantValidator) *TenantContextService {
	if validator == nil {
		validator = NewDefaultTenantValidator()
	}
	return &TenantContextService{
		defaultTenantID: defaultTenantID,
		tenantValidator: validator,
	}
}

// WithTenantContext adds tenant context to the given context.
func (s *TenantContextService) WithTenantContext(ctx context.Context, tenantCtx TenantContext) context.Context {
	ctx = context.WithValue(ctx, TenantIDKey, tenantCtx.TenantID)
	ctx = context.WithValue(ctx, UserIDKey, tenantCtx.UserID)
	ctx = context.WithValue(ctx, CorrelationIDKey, tenantCtx.CorrelationID)
	ctx = context.WithValue(ctx, CausationIDKey, tenantCtx.CausationID)
	ctx = context.WithValue(ctx, RequestIDKey, tenantCtx.RequestID)
	ctx = context.WithValue(ctx, SessionIDKey, tenantCtx.SessionID)
	return ctx
}

// GetTenantID extracts tenant ID from context.
func (s *TenantContextService) GetTenantID(ctx context.Context) (string, error) {
	tenantID, ok := ctx.Value(TenantIDKey).(string)
	if !ok || tenantID == "" {
		return s.defaultTenantID, nil
	}
	return tenantID, nil
}

// RequireTenantID extracts and validates tenant ID from context.
func (s *TenantContextService) RequireTenantID(ctx context.Context) (string, error) {
	tenantID, err := s.GetTenantID(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get tenant ID: %w", err)
	}
	if tenantID == "" {
		return "", fmt.Errorf("tenant ID is required but was empty")
	}
	if tenantID == DefaultTenantID || tenantID == s.defaultTenantID {
		return "", fmt.Errorf("tenant ID cannot be default tenant ID: %s", tenantID)
	}
	if err := s.tenantValidator.ValidateTenantID(tenantID); err != nil {
		return "", fmt.Errorf("invalid tenant ID: %w", err)
	}
	return tenantID, nil
}

// GetUserID extracts user ID from context.
func (s *TenantContextService) GetUserID(ctx context.Context) (string, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok || userID == "" {
		return "", fmt.Errorf("user ID not found in context")
	}
	return userID, nil
}

// GetCorrelationID extracts correlation ID from context.
func (s *TenantContextService) GetCorrelationID(ctx context.Context) (string, error) {
	correlationID, ok := ctx.Value(CorrelationIDKey).(string)
	if !ok || correlationID == "" {
		return "", fmt.Errorf("correlation ID not found in context")
	}
	return correlationID, nil
}

// GetCausationID extracts causation ID from context.
func (s *TenantContextService) GetCausationID(ctx context.Context) (string, error) {
	causationID, ok := ctx.Value(CausationIDKey).(string)
	if !ok || causationID == "" {
		return "", fmt.Errorf("causation ID not found in context")
	}
	return causationID, nil
}

// GetRequestID extracts request ID from context.
func (s *TenantContextService) GetRequestID(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok || requestID == "" {
		return "", fmt.Errorf("request ID not found in context")
	}
	return requestID, nil
}

// GetSessionID extracts session ID from context.
func (s *TenantContextService) GetSessionID(ctx context.Context) (string, error) {
	sessionID, ok := ctx.Value(SessionIDKey).(string)
	if !ok || sessionID == "" {
		return "", fmt.Errorf("session ID not found in context")
	}
	return sessionID, nil
}

// GetTenantContext extracts all tenant context from context.
func (s *TenantContextService) GetTenantContext(ctx context.Context) (TenantContext, error) {
	tenantID, _ := s.GetTenantID(ctx)
	userID, _ := s.GetUserID(ctx)
	correlationID, _ := s.GetCorrelationID(ctx)
	causationID, _ := s.GetCausationID(ctx)
	requestID, _ := s.GetRequestID(ctx)
	sessionID, _ := s.GetSessionID(ctx)
	return TenantContext{
		TenantID:      tenantID,
		UserID:        userID,
		CorrelationID: correlationID,
		CausationID:   causationID,
		RequestID:     requestID,
		SessionID:     sessionID,
	}, nil
}

// ValidateTenantContext validates the tenant context.
func (s *TenantContextService) ValidateTenantContext(tenantCtx TenantContext) error {
	if err := s.tenantValidator.ValidateTenantID(tenantCtx.TenantID); err != nil {
		return fmt.Errorf("invalid tenant ID: %w", err)
	}
	if tenantCtx.UserID != "" {
		if err := s.tenantValidator.ValidateUserID(tenantCtx.UserID); err != nil {
			return fmt.Errorf("invalid user ID: %w", err)
		}
	}
	return nil
}

// SetDefaultTenantID sets the default tenant ID.
func (s *TenantContextService) SetDefaultTenantID(tenantID string) error {
	if err := s.tenantValidator.ValidateTenantID(tenantID); err != nil {
		return fmt.Errorf("invalid default tenant ID: %w", err)
	}
	s.defaultTenantID = tenantID
	return nil
}

// GetDefaultTenantID returns the default tenant ID.
func (s *TenantContextService) GetDefaultTenantID() string {
	return s.defaultTenantID
}

// ClearTenantContext removes tenant context from the given context.
func (s *TenantContextService) ClearTenantContext(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, TenantIDKey, nil)
	ctx = context.WithValue(ctx, UserIDKey, nil)
	ctx = context.WithValue(ctx, CorrelationIDKey, nil)
	ctx = context.WithValue(ctx, CausationIDKey, nil)
	ctx = context.WithValue(ctx, RequestIDKey, nil)
	ctx = context.WithValue(ctx, SessionIDKey, nil)
	return ctx
}

// CopyTenantContext copies tenant context from one context to another.
func (s *TenantContextService) CopyTenantContext(from, to context.Context) context.Context {
	if tenantCtx, err := s.GetTenantContext(from); err == nil {
		return s.WithTenantContext(to, tenantCtx)
	}
	return to
}

// MergeTenantContext merges tenant context from one context into another.
func (s *TenantContextService) MergeTenantContext(from, to context.Context) context.Context {
	fromCtx, err := s.GetTenantContext(from)
	if err != nil {
		return to
	}
	toCtx, err := s.GetTenantContext(to)
	if err != nil {
		return s.WithTenantContext(to, fromCtx)
	}
	if fromCtx.TenantID != "" {
		toCtx.TenantID = fromCtx.TenantID
	}
	if fromCtx.UserID != "" {
		toCtx.UserID = fromCtx.UserID
	}
	if fromCtx.CorrelationID != "" {
		toCtx.CorrelationID = fromCtx.CorrelationID
	}
	if fromCtx.CausationID != "" {
		toCtx.CausationID = fromCtx.CausationID
	}
	if fromCtx.RequestID != "" {
		toCtx.RequestID = fromCtx.RequestID
	}
	if fromCtx.SessionID != "" {
		toCtx.SessionID = fromCtx.SessionID
	}
	return s.WithTenantContext(to, toCtx)
}

// IsTenantContextPresent checks if tenant context is present in the given context.
func (s *TenantContextService) IsTenantContextPresent(ctx context.Context) bool {
	_, hasTenantID := ctx.Value(TenantIDKey).(string)
	_, hasUserID := ctx.Value(UserIDKey).(string)
	return hasTenantID || hasUserID
}

// WithCorrelationID adds correlation ID to the given context.
func WithCorrelationID(ctx context.Context, correlationID string) context.Context {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	tenantCtx.CorrelationID = correlationID
	return service.WithTenantContext(ctx, tenantCtx)
}

// WithCausationID adds causation ID to the given context.
func WithCausationID(ctx context.Context, causationID string) context.Context {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	tenantCtx.CausationID = causationID
	return service.WithTenantContext(ctx, tenantCtx)
}

// GetTenantID extracts tenant ID from context (package-level function).
func GetTenantID(ctx context.Context) (string, error) {
	service := NewTenantContextService()
	return service.GetTenantID(ctx)
}

// GetCorrelationID extracts correlation ID from context (package-level function).
func GetCorrelationID(ctx context.Context) (string, error) {
	service := NewTenantContextService()
	return service.GetCorrelationID(ctx)
}

// GetCausationID extracts causation ID from context (package-level function).
func GetCausationID(ctx context.Context) (string, error) {
	service := NewTenantContextService()
	return service.GetCausationID(ctx)
}

// CorrelationIDFromContext extracts correlation ID from context. Returns (id, true) if found.
func CorrelationIDFromContext(ctx context.Context) (string, bool) {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	if tenantCtx.CorrelationID != "" {
		return tenantCtx.CorrelationID, true
	}
	return "", false
}

// CausationIDFromContext extracts causation ID from context. Returns (id, true) if found.
func CausationIDFromContext(ctx context.Context) (string, bool) {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	if tenantCtx.CausationID != "" {
		return tenantCtx.CausationID, true
	}
	return "", false
}

// EnsureCorrelationID ensures correlation_id is present in context. Returns the updated context and the correlation ID (empty if missing).
func EnsureCorrelationID(ctx context.Context) (context.Context, string) {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	correlationID := tenantCtx.CorrelationID
	if correlationID == "" {
		return ctx, ""
	}
	tenantCtx.CorrelationID = correlationID
	return service.WithTenantContext(ctx, tenantCtx), correlationID
}

// EnsureTraceIDs ensures both correlation_id and causation_id are present in context.
func EnsureTraceIDs(ctx context.Context, correlationID *string, causationID *string) context.Context {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	if correlationID != nil && *correlationID == "" && tenantCtx.CorrelationID != "" {
		*correlationID = tenantCtx.CorrelationID
	}
	if causationID != nil && *causationID == "" {
		if tenantCtx.CausationID != "" {
			*causationID = tenantCtx.CausationID
		} else if correlationID != nil && *correlationID != "" {
			*causationID = *correlationID
		}
	}
	if correlationID != nil {
		tenantCtx.CorrelationID = *correlationID
	}
	if causationID != nil {
		tenantCtx.CausationID = *causationID
	}
	return service.WithTenantContext(ctx, tenantCtx)
}

// EnsureCausationID ensures causation_id is present in context.
func EnsureCausationID(ctx context.Context, correlationID string) (context.Context, string) {
	service := NewTenantContextService()
	tenantCtx, _ := service.GetTenantContext(ctx)
	causationID := tenantCtx.CausationID
	if causationID == "" {
		if correlationID != "" {
			causationID = correlationID
		} else if tenantCtx.CorrelationID != "" {
			causationID = tenantCtx.CorrelationID
		}
	}
	tenantCtx.CausationID = causationID
	return service.WithTenantContext(ctx, tenantCtx), causationID
}
