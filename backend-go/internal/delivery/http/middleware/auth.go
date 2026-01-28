package middleware

import (
	"strings"
	"time"

	"github.com/bank-melli/tpa/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents JWT token claims
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	RoleID   uint   `json:"role_id"`
	RoleName string `json:"role_name"`
	TenantID uint   `json:"tenant_id"`
	CenterID *uint  `json:"center_id,omitempty"`
	jwt.RegisteredClaims
}

// AuthMiddleware creates authentication middleware
func AuthMiddleware(cfg *config.JWTConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "missing authorization header")
		}

		// Extract token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid authorization header format")
		}
		tokenString := parts[1]

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
			}
			return []byte(cfg.Secret), nil
		})

		if err != nil {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token: "+err.Error())
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			return fiber.NewError(fiber.StatusUnauthorized, "invalid token claims")
		}

		// Check expiration
		if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
			return fiber.NewError(fiber.StatusUnauthorized, "token expired")
		}

		// Set user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role_id", claims.RoleID)
		c.Locals("role_name", claims.RoleName)
		c.Locals("tenant_id", claims.TenantID)
		c.Locals("center_id", claims.CenterID)

		return c.Next()
	}
}

// TenantMiddleware extracts tenant ID from header or JWT
// Deprecated: Use tenant.NewResolver().Middleware() instead
func TenantMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// First check X-Tenant-ID header
		tenantIDHeader := c.Get("X-Tenant-ID")
		if tenantIDHeader != "" {
			var tenantID uint
			if _, err := c.ParamsInt("tenant_id"); err == nil {
				// Already set from JWT, verify it matches
				jwtTenantID := c.Locals("tenant_id")
				if jwtTenantID != nil && jwtTenantID.(uint) != tenantID {
					// Allow only system admin to switch tenants
					roleName := c.Locals("role_name")
					if roleName == nil || roleName.(string) != "system_admin" {
						return fiber.NewError(fiber.StatusForbidden, "cannot access other tenant's data")
					}
				}
			}
		}

		// Use tenant ID from JWT if not set
		if c.Locals("tenant_id") == nil {
			c.Locals("tenant_id", uint(1)) // Default tenant
		}

		return c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		roleName := c.Locals("role_name")
		if roleName == nil {
			return fiber.NewError(fiber.StatusForbidden, "access denied")
		}

		userRole := roleName.(string)
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return fiber.NewError(fiber.StatusForbidden, "access denied")
	}
}

// PermissionMiddleware checks if user has required permission
func PermissionMiddleware(requiredPermission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO: Implement permission check from database or cache
		// For now, allow all authenticated users
		return c.Next()
	}
}

// CenterAccessMiddleware restricts access to center's own data
func CenterAccessMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userCenterID := c.Locals("center_id")
		roleName := c.Locals("role_name")

		// System admin and insurer admin can access all centers
		if roleName != nil {
			role := roleName.(string)
			if role == "system_admin" || role == "insurer_admin" || role == "supervisor" {
				return c.Next()
			}
		}

		// For center users, restrict to their own center
		if userCenterID != nil {
			requestCenterID := c.Query("center_id")
			if requestCenterID != "" {
				// Validate center access
				// TODO: Parse and compare center IDs
			}
		}

		return c.Next()
	}
}

// GenerateToken generates a new JWT token
func GenerateToken(claims *JWTClaims, secret string, expiresIn time.Duration) (string, error) {
	claims.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(userID uint, secret string, expiresIn time.Duration) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   string(rune(userID)),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
