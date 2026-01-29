package entity

import "time"

// EmployeeImportTemp - جدول موقت برای import کارمندان از سرور HR
type EmployeeImportTemp struct {
	ID       uint `gorm:"primaryKey"`
	TenantID uint `gorm:"not null;index:idx_employees_temp_tenant"`

	// Same structure as Employee
	ParentID                *uint
	RelationTypeID          *uint
	CustomEmployeeCodeID    *uint
	SpecialEmployeeTypeID   *uint
	GuardianshipTypeID      *uint

	PersonnelCode   string
	NationalCode    string
	FirstName       string
	LastName        string
	FatherName      *string
	BirthDate       *time.Time
	Gender          string
	MaritalStatus   string
	IDNumber        *string

	Phone   *string
	Mobile  *string
	Email   *string
	Address *string

	BranchID        *int
	LocationID      *int
	WorkLocationID  *int
	AccountNumber   *string
	RecruitmentDate *time.Time
	TerminationDate *time.Time
	RetirementDate  *time.Time

	IsActive bool
	Priority int
	Status   string
	Picture  *string
	Description *string

	// Import metadata
	ImportBatchID string    `gorm:"size:100;index:idx_employees_temp_batch"`
	ImportDate    time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// TableName specifies the table name
func (EmployeeImportTemp) TableName() string {
	return "employees_import_temp"
}

// EmployeeImportHistory - تاریخچه import کارمندان
type EmployeeImportHistory struct {
	ID       uint `gorm:"primaryKey"`
	TenantID uint `gorm:"not null"`

	BatchID        string    `gorm:"size:100;uniqueIndex;not null"`
	ImportDate     time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	Source         string    `gorm:"size:100"` // hr_server, csv_file, manual
	TotalRecords   int       `gorm:"default:0"`
	NewRecords     int       `gorm:"default:0"`
	UpdatedRecords int       `gorm:"default:0"`
	FailedRecords  int       `gorm:"default:0"`
	Status         string    `gorm:"size:50;default:pending"` // pending, processing, completed, failed
	Notes          *string   `gorm:"type:text"`
	ImportedByUserID *uint

	CreatedAt time.Time
}

// TableName specifies the table name
func (EmployeeImportHistory) TableName() string {
	return "employee_import_history"
}

// ImportStatus constants
const (
	ImportStatusPending    = "pending"
	ImportStatusProcessing = "processing"
	ImportStatusCompleted  = "completed"
	ImportStatusFailed     = "failed"
)
