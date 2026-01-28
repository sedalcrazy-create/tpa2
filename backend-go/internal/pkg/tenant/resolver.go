package tenant

import (
	"strconv"
	"sync"

	"github.com/gofiber/fiber/v2"
)

// ResolverConfig configures the tenant resolver
type ResolverConfig struct {
	// HeaderName is the header to check for tenant ID (default: X-Tenant-ID)
	HeaderName string

	// QueryParam is the query parameter to check for tenant ID (default: tenant_id)
	QueryParam string

	// DefaultTenantID is used when no tenant is specified (0 means required)
	DefaultTenantID uint

	// AllowCrossTenant allows system admins to access other tenants
	AllowCrossTenant bool

	// TenantLoader loads full tenant info from database
	TenantLoader TenantLoader

	// Cache configuration
	CacheEnabled bool
	CacheTTL     int // seconds
}

// TenantLoader is a function that loads tenant info from database
type TenantLoader func(tenantID uint) (*TenantInfo, error)

// DefaultConfig returns default resolver configuration
func DefaultConfig() ResolverConfig {
	return ResolverConfig{
		HeaderName:       "X-Tenant-ID",
		QueryParam:       "tenant_id",
		DefaultTenantID:  0,
		AllowCrossTenant: true,
		CacheEnabled:     true,
		CacheTTL:         300, // 5 minutes
	}
}

// Resolver handles tenant resolution from requests
type Resolver struct {
	config ResolverConfig
	cache  sync.Map // Simple in-memory cache
}

// NewResolver creates a new tenant resolver
func NewResolver(config ...ResolverConfig) *Resolver {
	cfg := DefaultConfig()
	if len(config) > 0 {
		cfg = config[0]
	}
	return &Resolver{
		config: cfg,
	}
}

// Middleware returns Fiber middleware for tenant resolution
func (r *Resolver) Middleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Try to get tenant ID from JWT (already set by auth middleware)
		jwtTenantID := c.Locals("tenant_id")

		// 2. Try to get tenant ID from header
		var headerTenantID uint
		headerValue := c.Get(r.config.HeaderName)
		if headerValue != "" {
			if id, err := strconv.ParseUint(headerValue, 10, 32); err == nil {
				headerTenantID = uint(id)
			}
		}

		// 3. Try to get tenant ID from query param
		var queryTenantID uint
		queryValue := c.Query(r.config.QueryParam)
		if queryValue != "" {
			if id, err := strconv.ParseUint(queryValue, 10, 32); err == nil {
				queryTenantID = uint(id)
			}
		}

		// 4. Determine final tenant ID (priority: header > query > JWT > default)
		var tenantID uint
		var requestedDifferentTenant bool

		if jwtTenantID != nil {
			tenantID = jwtTenantID.(uint)
		}

		if headerTenantID > 0 {
			if tenantID > 0 && headerTenantID != tenantID {
				requestedDifferentTenant = true
			}
			tenantID = headerTenantID
		} else if queryTenantID > 0 {
			if tenantID > 0 && queryTenantID != tenantID {
				requestedDifferentTenant = true
			}
			tenantID = queryTenantID
		}

		// Use default if still not set
		if tenantID == 0 {
			tenantID = r.config.DefaultTenantID
		}

		// 5. Validate tenant ID
		if tenantID == 0 {
			return fiber.NewError(fiber.StatusBadRequest, "tenant ID is required")
		}

		// 6. Check cross-tenant access
		if requestedDifferentTenant && r.config.AllowCrossTenant {
			roleName := c.Locals("role_name")
			if roleName == nil || roleName.(string) != "system_admin" {
				return fiber.NewError(fiber.StatusForbidden, "cross-tenant access denied")
			}
		}

		// 7. Build tenant info
		info := &TenantInfo{
			ID: tenantID,
		}

		// Add user info from JWT
		if userID := c.Locals("user_id"); userID != nil {
			info.UserID = userID.(uint)
		}
		if roleName := c.Locals("role_name"); roleName != nil {
			info.UserRole = roleName.(string)
		}
		if centerID := c.Locals("center_id"); centerID != nil {
			info.CenterID = centerID.(*uint)
		}

		// 8. Load full tenant info if loader is configured
		if r.config.TenantLoader != nil {
			fullInfo, err := r.loadTenantInfo(tenantID)
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "failed to load tenant info")
			}
			if fullInfo != nil {
				if !fullInfo.IsActive {
					return fiber.NewError(fiber.StatusForbidden, "tenant is not active")
				}
				// Merge with user info
				fullInfo.UserID = info.UserID
				fullInfo.UserRole = info.UserRole
				fullInfo.CenterID = info.CenterID
				info = fullInfo
			}
		}

		// 9. Set tenant info in context
		c.Locals("tenant_info", info)
		c.Locals("tenant_id", tenantID)

		// 10. Set context for downstream handlers
		ctx := WithTenantInfo(c.UserContext(), info)
		c.SetUserContext(ctx)

		return c.Next()
	}
}

// loadTenantInfo loads tenant info with caching
func (r *Resolver) loadTenantInfo(tenantID uint) (*TenantInfo, error) {
	// Check cache first
	if r.config.CacheEnabled {
		if cached, ok := r.cache.Load(tenantID); ok {
			return cached.(*TenantInfo), nil
		}
	}

	// Load from database
	if r.config.TenantLoader != nil {
		info, err := r.config.TenantLoader(tenantID)
		if err != nil {
			return nil, err
		}
		// Cache the result
		if r.config.CacheEnabled && info != nil {
			r.cache.Store(tenantID, info)
		}
		return info, nil
	}

	return nil, nil
}

// InvalidateCache removes tenant from cache
func (r *Resolver) InvalidateCache(tenantID uint) {
	r.cache.Delete(tenantID)
}

// ClearCache removes all cached tenants
func (r *Resolver) ClearCache() {
	r.cache = sync.Map{}
}

// GetTenantFromContext helper to get tenant info from Fiber context
func GetTenantFromContext(c *fiber.Ctx) (*TenantInfo, error) {
	info := c.Locals("tenant_info")
	if info == nil {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return nil, ErrNoTenantInContext
		}
		return &TenantInfo{ID: tenantID.(uint)}, nil
	}
	return info.(*TenantInfo), nil
}

// GetTenantIDFromContext helper to get only tenant ID from Fiber context
func GetTenantIDFromContext(c *fiber.Ctx) (uint, error) {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return 0, ErrNoTenantInContext
	}
	return tenantID.(uint), nil
}

// MustGetTenantIDFromContext helper that panics if no tenant
func MustGetTenantIDFromContext(c *fiber.Ctx) uint {
	tenantID, err := GetTenantIDFromContext(c)
	if err != nil {
		panic(err)
	}
	return tenantID
}
