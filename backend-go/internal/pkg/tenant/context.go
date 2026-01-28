package tenant

import (
	"context"
	"errors"
)

// Context keys
type contextKey string

const (
	tenantIDKey   contextKey = "tenant_id"
	tenantInfoKey contextKey = "tenant_info"
)

// TenantInfo contains full tenant context information
type TenantInfo struct {
	ID           uint   `json:"id"`
	Code         string `json:"code"`
	Name         string `json:"name"`
	SchemaName   string `json:"schema_name"` // For schema-based multi-tenancy
	IsActive     bool   `json:"is_active"`
	UserID       uint   `json:"user_id"`
	UserRole     string `json:"user_role"`
	CenterID     *uint  `json:"center_id,omitempty"`
	Permissions  []string `json:"permissions,omitempty"`
}

// Errors
var (
	ErrNoTenantInContext    = errors.New("no tenant found in context")
	ErrInvalidTenantID      = errors.New("invalid tenant ID")
	ErrTenantAccessDenied   = errors.New("tenant access denied")
	ErrTenantNotActive      = errors.New("tenant is not active")
	ErrCrossTenantAccess    = errors.New("cross-tenant access not allowed")
)

// WithTenantID adds tenant ID to context
func WithTenantID(ctx context.Context, tenantID uint) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

// GetTenantID extracts tenant ID from context
func GetTenantID(ctx context.Context) (uint, error) {
	value := ctx.Value(tenantIDKey)
	if value == nil {
		return 0, ErrNoTenantInContext
	}
	tenantID, ok := value.(uint)
	if !ok {
		return 0, ErrInvalidTenantID
	}
	return tenantID, nil
}

// MustGetTenantID extracts tenant ID or panics
func MustGetTenantID(ctx context.Context) uint {
	tenantID, err := GetTenantID(ctx)
	if err != nil {
		panic(err)
	}
	return tenantID
}

// WithTenantInfo adds full tenant info to context
func WithTenantInfo(ctx context.Context, info *TenantInfo) context.Context {
	ctx = WithTenantID(ctx, info.ID)
	return context.WithValue(ctx, tenantInfoKey, info)
}

// GetTenantInfo extracts full tenant info from context
func GetTenantInfo(ctx context.Context) (*TenantInfo, error) {
	value := ctx.Value(tenantInfoKey)
	if value == nil {
		// Try to create minimal info from tenant ID
		tenantID, err := GetTenantID(ctx)
		if err != nil {
			return nil, ErrNoTenantInContext
		}
		return &TenantInfo{ID: tenantID}, nil
	}
	info, ok := value.(*TenantInfo)
	if !ok {
		return nil, ErrNoTenantInContext
	}
	return info, nil
}

// HasPermission checks if tenant context has specific permission
func HasPermission(ctx context.Context, permission string) bool {
	info, err := GetTenantInfo(ctx)
	if err != nil {
		return false
	}
	for _, p := range info.Permissions {
		if p == permission || p == "*" {
			return true
		}
	}
	return false
}

// IsSystemAdmin checks if user is system admin
func IsSystemAdmin(ctx context.Context) bool {
	info, err := GetTenantInfo(ctx)
	if err != nil {
		return false
	}
	return info.UserRole == "system_admin"
}

// CanAccessTenant checks if current context can access specified tenant
func CanAccessTenant(ctx context.Context, targetTenantID uint) bool {
	info, err := GetTenantInfo(ctx)
	if err != nil {
		return false
	}
	// System admin can access any tenant
	if info.UserRole == "system_admin" {
		return true
	}
	// Others can only access their own tenant
	return info.ID == targetTenantID
}

// ValidateTenantAccess returns error if access is not allowed
func ValidateTenantAccess(ctx context.Context, targetTenantID uint) error {
	if !CanAccessTenant(ctx, targetTenantID) {
		return ErrCrossTenantAccess
	}
	return nil
}
