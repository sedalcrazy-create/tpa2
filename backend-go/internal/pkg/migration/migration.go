package migration

import (
	"context"
	"fmt"
	"sort"
	"time"

	"gorm.io/gorm"
)

// MigrationStatus represents migration execution status
type MigrationStatus string

const (
	MigrationPending   MigrationStatus = "pending"
	MigrationRunning   MigrationStatus = "running"
	MigrationCompleted MigrationStatus = "completed"
	MigrationFailed    MigrationStatus = "failed"
	MigrationRolledBack MigrationStatus = "rolled_back"
)

// MigrationRecord tracks migration history per tenant
type MigrationRecord struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	TenantID    uint            `gorm:"index" json:"tenant_id"`
	Version     string          `gorm:"size:100;index" json:"version"`
	Name        string          `gorm:"size:255" json:"name"`
	Description string          `gorm:"size:1000" json:"description"`
	Status      MigrationStatus `gorm:"size:50" json:"status"`
	StartedAt   time.Time       `json:"started_at"`
	CompletedAt *time.Time      `json:"completed_at"`
	Error       string          `gorm:"size:2000" json:"error,omitempty"`
	Checksum    string          `gorm:"size:64" json:"checksum"`
	ExecutedBy  string          `gorm:"size:100" json:"executed_by"`
	BatchNo     int             `json:"batch_no"`
}

// TableName returns table name for GORM
func (MigrationRecord) TableName() string {
	return "tenant_migrations"
}

// Migration defines a single migration
type Migration struct {
	Version     string
	Name        string
	Description string
	Up          func(ctx context.Context, db *gorm.DB, tenantID uint) error
	Down        func(ctx context.Context, db *gorm.DB, tenantID uint) error
}

// Migrator handles per-tenant migrations
type Migrator struct {
	db         *gorm.DB
	migrations []Migration
}

// NewMigrator creates a new migrator
func NewMigrator(db *gorm.DB) *Migrator {
	return &Migrator{
		db:         db,
		migrations: make([]Migration, 0),
	}
}

// Register adds a migration to the migrator
func (m *Migrator) Register(migration Migration) {
	m.migrations = append(m.migrations, migration)
}

// RegisterMany adds multiple migrations
func (m *Migrator) RegisterMany(migrations []Migration) {
	m.migrations = append(m.migrations, migrations...)
}

// Initialize creates the migrations tracking table
func (m *Migrator) Initialize(ctx context.Context) error {
	return m.db.WithContext(ctx).AutoMigrate(&MigrationRecord{})
}

// MigrateAll runs all pending migrations for all tenants
func (m *Migrator) MigrateAll(ctx context.Context) error {
	// Get all tenant IDs
	var tenantIDs []uint
	err := m.db.WithContext(ctx).
		Raw("SELECT DISTINCT tenant_id FROM insurers WHERE deleted_at IS NULL").
		Pluck("tenant_id", &tenantIDs).Error
	if err != nil {
		// If insurers table doesn't exist, use default tenant
		tenantIDs = []uint{1}
	}

	for _, tenantID := range tenantIDs {
		if err := m.Migrate(ctx, tenantID); err != nil {
			return fmt.Errorf("migration failed for tenant %d: %w", tenantID, err)
		}
	}
	return nil
}

// Migrate runs all pending migrations for a specific tenant
func (m *Migrator) Migrate(ctx context.Context, tenantID uint) error {
	// Sort migrations by version
	sort.Slice(m.migrations, func(i, j int) bool {
		return m.migrations[i].Version < m.migrations[j].Version
	})

	// Get last batch number
	var lastBatch int
	m.db.WithContext(ctx).
		Model(&MigrationRecord{}).
		Where("tenant_id = ?", tenantID).
		Select("COALESCE(MAX(batch_no), 0)").
		Scan(&lastBatch)
	newBatch := lastBatch + 1

	// Run each migration
	for _, migration := range m.migrations {
		if m.hasRun(ctx, tenantID, migration.Version) {
			continue
		}

		if err := m.runMigration(ctx, tenantID, migration, newBatch); err != nil {
			return err
		}
	}

	return nil
}

// hasRun checks if a migration has already been executed for a tenant
func (m *Migrator) hasRun(ctx context.Context, tenantID uint, version string) bool {
	var count int64
	m.db.WithContext(ctx).
		Model(&MigrationRecord{}).
		Where("tenant_id = ? AND version = ? AND status = ?", tenantID, version, MigrationCompleted).
		Count(&count)
	return count > 0
}

// runMigration executes a single migration
func (m *Migrator) runMigration(ctx context.Context, tenantID uint, migration Migration, batchNo int) error {
	// Create record
	record := &MigrationRecord{
		TenantID:    tenantID,
		Version:     migration.Version,
		Name:        migration.Name,
		Description: migration.Description,
		Status:      MigrationRunning,
		StartedAt:   time.Now(),
		BatchNo:     batchNo,
		ExecutedBy:  "system",
	}
	if err := m.db.WithContext(ctx).Create(record).Error; err != nil {
		return err
	}

	// Run migration in transaction
	err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return migration.Up(ctx, tx, tenantID)
	})

	// Update record
	now := time.Now()
	if err != nil {
		record.Status = MigrationFailed
		record.Error = err.Error()
	} else {
		record.Status = MigrationCompleted
	}
	record.CompletedAt = &now
	m.db.WithContext(ctx).Save(record)

	return err
}

// Rollback rolls back the last batch of migrations for a tenant
func (m *Migrator) Rollback(ctx context.Context, tenantID uint) error {
	// Get last batch
	var lastBatch int
	m.db.WithContext(ctx).
		Model(&MigrationRecord{}).
		Where("tenant_id = ? AND status = ?", tenantID, MigrationCompleted).
		Select("COALESCE(MAX(batch_no), 0)").
		Scan(&lastBatch)

	if lastBatch == 0 {
		return nil // Nothing to rollback
	}

	// Get migrations in this batch (reverse order)
	var records []MigrationRecord
	m.db.WithContext(ctx).
		Where("tenant_id = ? AND batch_no = ? AND status = ?", tenantID, lastBatch, MigrationCompleted).
		Order("version DESC").
		Find(&records)

	// Rollback each
	for _, record := range records {
		// Find migration definition
		var migration *Migration
		for _, mig := range m.migrations {
			if mig.Version == record.Version {
				migration = &mig
				break
			}
		}

		if migration == nil || migration.Down == nil {
			continue // Skip if no down migration
		}

		// Run down migration
		err := m.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return migration.Down(ctx, tx, tenantID)
		})

		// Update record
		if err != nil {
			record.Status = MigrationFailed
			record.Error = "rollback failed: " + err.Error()
		} else {
			record.Status = MigrationRolledBack
		}
		now := time.Now()
		record.CompletedAt = &now
		m.db.WithContext(ctx).Save(&record)

		if err != nil {
			return err
		}
	}

	return nil
}

// Status returns migration status for a tenant
func (m *Migrator) Status(ctx context.Context, tenantID uint) ([]MigrationRecord, error) {
	var records []MigrationRecord
	err := m.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("version ASC").
		Find(&records).Error
	return records, err
}

// Pending returns pending migrations for a tenant
func (m *Migrator) Pending(ctx context.Context, tenantID uint) []Migration {
	var pending []Migration
	for _, migration := range m.migrations {
		if !m.hasRun(ctx, tenantID, migration.Version) {
			pending = append(pending, migration)
		}
	}
	return pending
}
