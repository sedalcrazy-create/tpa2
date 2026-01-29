package entity

import "time"

// Employee - کارمند یا فرد تحت تکفل (Compatible with Refah/Yii structure)
type Employee struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_employees_tenant,where:deleted_at IS NULL"`

	// Parent Relationship (for family members)
	// NULL = کارمند اصلی (Main Employee)
	// NOT NULL = فرد تحت تکفل (Family Member)
	ParentID       *uint `gorm:"index:idx_employees_parent,where:deleted_at IS NULL"`
	Parent         *Employee
	FamilyMembers  []Employee `gorm:"foreignKey:ParentID"`
	RelationTypeID *uint      `gorm:"index:idx_employees_relation_type,where:deleted_at IS NULL"`
	RelationType   *RelationType

	// Codes & Types
	CustomEmployeeCodeID    *uint `json:"custom_employee_code_id,omitempty"`
	CustomEmployeeCode      *CustomEmployeeCode
	SpecialEmployeeTypeID   *uint `json:"special_employee_type_id,omitempty"`
	SpecialEmployeeType     *SpecialEmployeeType
	GuardianshipTypeID      *uint
	GuardianshipType        *GuardianshipType

	// Personal Info
	PersonnelCode  string  `gorm:"size:50;index:idx_employees_personnel_code,where:deleted_at IS NULL"`
	NationalCode   string  `gorm:"size:10;index:idx_employees_national_code,where:deleted_at IS NULL"`
	FirstName      string  `gorm:"size:255;not null"`
	LastName       string  `gorm:"size:255;not null"`
	FatherName     *string `gorm:"size:255"`
	BirthDate      *time.Time
	Gender         string `gorm:"size:10"` // male, female
	MaritalStatus  string `gorm:"size:20"` // single, married
	IDNumber       *string `gorm:"size:50"` // شماره شناسنامه

	// Contact
	Phone   *string `gorm:"size:20"`
	Mobile  *string `gorm:"size:20"`
	Email   *string `gorm:"size:255"`
	Address *string `gorm:"type:text"`

	// Employment Info
	BranchID         *int
	LocationID       *int // محل سکونت
	WorkLocationID   *int // محل خدمت
	AccountNumber    *string `gorm:"size:50"`
	RecruitmentDate  *time.Time // تاریخ استخدام
	TerminationDate  *time.Time // تاریخ خروج
	RetirementDate   *time.Time // تاریخ بازنشستگی

	// Status
	IsActive bool   `gorm:"default:true"`
	Priority int    `gorm:"default:1"`
	Status   string `gorm:"size:20;default:active"` // active, inactive, retired

	// Extra
	Picture     *string `gorm:"size:255"`
	Description *string `gorm:"type:text"`

	// Relations
	Tenant *Insurer `gorm:"foreignKey:TenantID"`
}

// TableName specifies the table name
func (Employee) TableName() string {
	return "employees"
}

// IsMainEmployee checks if this is a main employee (not a family member)
func (e *Employee) IsMainEmployee() bool {
	return e.ParentID == nil || (e.RelationTypeID != nil && e.RelationType != nil && e.RelationType.IsMainEmployee())
}

// IsFamilyMember checks if this is a family member
func (e *Employee) IsFamilyMember() bool {
	return !e.IsMainEmployee()
}

// GetFullName returns full name
func (e *Employee) GetFullName() string {
	return e.FirstName + " " + e.LastName
}

// GenerateEmployeeTypeCode generates type code based on Refah/Yii logic
// Formula from Employee.php:
// empTypeCode = (id_set * 1000) + (isRetired ? 100 : 200) + id_cec
func (e *Employee) GenerateEmployeeTypeCode() *int {
	if e.CustomEmployeeCodeID == nil {
		return nil
	}

	code := 0

	// Add special employee type code (id_set * 1000)
	if e.SpecialEmployeeTypeID != nil {
		code += int(*e.SpecialEmployeeTypeID) * 1000
	}

	// Check if retired from custom employee code
	if e.CustomEmployeeCode != nil && e.CustomEmployeeCode.IsRetired {
		code += 100
	} else {
		code += 200
	}

	// Add custom employee code ID
	code += int(*e.CustomEmployeeCodeID)

	return &code
}

// CalculateNewPCode calculates employee code from stored procedure logic
// Main Employee: festno
// Family Member: (9000000 + parent_festno) * 100 + child_number
func CalculateNewPCode(personnelCode string, parentPersonnelCode *string, childNumber int) int64 {
	// TODO: Implement based on stored procedure logic
	// This is a placeholder
	return 0
}
