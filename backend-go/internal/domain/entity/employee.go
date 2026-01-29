package entity

import "time"

// Employee - کارمند یا فرد تحت تکفل (Simplified structure)
type Employee struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_employees_tenant,where:deleted_at IS NULL" json:"tenant_id"`

	// Parent Relationship (for family members)
	// NULL = کارمند اصلی (Main Employee)
	// NOT NULL = فرد تحت تکفل (Family Member)
	ParentID      *uint      `gorm:"index:idx_employees_parent,where:deleted_at IS NULL" json:"parent_id"`
	Parent        *Employee  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	FamilyMembers []Employee `gorm:"foreignKey:ParentID" json:"family_members,omitempty"`

	// Relation Type - نسبت با سرپرست خانوار
	RelationType *string `gorm:"size:20" json:"relation_type,omitempty"` // SELF, SPOUSE_MALE, SPOUSE_FEMALE, CHILD, etc.

	// Personal Info
	PersonnelCode string     `gorm:"size:50;index:idx_employees_personnel_code,where:deleted_at IS NULL" json:"personnel_code"`
	NationalCode  string     `gorm:"size:10;index:idx_employees_national_code,where:deleted_at IS NULL" json:"national_code"`
	FirstName     string     `gorm:"size:255;not null" json:"first_name"`
	LastName      string     `gorm:"size:255;not null" json:"last_name"`
	FatherName    *string    `gorm:"size:255" json:"father_name,omitempty"`
	BirthDate     *time.Time `json:"birth_date,omitempty"`
	Gender        string     `gorm:"size:10" json:"gender"`         // male, female
	MaritalStatus string     `gorm:"size:20" json:"marital_status"` // single, married

	// Contact
	Phone   *string `gorm:"size:20" json:"phone,omitempty"`
	Mobile  *string `gorm:"size:20" json:"mobile,omitempty"`
	Email   *string `gorm:"size:255" json:"email,omitempty"`
	Address *string `gorm:"type:text" json:"address,omitempty"`

	// Employment Info (only for main employees)
	RecruitmentDate *time.Time `json:"recruitment_date,omitempty"` // تاریخ استخدام
	RetirementDate  *time.Time `json:"retirement_date,omitempty"`  // تاریخ بازنشستگی

	// Status
	IsActive bool   `gorm:"default:true" json:"is_active"`
	Status   string `gorm:"size:20;default:active" json:"status"` // active, inactive, retired

	// Relations
	Tenant *Insurer `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
}

// TableName specifies the table name
func (Employee) TableName() string {
	return "employees"
}

// IsMainEmployee checks if this is a main employee (not a family member)
func (e *Employee) IsMainEmployee() bool {
	return e.ParentID == nil
}

// IsFamilyMember checks if this is a family member
func (e *Employee) IsFamilyMember() bool {
	return e.ParentID != nil
}

// GetFullName returns full name
func (e *Employee) GetFullName() string {
	return e.FirstName + " " + e.LastName
}
