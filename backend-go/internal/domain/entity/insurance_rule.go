package entity

import "time"

// InsuranceRule represents business rules for insurance policies
// Defines coverage limits, deductibles, waiting periods, and other policy rules
type InsuranceRule struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_insurance_rule_tenant"`

	// Insurance reference
	InsuranceID uint      `gorm:"not null;index:idx_insurance_rule_insurance"`
	Insurance   Insurance `gorm:"foreignKey:InsuranceID"`

	// Rule identification
	Code        string `gorm:"size:50;not null;index"`
	Title       string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`
	RuleType    string `gorm:"size:50;not null;index"` // COVERAGE, LIMIT, WAITING_PERIOD, DEDUCTIBLE, etc.

	// Coverage limits (سقف تعهدات)
	AnnualLimit         *int64 `gorm:"type:bigint"` // سقف سالانه کل
	PerClaimLimit       *int64 `gorm:"type:bigint"` // سقف هر ادعا
	LifetimeLimit       *int64 `gorm:"type:bigint"` // سقف مادام العمر
	DailyLimit          *int64 `gorm:"type:bigint"` // سقف روزانه
	MonthlyLimit        *int64 `gorm:"type:bigint"` // سقف ماهانه
	PerServiceLimit     *int64 `gorm:"type:bigint"` // سقف هر خدمت
	PerDrugLimit        *int64 `gorm:"type:bigint"` // سقف هر دارو
	HospitalizationDays *int   `gorm:"type:int"`    // حداکثر روزهای بستری

	// Service-specific limits
	DrugLimit         *int64 `gorm:"type:bigint"` // سقف دارو
	DentalLimit       *int64 `gorm:"type:bigint"` // سقف دندان
	OpticalLimit      *int64 `gorm:"type:bigint"` // سقف عینک/لنز
	PhysiotherapyLimit *int64 `gorm:"type:bigint"` // سقف فیزیوتراپی
	LabLimit          *int64 `gorm:"type:bigint"` // سقف آزمایش
	ImagingLimit      *int64 `gorm:"type:bigint"` // سقف تصویربرداری
	SurgeryLimit      *int64 `gorm:"type:bigint"` // سقف جراحی

	// Deductible (فرانشیز)
	DeductibleAmount     *int64   `gorm:"type:bigint"`       // مبلغ فرانشیز
	DeductiblePercentage *float64 `gorm:"type:decimal(5,2)"` // درصد فرانشیز
	DeductibleType       *string  `gorm:"size:20"`           // FIXED, PERCENTAGE, COMBINED

	// Waiting periods (دوره انتظار)
	GeneralWaitingDays   *int `gorm:"type:int"` // دوره انتظار عمومی
	DentalWaitingDays    *int `gorm:"type:int"` // دوره انتظار دندان
	OpticalWaitingDays   *int `gorm:"type:int"` // دوره انتظار عینک
	MaternityWaitingDays *int `gorm:"type:int"` // دوره انتظار زایمان
	SurgeryWaitingDays   *int `gorm:"type:int"` // دوره انتظار جراحی

	// Co-payment (همیار پرداخت)
	CoPaymentPercentage *float64 `gorm:"type:decimal(5,2)"` // درصد پرداخت بیمه‌شده
	CoPaymentAmount     *int64   `gorm:"type:bigint"`       // مبلغ ثابت پرداخت

	// Age/Gender restrictions
	MinAge          *int    `gorm:"type:int"`
	MaxAge          *int    `gorm:"type:int"`
	GenderLimit     *string `gorm:"size:10"`  // MALE, FEMALE, ALL
	AppliestoMain   bool    `gorm:"default:true"`
	AppliesToDep    bool    `gorm:"default:true"`
	DependentLimit  *int    `gorm:"type:int"` // حداکثر تعداد افراد تبعی
	ChildAgeLimit   *int    `gorm:"type:int"` // سن فرزند

	// Special conditions
	RequiresPreAuth       bool `gorm:"default:false"` // نیاز به تایید قبلی
	RequiresPrescription  bool `gorm:"default:false"` // نیاز به نسخه
	RequiresReferal       bool `gorm:"default:false"` // نیاز به ارجاع
	AllowEmergency        bool `gorm:"default:true"`  // پوشش اورژانس
	AllowChronicIllness   bool `gorm:"default:true"`  // پوشش بیماری مزمن
	AllowPreExisting      bool `gorm:"default:false"` // پوشش بیماری قبلی
	NetworkOnly           bool `gorm:"default:false"` // فقط شبکه طرف قرارداد

	// Exclusions
	ExcludedServices string `gorm:"type:text"` // JSON array of excluded service codes
	ExcludedDrugs    string `gorm:"type:text"` // JSON array of excluded drug codes
	ExcludedICD10    string `gorm:"type:text"` // JSON array of excluded diagnosis codes

	// Renewal rules
	AutoRenewal       bool `gorm:"default:false"` // تمدید خودکار
	RenewalGraceDays  *int `gorm:"type:int"`      // مهلت تمدید
	PremiumIncrease   *float64 `gorm:"type:decimal(5,2)"` // درصد افزایش حق بیمه

	// Priority and status
	Priority int  `gorm:"default:0;index"`
	IsActive bool `gorm:"default:true;index"`

	// Validity period
	EffectiveDate *time.Time `gorm:"index"`
	ExpiryDate    *time.Time `gorm:"index"`

	// Rule engine integration
	RuleEngineCode string `gorm:"type:text"` // Grule DSL code
}

// TableName specifies the table name for InsuranceRule
func (InsuranceRule) TableName() string {
	return "insurance_rules"
}

// IsValid checks if the rule is currently valid
func (ir *InsuranceRule) IsValid() bool {
	if !ir.IsActive {
		return false
	}

	now := time.Now()

	if ir.EffectiveDate != nil && now.Before(*ir.EffectiveDate) {
		return false
	}

	if ir.ExpiryDate != nil && now.After(*ir.ExpiryDate) {
		return false
	}

	return true
}

// CheckWaitingPeriod checks if waiting period has passed
func (ir *InsuranceRule) CheckWaitingPeriod(policyStartDate time.Time, serviceType string) bool {
	now := time.Now()
	daysSinceStart := int(now.Sub(policyStartDate).Hours() / 24)

	switch serviceType {
	case "DENTAL":
		if ir.DentalWaitingDays != nil && daysSinceStart < *ir.DentalWaitingDays {
			return false
		}
	case "OPTICAL":
		if ir.OpticalWaitingDays != nil && daysSinceStart < *ir.OpticalWaitingDays {
			return false
		}
	case "MATERNITY":
		if ir.MaternityWaitingDays != nil && daysSinceStart < *ir.MaternityWaitingDays {
			return false
		}
	case "SURGERY":
		if ir.SurgeryWaitingDays != nil && daysSinceStart < *ir.SurgeryWaitingDays {
			return false
		}
	default:
		if ir.GeneralWaitingDays != nil && daysSinceStart < *ir.GeneralWaitingDays {
			return false
		}
	}

	return true
}

// GetServiceLimit returns the limit for a specific service type
func (ir *InsuranceRule) GetServiceLimit(serviceType string) *int64 {
	switch serviceType {
	case "DRUG":
		return ir.DrugLimit
	case "DENTAL":
		return ir.DentalLimit
	case "OPTICAL":
		return ir.OpticalLimit
	case "PHYSIOTHERAPY":
		return ir.PhysiotherapyLimit
	case "LAB":
		return ir.LabLimit
	case "IMAGING":
		return ir.ImagingLimit
	case "SURGERY":
		return ir.SurgeryLimit
	default:
		return nil
	}
}
