package tenant

import (
	"context"

	"gorm.io/gorm"
)

// ScopedDB returns a GORM DB scoped to the tenant in context
func ScopedDB(ctx context.Context, db *gorm.DB) *gorm.DB {
	tenantID, err := GetTenantID(ctx)
	if err != nil {
		// Return unscoped if no tenant (will fail on query)
		return db.WithContext(ctx)
	}
	return db.WithContext(ctx).Where("tenant_id = ?", tenantID)
}

// ScopedDBWithID returns a GORM DB scoped to specific tenant ID
func ScopedDBWithID(ctx context.Context, db *gorm.DB, tenantID uint) *gorm.DB {
	return db.WithContext(ctx).Where("tenant_id = ?", tenantID)
}

// TenantScope is a GORM scope function for tenant filtering
func TenantScope(tenantID uint) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("tenant_id = ?", tenantID)
	}
}

// TenantScopeFromContext creates a scope from context
func TenantScopeFromContext(ctx context.Context) func(*gorm.DB) *gorm.DB {
	tenantID, err := GetTenantID(ctx)
	if err != nil {
		return func(db *gorm.DB) *gorm.DB { return db }
	}
	return TenantScope(tenantID)
}

// AutoSetTenant is a GORM callback to automatically set tenant_id on create
func AutoSetTenant(ctx context.Context) func(*gorm.DB) {
	return func(db *gorm.DB) {
		tenantID, err := GetTenantID(ctx)
		if err != nil {
			return
		}
		db.Statement.SetColumn("tenant_id", tenantID)
	}
}

// RegisterCallbacks registers GORM callbacks for automatic tenant handling
func RegisterCallbacks(db *gorm.DB) error {
	// Before create: set tenant_id from context
	err := db.Callback().Create().Before("gorm:create").Register("tenant:before_create", func(db *gorm.DB) {
		ctx := db.Statement.Context
		if ctx == nil {
			return
		}
		tenantID, err := GetTenantID(ctx)
		if err != nil {
			return
		}
		// Only set if model has tenant_id field and it's not already set
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("tenant_id"); field != nil {
				if val, _ := field.ValueOf(ctx, db.Statement.ReflectValue); val == uint(0) {
					db.Statement.SetColumn("TenantID", tenantID)
				}
			}
		}
	})
	if err != nil {
		return err
	}

	// Before query: add tenant filter
	err = db.Callback().Query().Before("gorm:query").Register("tenant:before_query", func(db *gorm.DB) {
		ctx := db.Statement.Context
		if ctx == nil {
			return
		}
		// Skip if explicitly disabled
		if skip, ok := ctx.Value("skip_tenant_scope").(bool); ok && skip {
			return
		}
		// Check if model has tenant_id field
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("tenant_id"); field != nil {
				tenantID, err := GetTenantID(ctx)
				if err == nil && tenantID > 0 {
					db.Where("tenant_id = ?", tenantID)
				}
			}
		}
	})
	if err != nil {
		return err
	}

	// Before update: ensure tenant match
	err = db.Callback().Update().Before("gorm:update").Register("tenant:before_update", func(db *gorm.DB) {
		ctx := db.Statement.Context
		if ctx == nil {
			return
		}
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("tenant_id"); field != nil {
				tenantID, err := GetTenantID(ctx)
				if err == nil && tenantID > 0 {
					db.Where("tenant_id = ?", tenantID)
				}
			}
		}
	})
	if err != nil {
		return err
	}

	// Before delete: ensure tenant match
	err = db.Callback().Delete().Before("gorm:delete").Register("tenant:before_delete", func(db *gorm.DB) {
		ctx := db.Statement.Context
		if ctx == nil {
			return
		}
		if db.Statement.Schema != nil {
			if field := db.Statement.Schema.LookUpField("tenant_id"); field != nil {
				tenantID, err := GetTenantID(ctx)
				if err == nil && tenantID > 0 {
					db.Where("tenant_id = ?", tenantID)
				}
			}
		}
	})

	return err
}

// SkipTenantScope returns context that skips tenant scope (for admin operations)
func SkipTenantScope(ctx context.Context) context.Context {
	return context.WithValue(ctx, "skip_tenant_scope", true)
}
