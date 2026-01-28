package migration

import (
	"context"

	"gorm.io/gorm"
)

// GetTenantMigrations returns all tenant-specific migrations
func GetTenantMigrations() []Migration {
	return []Migration{
		{
			Version:     "2024_01_01_000001",
			Name:        "create_tenant_settings",
			Description: "Creates per-tenant settings table",
			Up: func(ctx context.Context, db *gorm.DB, tenantID uint) error {
				// This runs within the tenant's schema context
				return db.Exec(`
					INSERT INTO tenant_settings (tenant_id, setting_key, setting_value, created_at, updated_at)
					VALUES (?, 'claim_auto_approve_threshold', '0', NOW(), NOW())
					ON CONFLICT (tenant_id, setting_key) DO NOTHING
				`, tenantID).Error
			},
			Down: func(ctx context.Context, db *gorm.DB, tenantID uint) error {
				return db.Exec(`
					DELETE FROM tenant_settings
					WHERE tenant_id = ? AND setting_key = 'claim_auto_approve_threshold'
				`, tenantID).Error
			},
		},
		{
			Version:     "2024_01_01_000002",
			Name:        "initialize_tenant_sequences",
			Description: "Initializes sequence numbers for claims, packages per tenant",
			Up: func(ctx context.Context, db *gorm.DB, tenantID uint) error {
				return db.Exec(`
					INSERT INTO tenant_sequences (tenant_id, sequence_name, current_value, prefix, created_at, updated_at)
					VALUES
						(?, 'claim', 0, 'CLM', NOW(), NOW()),
						(?, 'package', 0, 'PKG', NOW(), NOW()),
						(?, 'settlement', 0, 'STL', NOW(), NOW())
					ON CONFLICT (tenant_id, sequence_name) DO NOTHING
				`, tenantID, tenantID, tenantID).Error
			},
			Down: func(ctx context.Context, db *gorm.DB, tenantID uint) error {
				return db.Exec(`
					DELETE FROM tenant_sequences WHERE tenant_id = ?
				`, tenantID).Error
			},
		},
		{
			Version:     "2024_01_02_000001",
			Name:        "add_default_reason_codes",
			Description: "Adds default deduction reason codes for tenant",
			Up: func(ctx context.Context, db *gorm.DB, tenantID uint) error {
				// Insert default reason codes if not exists
				reasonCodes := []struct {
					Code    string
					Title   string
				}{
					{"R001", "خارج از تعهد بیمه"},
					{"R002", "مدارک ناقص"},
					{"R003", "عدم تطابق با تعرفه"},
					{"R004", "تکراری بودن خدمت"},
					{"R005", "عدم پوشش بیمه‌ای"},
				}

				for _, rc := range reasonCodes {
					err := db.Exec(`
						INSERT INTO reason_codes (tenant_id, code, title_fa, is_active, created_at, updated_at)
						VALUES (?, ?, ?, true, NOW(), NOW())
						ON CONFLICT DO NOTHING
					`, tenantID, rc.Code, rc.Title).Error
					if err != nil {
						return err
					}
				}
				return nil
			},
			Down: func(ctx context.Context, db *gorm.DB, tenantID uint) error {
				return db.Exec(`
					DELETE FROM reason_codes
					WHERE tenant_id = ? AND code IN ('R001', 'R002', 'R003', 'R004', 'R005')
				`, tenantID).Error
			},
		},
	}
}

// TenantSettingsTable SQL for tenant settings table
const TenantSettingsTable = `
CREATE TABLE IF NOT EXISTS tenant_settings (
	id SERIAL PRIMARY KEY,
	tenant_id INTEGER NOT NULL,
	setting_key VARCHAR(100) NOT NULL,
	setting_value TEXT,
	setting_type VARCHAR(50) DEFAULT 'string',
	description TEXT,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(tenant_id, setting_key)
);
CREATE INDEX IF NOT EXISTS idx_tenant_settings_tenant ON tenant_settings(tenant_id);
`

// TenantSequencesTable SQL for tenant sequences table
const TenantSequencesTable = `
CREATE TABLE IF NOT EXISTS tenant_sequences (
	id SERIAL PRIMARY KEY,
	tenant_id INTEGER NOT NULL,
	sequence_name VARCHAR(100) NOT NULL,
	current_value BIGINT DEFAULT 0,
	prefix VARCHAR(20),
	suffix VARCHAR(20),
	pad_length INTEGER DEFAULT 6,
	created_at TIMESTAMP DEFAULT NOW(),
	updated_at TIMESTAMP DEFAULT NOW(),
	UNIQUE(tenant_id, sequence_name)
);
CREATE INDEX IF NOT EXISTS idx_tenant_sequences_tenant ON tenant_sequences(tenant_id);
`

// InitializeTenantTables creates base tables for tenant migrations
func InitializeTenantTables(db *gorm.DB) error {
	if err := db.Exec(TenantSettingsTable).Error; err != nil {
		return err
	}
	return db.Exec(TenantSequencesTable).Error
}
