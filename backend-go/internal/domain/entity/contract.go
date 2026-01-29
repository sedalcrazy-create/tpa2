package entity

import "time"

// ContractType represents types of insurance contracts
type ContractType struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_contract_type_tenant"`

	// Code and naming
	Code        string `gorm:"size:50;uniqueIndex:idx_contract_type_code_tenant;not null"`
	Title       string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	// Category
	Category string `gorm:"size:50;index"` // INDIVIDUAL, GROUP, CORPORATE, GOVERNMENT

	// Default settings
	DefaultDurationMonths *int    `gorm:"type:int"` // مدت پیش‌فرض (ماه)
	DefaultRenewalType    *string `gorm:"size:20"`  // AUTO, MANUAL
	DefaultGracePeriodDays *int   `gorm:"type:int"` // مهلت پیش‌فرض

	// Template
	ContractTemplate *string `gorm:"type:text"` // قالب قرارداد (HTML/JSON)
	TermsTemplate    *string `gorm:"type:text"` // شرایط و ضوابط

	// Status
	IsActive  bool `gorm:"default:true;index"`
	SortOrder int  `gorm:"default:0"`

	// Relations
	Contracts []Contract `gorm:"foreignKey:ContractTypeID"`
}

// TableName specifies the table name for ContractType
func (ContractType) TableName() string {
	return "contract_types"
}

// Contract represents an insurance contract between insurer and employer/individual
type Contract struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_contract_tenant"`

	// Contract identification
	ContractNumber string `gorm:"size:50;uniqueIndex:idx_contract_number_tenant;not null"`
	ContractTypeID uint   `gorm:"not null;index:idx_contract_type"`
	ContractType   ContractType `gorm:"foreignKey:ContractTypeID"`

	// Parties
	// Insurer side (always the tenant/insurer)
	InsurerName      string  `gorm:"size:200;not null"`
	InsurerCode      *string `gorm:"size:50"`
	InsurerRepName   *string `gorm:"size:200"` // نماینده بیمه‌گر
	InsurerRepTitle  *string `gorm:"size:100"` // سمت نماینده

	// Employer/Insured side
	EmployerName     string  `gorm:"size:200;not null;index"` // نام کارفرما/بیمه‌گذار
	EmployerCode     *string `gorm:"size:50;index"`           // کد کارفرما
	EmployerNationalID *string `gorm:"size:11;index"`        // شناسه ملی
	EmployerEconomicCode *string `gorm:"size:14;index"`      // کد اقتصادی
	EmployerRepName  *string `gorm:"size:200"`                // نماینده کارفرما
	EmployerRepTitle *string `gorm:"size:100"`                // سمت نماینده
	EmployerPhone    *string `gorm:"size:20"`
	EmployerEmail    *string `gorm:"size:100"`
	EmployerAddress  *string `gorm:"type:text"`

	// Contract period
	StartDate     time.Time  `gorm:"not null;index"` // تاریخ شروع
	EndDate       *time.Time `gorm:"index"`          // تاریخ پایان
	SignDate      *time.Time `gorm:"index"`          // تاریخ امضا
	EffectiveDate *time.Time `gorm:"index"`          // تاریخ اجرا

	// Contract details
	TotalInsured      int    `gorm:"default:0"`           // تعداد کل بیمه‌شدگان
	MainInsured       int    `gorm:"default:0"`           // تعداد بیمه‌شده اصلی
	DependentsAllowed int    `gorm:"default:0"`           // تعداد افراد تبعی مجاز
	MaxDependentsPerMain int `gorm:"default:4"`          // حداکثر تبعی به ازای هر اصلی

	// Financial terms
	TotalPremiumAmount   *int64   `gorm:"type:bigint"`       // کل حق بیمه
	PremiumPerPerson     *int64   `gorm:"type:bigint"`       // حق بیمه به ازای هر نفر
	PaymentMethod        *string  `gorm:"size:50"`           // MONTHLY, QUARTERLY, ANNUALLY, LUMP_SUM
	PaymentDay           *int     `gorm:"type:int"`          // روز پرداخت
	AdvancePaymentPercent *float64 `gorm:"type:decimal(5,2)"` // درصد پیش‌پرداخت
	AdvancePaymentAmount *int64   `gorm:"type:bigint"`       // مبلغ پیش‌پرداخت

	// Coverage limits
	AnnualCoverageLimit  *int64 `gorm:"type:bigint"` // سقف سالانه کل
	PerClaimLimit        *int64 `gorm:"type:bigint"` // سقف هر ادعا
	FranchiseAmount      *int64 `gorm:"type:bigint"` // فرانشیز
	FranchisePercentage  *float64 `gorm:"type:decimal(5,2)"` // درصد فرانشیز

	// Renewal terms
	RenewalType        string `gorm:"size:20;default:'MANUAL'"` // AUTO, MANUAL
	GracePeriodDays    int    `gorm:"default:30"`               // مهلت تجدید
	RenewalNotifyDays  int    `gorm:"default:60"`               // اطلاع‌رسانی قبل از انقضا
	AllowEarlyTermination bool `gorm:"default:false"`           // امکان فسخ زودهنگام
	TerminationPenalty *int64 `gorm:"type:bigint"`              // جریمه فسخ

	// Addendums and amendments
	AddendumCount    int     `gorm:"default:0"`       // تعداد الحاقیه‌ها
	LastAddendumDate *time.Time                       // تاریخ آخرین الحاقیه
	LastAddendumNo   *string `gorm:"size:50"`        // شماره آخرین الحاقیه

	// Documents
    ContractFilePath     *string `gorm:"size:500"` // مسیر فایل قرارداد
	SignedContractPath   *string `gorm:"size:500"` // قرارداد امضا شده
	AddendumsPath        *string `gorm:"type:text"` // الحاقیه‌ها (JSON)
	AttachmentsPath      *string `gorm:"type:text"` // ضمائم (JSON)

	// Approval workflow
	Status           string     `gorm:"size:20;not null;default:'DRAFT';index"` // DRAFT, PENDING, APPROVED, ACTIVE, SUSPENDED, TERMINATED
	ApprovedBy       *uint                                                       // کاربر تایید کننده
	ApprovedAt       *time.Time                                                  // تاریخ تایید
	ApprovalNotes    *string    `gorm:"type:text"`                               // یادداشت تایید

	// Termination
	TerminationDate   *time.Time                 // تاریخ فسخ
	TerminationReason *string `gorm:"type:text"` // دلیل فسخ
	TerminatedBy      *uint                      // کاربر فسخ کننده

	// Notes
	Terms            string  `gorm:"type:text"`     // شرایط و ضوابط
	SpecialConditions *string `gorm:"type:text"`    // شرایط خاص
	Notes            *string `gorm:"type:text"`     // یادداشت‌ها
	InternalNotes    *string `gorm:"type:text"`     // یادداشت داخلی

	// System tracking
	CreatedBy        *uint                         // کاربر ثبت کننده
	IsActive         bool   `gorm:"default:true;index"`

	// Relations
	Policies         []Policy `gorm:"foreignKey:ContractID"`
	InsuranceHistories []InsuranceHistory `gorm:"foreignKey:PolicyID"` // Through Policy
}

