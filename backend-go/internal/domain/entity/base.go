package entity

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel contains common fields for all entities
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// AuditModel adds audit fields
type AuditModel struct {
	BaseModel
	CreatedBy uint `json:"created_by"`
	UpdatedBy uint `json:"updated_by"`
}

// TenantModel adds tenant isolation
type TenantModel struct {
	AuditModel
	TenantID uint `gorm:"index" json:"tenant_id"` // Insurer ID for multi-tenancy
}
