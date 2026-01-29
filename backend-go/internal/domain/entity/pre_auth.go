package entity

import "time"

// PreAuth represents pre-authorization (علی‌الحساب) - advance payments
type PreAuth struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_pre_auths_tenant,where:deleted_at IS NULL" json:"tenant_id"`

	PersonID uint   `gorm:"not null;index:idx_pre_auths_person,where:deleted_at IS NULL" json:"person_id"`
	Subject  string `gorm:"size:500" json:"subject,omitempty"`
	Amount   int64  `gorm:"not null" json:"amount"`
	Type     *int16 `json:"type,omitempty"`

	PaymentDate *time.Time `json:"payment_date,omitempty"`

	RegisterUserID uint      `gorm:"not null" json:"register_user_id"`
	RegisterDate   time.Time `gorm:"not null" json:"register_date"`

	ClaimID *uint `gorm:"index:idx_pre_auths_claim,where:deleted_at IS NULL" json:"claim_id,omitempty"`

	// Relations
	Tenant       *Insurer `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Person       *Person  `gorm:"foreignKey:PersonID" json:"person,omitempty"`
	RegisterUser *User    `gorm:"foreignKey:RegisterUserID" json:"register_user,omitempty"`
	// Claim relation loaded manually to avoid circular dependency (Claim also references PreAuth)
	Claim        *Claim   `gorm:"-" json:"claim,omitempty"`
}

// TableName specifies the table name for PreAuth
func (PreAuth) TableName() string {
	return "pre_auths"
}