// TableName specifies the table name for Contract
func (Contract) TableName() string {
	return "contracts"
}

// IsCurrentlyActive checks if contract is currently active
func (c *Contract) IsCurrentlyActive() bool {
	if !c.IsActive || c.Status != "ACTIVE" {
		return false
	}

	now := time.Now()

	if now.Before(c.StartDate) {
		return false
	}

	if c.EndDate != nil && now.After(*c.EndDate) {
		return false
	}

	return true
}

// DaysUntilExpiry returns days until contract expires
func (c *Contract) DaysUntilExpiry() int {
	if c.EndDate == nil {
		return -1 // No expiry
	}

	days := int(time.Until(*c.EndDate).Hours() / 24)
	if days < 0 {
		return 0 // Already expired
	}

	return days
}

// NeedsRenewalNotification checks if renewal notification should be sent
func (c *Contract) NeedsRenewalNotification() bool {
	if c.EndDate == nil {
		return false
	}

	daysUntil := c.DaysUntilExpiry()
	return daysUntil > 0 && daysUntil <= c.RenewalNotifyDays
}

// IsInGracePeriod checks if contract is in grace period after expiry
func (c *Contract) IsInGracePeriod() bool {
	if c.EndDate == nil {
		return false
	}

	now := time.Now()
	if now.Before(*c.EndDate) {
		return false // Not expired yet
	}

	daysSinceExpiry := int(now.Sub(*c.EndDate).Hours() / 24)
	return daysSinceExpiry <= c.GracePeriodDays
}
