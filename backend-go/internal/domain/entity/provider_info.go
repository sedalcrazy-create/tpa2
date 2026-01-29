package entity

// ProviderInfo represents detailed information about healthcare providers (physicians, pharmacists, etc.)
// This is separate from Center - a center can have multiple providers
type ProviderInfo struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_provider_tenant"`

	// Provider identification
	NationalCode   string  `gorm:"size:10;uniqueIndex:idx_provider_national;not null"` // کد ملی
	MedicalCode    *string `gorm:"size:20;uniqueIndex:idx_provider_medical"` // کد نظام پزشکی
	PharmacyCode   *string `gorm:"size:20;uniqueIndex:idx_provider_pharmacy"` // کد داروخانه
	LicenseNumber  *string `gorm:"size:50;index"` // شماره پروانه
	InsuranceCode  *string `gorm:"size:20;index"` // کد بیمه

	// Personal information
	FirstName      string `gorm:"size:100;not null"`
	LastName       string `gorm:"size:100;not null"`
	FatherName     *string `gorm:"size:100"`
	BirthDate      *string `gorm:"size:10"` // YYYY-MM-DD
	Gender         string `gorm:"size:10;index"` // MALE, FEMALE

	// Professional information
	ProviderType   string `gorm:"size:50;not null;index"` // PHYSICIAN, PHARMACIST, DENTIST, NURSE, TECHNICIAN, etc.
	Specialty      *string `gorm:"size:100;index"` // تخصص
	SubSpecialty   *string `gorm:"size:100"` // فوق تخصص
	AcademicDegree *string `gorm:"size:50"` // GENERAL, SPECIALIST, SUBSPECIALIST, etc.

	// Education
	UniversityName     *string `gorm:"size:200"`
	GraduationYear     *int    `gorm:"type:int"`
	MedicalCouncilName *string `gorm:"size:100"` // نظام پزشکی استان

	// License and certification
	LicenseIssueDate   *string `gorm:"size:10"` // تاریخ صدور پروانه
	LicenseExpiryDate  *string `gorm:"size:10"` // تاریخ انقضا پروانه
	CertificationBody  *string `gorm:"size:200"` // مرجع صدور
	IsLicenseValid     bool    `gorm:"default:true;index"`

	// Centers association
	Centers []Center `gorm:"many2many:provider_center_mappings;"`

	// Contact information
	Phone       *string `gorm:"size:20"`
	Mobile      *string `gorm:"size:20;index"`
	Email       *string `gorm:"size:100;index"`
	Website     *string `gorm:"size:200"`

	// Address
	ProvinceID *uint     `gorm:"index"`
	Province   *Province `gorm:"foreignKey:ProvinceID"`
	CityID     *uint     `gorm:"index"`
	City       *City     `gorm:"foreignKey:CityID"`
	Address    *string   `gorm:"type:text"`
	PostalCode *string   `gorm:"size:10"`

	// Profile
	Biography          *string `gorm:"type:text"` // بیوگرافی
	PhotoPath          *string `gorm:"size:500"` // مسیر عکس
	SignaturePath      *string `gorm:"size:500"` // مسیر امضا
	Achievements       *string `gorm:"type:text"` // افتخارات
	ResearchInterests  *string `gorm:"type:text"` // زمینه‌های تحقیقاتی
	Languages          *string `gorm:"size:200"` // زبان‌های تسلط

	// Statistics
	TotalPrescriptions *int `gorm:"type:int;default:0"` // تعداد نسخه
	TotalClaims        *int `gorm:"type:int;default:0"` // تعداد ادعا
	AvgClaimAmount     *int64 `gorm:"type:bigint"` // میانگین مبلغ ادعا
	LastPrescriptionDate *string `gorm:"size:10"` // آخرین نسخه

	// Verification
	IsVerified         bool    `gorm:"default:false;index"` // تایید شده
	VerifiedAt         *string `gorm:"size:10"` // تاریخ تایید
	VerifiedBy         *uint   // کاربر تایید کننده
	VerificationNotes  *string `gorm:"type:text"`

	// Flags and status
	IsActive           bool `gorm:"default:true;index"`
	IsBlacklisted      bool `gorm:"default:false;index"` // لیست سیاه
	BlacklistReason    *string `gorm:"type:text"`
	BlacklistedAt      *string `gorm:"size:10"`
	IsSuspended        bool `gorm:"default:false;index"` // معلق
	SuspensionReason   *string `gorm:"type:text"`
	SuspendedAt        *string `gorm:"size:10"`

	// Notes
	Notes              *string `gorm:"type:text"`
	InternalNotes      *string `gorm:"type:text"`

	// Relations
	Prescriptions []Prescription `gorm:"foreignKey:PhysicianMedicalCode;references:MedicalCode"`
}

// TableName specifies the table name for ProviderInfo
func (ProviderInfo) TableName() string {
	return "provider_infos"
}

// ProviderCenterMapping represents the many-to-many relationship between providers and centers
type ProviderCenterMapping struct {
	BaseModel

	ProviderID uint         `gorm:"not null;index"`
	Provider   ProviderInfo `gorm:"foreignKey:ProviderID"`
	CenterID   uint         `gorm:"not null;index"`
	Center     Center       `gorm:"foreignKey:CenterID"`

	// Mapping properties
	Role          *string `gorm:"size:50"` // HEAD, CONSULTANT, STAFF, etc.
	StartDate     *string `gorm:"size:10"` // تاریخ شروع همکاری
	EndDate       *string `gorm:"size:10"` // تاریخ پایان همکاری
	IsActive      bool    `gorm:"default:true"`
	IsPrimary     bool    `gorm:"default:false"` // محل کار اصلی
	WorkSchedule  *string `gorm:"type:text"` // برنامه کاری (JSON)
	Notes         *string `gorm:"type:text"`
}

// TableName specifies the table name for ProviderCenterMapping
func (ProviderCenterMapping) TableName() string {
	return "provider_center_mappings"
}

// GetFullName returns provider's full name
func (p *ProviderInfo) GetFullName() string {
	return p.FirstName + " " + p.LastName
}

// IsPhysician checks if provider is a physician
func (p *ProviderInfo) IsPhysician() bool {
	return p.ProviderType == "PHYSICIAN"
}

// IsPharmacist checks if provider is a pharmacist
func (p *ProviderInfo) IsPharmacist() bool {
	return p.ProviderType == "PHARMACIST"
}

// CanPrescribe checks if provider can write prescriptions
func (p *ProviderInfo) CanPrescribe() bool {
	return p.IsActive && !p.IsBlacklisted && !p.IsSuspended &&
	       p.IsLicenseValid && (p.IsPhysician() || p.ProviderType == "DENTIST")
}
