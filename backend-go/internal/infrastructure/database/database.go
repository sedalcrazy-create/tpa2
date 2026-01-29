package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bank-melli/tpa/internal/config"
	"github.com/bank-melli/tpa/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database wraps the GORM database connection
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database connection
func NewDatabase(cfg *config.DatabaseConfig) (*Database, error) {
	// Configure GORM logger
	logLevel := logger.Silent
	if os.Getenv("APP_DEBUG") == "true" {
		logLevel = logger.Info
	}

	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	// Connect to database
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger:                 gormLogger,
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying SQL database
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

	return &Database{db}, nil
}

// AutoMigrate runs auto migration for all entities
func (db *Database) AutoMigrate() error {
	// SIMPLIFIED: Only migrate essential tables for testing
	return db.DB.AutoMigrate(
		// Base entities
		&entity.Insurer{},

		// Employees (single table for main employees and family members)
		&entity.Employee{},

		// Users & Auth
		&entity.Role{},
		&entity.Permission{},
		&entity.User{},
	)
}

// Close closes the database connection
func (db *Database) Close() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// HealthCheck checks if database connection is healthy
func (db *Database) HealthCheck() error {
	sqlDB, err := db.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// WithTenant returns a new DB instance scoped to a tenant
func (db *Database) WithTenant(tenantID uint) *gorm.DB {
	return db.DB.Where("tenant_id = ?", tenantID)
}

// Transaction executes a function within a transaction
func (db *Database) Transaction(fn func(tx *gorm.DB) error) error {
	return db.DB.Transaction(fn)
}

// Seed seeds initial data
func (db *Database) Seed() error {
	// Seed roles
	roles := []entity.Role{
		{Name: entity.RoleSystemAdmin, TitleFa: "مدیر سیستم", Level: 100, IsSystem: true, IsActive: true},
		{Name: entity.RoleInsurerAdmin, TitleFa: "مدیر بیمه‌گر", Level: 90, IsSystem: true, IsActive: true},
		{Name: entity.RoleSupervisor, TitleFa: "سرپرست", Level: 80, IsSystem: true, IsActive: true},
		{Name: entity.RoleClaimExaminer, TitleFa: "ارزیاب", Level: 50, IsSystem: false, IsActive: true},
		{Name: entity.RoleDrugExaminer, TitleFa: "ارزیاب دارو", Level: 50, IsSystem: false, IsActive: true},
		{Name: entity.RoleFinancialOfficer, TitleFa: "کارشناس مالی", Level: 60, IsSystem: false, IsActive: true},
		{Name: entity.RoleCenterUser, TitleFa: "کاربر مرکز", Level: 30, IsSystem: false, IsActive: true},
		{Name: entity.RoleReportViewer, TitleFa: "مشاهده‌کننده گزارش", Level: 20, IsSystem: false, IsActive: true},
	}

	for _, role := range roles {
		var existing entity.Role
		if err := db.DB.Where("name = ?", role.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.DB.Create(&role).Error; err != nil {
					return fmt.Errorf("failed to seed role %s: %w", role.Name, err)
				}
			} else {
				return err
			}
		}
	}

	// Seed permissions
	permissions := []entity.Permission{
		// Claims
		{Name: entity.PermClaimCreate, TitleFa: "ایجاد ادعا", Module: "claims", IsActive: true},
		{Name: entity.PermClaimRead, TitleFa: "مشاهده ادعا", Module: "claims", IsActive: true},
		{Name: entity.PermClaimUpdate, TitleFa: "ویرایش ادعا", Module: "claims", IsActive: true},
		{Name: entity.PermClaimDelete, TitleFa: "حذف ادعا", Module: "claims", IsActive: true},
		{Name: entity.PermClaimExamine, TitleFa: "ارزیابی ادعا", Module: "claims", IsActive: true},
		{Name: entity.PermClaimApprove, TitleFa: "تایید ادعا", Module: "claims", IsActive: true},
		{Name: entity.PermClaimReject, TitleFa: "رد ادعا", Module: "claims", IsActive: true},

		// Packages
		{Name: entity.PermPackageCreate, TitleFa: "ایجاد بسته", Module: "packages", IsActive: true},
		{Name: entity.PermPackageRead, TitleFa: "مشاهده بسته", Module: "packages", IsActive: true},
		{Name: entity.PermPackageUpdate, TitleFa: "ویرایش بسته", Module: "packages", IsActive: true},
		{Name: entity.PermPackageDelete, TitleFa: "حذف بسته", Module: "packages", IsActive: true},
		{Name: entity.PermPackageExamine, TitleFa: "ارزیابی بسته", Module: "packages", IsActive: true},
		{Name: entity.PermPackageApprove, TitleFa: "تایید بسته", Module: "packages", IsActive: true},

		// Centers
		{Name: entity.PermCenterCreate, TitleFa: "ایجاد مرکز", Module: "centers", IsActive: true},
		{Name: entity.PermCenterRead, TitleFa: "مشاهده مرکز", Module: "centers", IsActive: true},
		{Name: entity.PermCenterUpdate, TitleFa: "ویرایش مرکز", Module: "centers", IsActive: true},
		{Name: entity.PermCenterDelete, TitleFa: "حذف مرکز", Module: "centers", IsActive: true},

		// Settlements
		{Name: entity.PermSettlementCreate, TitleFa: "ایجاد تسویه", Module: "settlements", IsActive: true},
		{Name: entity.PermSettlementRead, TitleFa: "مشاهده تسویه", Module: "settlements", IsActive: true},
		{Name: entity.PermSettlementApprove, TitleFa: "تایید تسویه", Module: "settlements", IsActive: true},

		// Users
		{Name: entity.PermUserCreate, TitleFa: "ایجاد کاربر", Module: "users", IsActive: true},
		{Name: entity.PermUserRead, TitleFa: "مشاهده کاربر", Module: "users", IsActive: true},
		{Name: entity.PermUserUpdate, TitleFa: "ویرایش کاربر", Module: "users", IsActive: true},
		{Name: entity.PermUserDelete, TitleFa: "حذف کاربر", Module: "users", IsActive: true},

		// Reports
		{Name: entity.PermReportView, TitleFa: "مشاهده گزارش", Module: "reports", IsActive: true},
		{Name: entity.PermReportExport, TitleFa: "خروجی گزارش", Module: "reports", IsActive: true},

		// Settings
		{Name: entity.PermSettingsRead, TitleFa: "مشاهده تنظیمات", Module: "settings", IsActive: true},
		{Name: entity.PermSettingsUpdate, TitleFa: "ویرایش تنظیمات", Module: "settings", IsActive: true},
	}

	for _, perm := range permissions {
		var existing entity.Permission
		if err := db.DB.Where("name = ?", perm.Name).First(&existing).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.DB.Create(&perm).Error; err != nil {
					return fmt.Errorf("failed to seed permission %s: %w", perm.Name, err)
				}
			} else {
				return err
			}
		}
	}

	return nil
}
