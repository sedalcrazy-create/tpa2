package entity

import "time"

// Prescription represents a medical prescription (نسخه پزشکی)
// This is separate from Claim - prescriptions can later be converted to claims
type Prescription struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_prescription_tenant"`

	// Prescription number
	PrescriptionNumber string `gorm:"size:50;uniqueIndex:idx_prescription_number_tenant;not null"`

	// Patient information
	EmployeeID uint     `gorm:"not null;index:idx_prescription_employee"`
	Employee   Employee `gorm:"foreignKey:EmployeeID"`
	IsForMain  bool     `gorm:"default:true;index"` // true = اصلی, false = تبعی
	DependentID *uint   `gorm:"index"`
	Dependent  *FamilyMember `gorm:"foreignKey:DependentID"`

	// Physician information
	PhysicianName         string  `gorm:"size:200;not null"`
	PhysicianNationalCode *string `gorm:"size:10;index"`
	PhysicianMedicalCode  *string `gorm:"size:20;index"` // کد نظام پزشکی
	PhysicianSpecialty    *string `gorm:"size:100"`
	PhysicianPhone        *string `gorm:"size:20"`

	// Prescription details
	PrescriptionDate time.Time `gorm:"not null;index"` // تاریخ نسخه
	PrescriptionType string    `gorm:"size:50;not null;index"` // OUTPATIENT, INPATIENT, EMERGENCY, etc.
	IsElectronic     bool      `gorm:"default:false;index"` // نسخه الکترونیک
	IsUrgent         bool      `gorm:"default:false"`

	// Diagnosis
	MainDiagnosisCode     *string `gorm:"size:20;index"` // ICD-10 اصلی
	MainDiagnosisTitle    *string `gorm:"size:300"`
	SecondaryDiagnosisCodes *string `gorm:"type:text"` // JSON array of ICD-10 codes
	DiagnosisNotes        *string `gorm:"type:text"`

	// Center/Facility
	CenterID   *uint   `gorm:"index"`
	Center     *Center `gorm:"foreignKey:CenterID"`
	CenterName *string `gorm:"size:200"`
	CenterType *string `gorm:"size:50"` // PHARMACY, CLINIC, HOSPITAL, LAB, etc.

	// Prescription content
	Items []PrescriptionItem `gorm:"foreignKey:PrescriptionID"`

	// Approval and validation
	RequiresApproval bool       `gorm:"default:false"`
	ApprovedByPhysician bool    `gorm:"default:false"`
	ApprovalDate     *time.Time
	ApprovedBy       *uint  // User ID
	ApprovalNotes    *string `gorm:"type:text"`

	// Insurance validation
	IsValidatedByInsurer bool       `gorm:"default:false;index"`
	ValidationDate       *time.Time
	ValidatedBy          *uint // User ID
	ValidationNotes      *string `gorm:"type:text"`
	RejectionReason      *string `gorm:"type:text"`

	// Electronic prescription integration
	TaminPrescriptionID *string `gorm:"size:100;index"` // شناسه نسخه تامین
	SepasReferenceID    *string `gorm:"size:100;index"` // شناسه سپاس
	NationalRxNumber    *string `gorm:"size:100;index"` // شماره نسخه ملی

	// Files and attachments
	ScannedImagePath     *string `gorm:"size:500"`
	PhysicianSignaturePath *string `gorm:"size:500"`
	AdditionalDocumentsPath *string `gorm:"type:text"` // JSON array of file paths

	// Conversion to claim
	IsConvertedToClaim bool  `gorm:"default:false;index"`
	ClaimID            *uint `gorm:"index"`
	Claim              *Claim `gorm:"foreignKey:ClaimID"`
	ConversionDate     *time.Time

	// Status tracking
	Status       string     `gorm:"size:20;not null;default:'DRAFT';index"` // DRAFT, SUBMITTED, VALIDATED, REJECTED, CONVERTED
	StatusReason *string    `gorm:"type:text"`
	SubmittedAt  *time.Time
	CompletedAt  *time.Time

	// Notes
	PhysicianNotes  *string `gorm:"type:text"` // یادداشت پزشک
	PharmacyNotes   *string `gorm:"type:text"` // یادداشت داروخانه
	InsurerNotes    *string `gorm:"type:text"` // یادداشت بیمه
	InternalNotes   *string `gorm:"type:text"` // یادداشت داخلی

	// Metadata
	IsActive bool `gorm:"default:true;index"`
}

// TableName specifies the table name for Prescription
func (Prescription) TableName() string {
	return "prescriptions"
}

// PrescriptionItem represents an item in a prescription
type PrescriptionItem struct {
	BaseModel
	TenantID uint `gorm:"not null;index"`

	// Prescription reference
	PrescriptionID uint         `gorm:"not null;index:idx_prescription_item_prescription"`
	Prescription   Prescription `gorm:"foreignKey:PrescriptionID"`

	// Item reference (Drug or Service)
	ItemID   *uint `gorm:"index:idx_prescription_item_item"`
	Item     *Item `gorm:"foreignKey:ItemID"`
	ItemType string `gorm:"size:20;not null"` // DRUG, SERVICE, EQUIPMENT

	// Manual entry (if not in database)
	ManualItemName *string `gorm:"size:300"`
	ManualItemCode *string `gorm:"size:50"`

	// Prescription details
	Quantity         int     `gorm:"not null"`           // تعداد
	Dosage           *string `gorm:"size:100"`          // دوز
	Frequency        *string `gorm:"size:100"`          // دفعات مصرف
	Duration         *string `gorm:"size:100"`          // مدت مصرف
	Route            *string `gorm:"size:50"`           // راه مصرف
	InstructionID    *uint
	Instruction      *Instruction `gorm:"foreignKey:InstructionID"`
	SpecialNotes     *string      `gorm:"type:text"` // یادداشت ویژه

	// Body site (for procedures)
	BodySiteID *uint     `gorm:"index"`
	BodySite   *BodySite `gorm:"foreignKey:BodySiteID"`

	// Priority and urgency
	IsUrgent  bool `gorm:"default:false"`
	Priority  int  `gorm:"default:0"`
	SortOrder int  `gorm:"default:0"`

	// Validation
	IsValidated       bool    `gorm:"default:false"`
	ValidationNotes   *string `gorm:"type:text"`
	SubstitutionAllowed bool `gorm:"default:true"` // امکان جایگزینی

	// Status
	Status   string `gorm:"size:20;default:'PENDING'"` // PENDING, APPROVED, REJECTED, DISPENSED
	IsActive bool   `gorm:"default:true"`
}

// TableName specifies the table name for PrescriptionItem
func (PrescriptionItem) TableName() string {
	return "prescription_items"
}

// CanConvertToClaim checks if prescription can be converted to claim
func (p *Prescription) CanConvertToClaim() bool {
	return p.Status == "VALIDATED" && !p.IsConvertedToClaim && len(p.Items) > 0
}

// GetPatientName returns the name of the patient
func (p *Prescription) GetPatientName() string {
	if p.IsForMain {
		return p.Employee.FirstName + " " + p.Employee.LastName
	}
	if p.Dependent != nil {
		return p.Dependent.FirstName + " " + p.Dependent.LastName
	}
	return ""
}
