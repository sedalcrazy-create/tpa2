package entity

import "time"

// EmployeeIllness represents chronic illnesses and pre-existing conditions of employees
// Used for special coverage rules, waiting periods, and exclusions
type EmployeeIllness struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_emp_illness_tenant"`

	// Employee reference
	EmployeeID uint     `gorm:"not null;index:idx_emp_illness_employee"`
	Employee   Employee `gorm:"foreignKey:EmployeeID"`

	// Illness information
	ICD10Code     string  `gorm:"size:20;not null;index"` // ICD-10 diagnosis code
	ICD10Title    string  `gorm:"size:300;not null"`      // Disease name
	ICD10TitleFa  string  `gorm:"size:300;not null"`      // نام بیماری فارسی
	ICD10Category *string `gorm:"size:50;index"`          // Disease category

	// Diagnosis details
	DiagnosisDate    *time.Time `gorm:"index"` // تاریخ تشخیص
	DiagnosedBy      *string    `gorm:"size:200"` // پزشک تشخیص‌دهنده
	DiagnosisCenter  *string    `gorm:"size:200"` // مرکز تشخیص
	DiagnosisDetails *string    `gorm:"type:text"`

	// Severity and status
	Severity    string `gorm:"size:20;index"` // MILD, MODERATE, SEVERE, CRITICAL
	Status      string `gorm:"size:20;index"` // ACTIVE, CONTROLLED, CURED, DORMANT
	IsChronic   bool   `gorm:"default:false;index"` // بیماری مزمن
	IsCongenital bool  `gorm:"default:false;index"` // بیماری مادرزادی
	IsHereditary bool  `gorm:"default:false;index"` // بیماری ارثی

	// Insurance impact
	IsPreExisting        bool   `gorm:"default:false;index"` // بیماری قبل از بیمه
	RequiresDeclaration  bool   `gorm:"default:true"`        // نیاز به اظهار
	AffectsPremium       bool   `gorm:"default:false"`       // تاثیر بر حق بیمه
	PremiumLoadingPercent *float64 `gorm:"type:decimal(5,2)"`  // درصد افزایش حق بیمه

	// Coverage rules
	IsCovered            bool   `gorm:"default:true;index"`  // آیا تحت پوشش است
	CoverageStartDate    *time.Time `gorm:"index"`           // شروع پوشش
	CoverageEndDate      *time.Time                          // پایان پوشش
	WaitingPeriodDays    *int   `gorm:"type:int"`            // دوره انتظار (روز)
	ExclusionReason      *string `gorm:"type:text"`          // دلیل عدم پوشش
	SpecialConditions    *string `gorm:"type:text"`          // شرایط خاص پوشش
	AnnualCoverageLimit  *int64 `gorm:"type:bigint"`         // سقف سالانه اختصاصی

	// Treatment information
	CurrentTreatment     *string `gorm:"type:text"` // درمان فعلی
	RequiresMedication   bool    `gorm:"default:false"`
	RequiresMonitoring   bool    `gorm:"default:false"`
	LastVisitDate        *time.Time
	NextVisitDate        *time.Time

	// Related documents
	MedicalReportPath    *string `gorm:"size:500"` // مسیر گزارش پزشکی
	LabTestResultsPath   *string `gorm:"size:500"` // نتایج آزمایش
	ImagingResultsPath   *string `gorm:"size:500"` // نتایج تصویربرداری

	// Notes and remarks
	PhysicianNotes       *string `gorm:"type:text"` // یادداشت پزشک
	InsuranceNotes       *string `gorm:"type:text"` // یادداشت کارشناس بیمه
	InternalNotes        *string `gorm:"type:text"` // یادداشت داخلی

	// Approval workflow
	DeclaredByEmployee   bool       `gorm:"default:false"` // اظهار شده توسط کارمند
	DeclarationDate      *time.Time                        // تاریخ اظهار
	VerifiedByInsurer    bool       `gorm:"default:false"` // تایید شده
	VerificationDate     *time.Time                        // تاریخ تایید
	VerifiedBy           *uint                             // تایید کننده
	RejectionReason      *string    `gorm:"type:text"`     // دلیل رد

	// Status
	IsActive bool `gorm:"default:true;index"`
}

// TableName specifies the table name for EmployeeIllness
func (EmployeeIllness) TableName() string {
	return "employee_illnesses"
}

// IsWaitingPeriodPassed checks if waiting period has passed
func (ei *EmployeeIllness) IsWaitingPeriodPassed() bool {
	if !ei.IsCovered || ei.WaitingPeriodDays == nil {
		return true
	}

	if ei.CoverageStartDate == nil {
		return false
	}

	daysPassed := int(time.Since(*ei.CoverageStartDate).Hours() / 24)
	return daysPassed >= *ei.WaitingPeriodDays
}

// GetCoverageStatus returns the current coverage status
func (ei *EmployeeIllness) GetCoverageStatus() string {
	if !ei.IsActive {
		return "INACTIVE"
	}

	if !ei.IsCovered {
		return "EXCLUDED"
	}

	if !ei.IsWaitingPeriodPassed() {
		return "WAITING_PERIOD"
	}

	now := time.Now()
	if ei.CoverageEndDate != nil && now.After(*ei.CoverageEndDate) {
		return "EXPIRED"
	}

	return "COVERED"
}
