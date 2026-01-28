package entity

import (
	"time"
)

// Drug - دارو
type Drug struct {
	BaseModel

	// کدها
	IRCCode     string `gorm:"size:20;uniqueIndex" json:"irc_code"`    // کد IRC
	GenericCode string `gorm:"size:20" json:"generic_code"`           // کد ژنریک
	Barcode     string `gorm:"size:20" json:"barcode"`                // بارکد

	// نام‌ها
	TitleFa     string `gorm:"size:300" json:"title_fa"`      // نام فارسی
	TitleEn     string `gorm:"size:300" json:"title_en"`      // نام انگلیسی
	GenericName string `gorm:"size:200" json:"generic_name"`  // نام ژنریک

	// گروه‌بندی
	GroupID    *uint      `json:"group_id"`
	Group      *DrugGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	CategoryID *uint      `json:"category_id"`
	Category   *DrugGroup `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// فرم دارویی
	Form       string `gorm:"size:100" json:"form"`        // قرص، شربت، آمپول، ...
	Strength   string `gorm:"size:100" json:"strength"`    // قدرت/غلظت
	Unit       string `gorm:"size:50" json:"unit"`         // واحد
	PackageQty int    `json:"package_qty"`                 // تعداد در بسته

	// تولیدکننده
	ManufacturerID *uint         `json:"manufacturer_id"`
	Manufacturer   *Manufacturer `gorm:"foreignKey:ManufacturerID" json:"manufacturer,omitempty"`
	CountryID      *uint         `json:"country_id"`

	// قیمت‌گذاری
	BasePrice       int64   `json:"base_price"`        // قیمت پایه
	TechnicalFee    int64   `json:"technical_fee"`     // حق فنی
	InsurancePrice  int64   `json:"insurance_price"`   // قیمت بیمه‌ای
	FranchiseRate   float32 `json:"franchise_rate"`    // نرخ فرانشیز

	// نوع
	IsOTC          bool `gorm:"default:false" json:"is_otc"`           // OTC
	IsNarcotic     bool `gorm:"default:false" json:"is_narcotic"`      // مخدر
	IsRefrigerated bool `gorm:"default:false" json:"is_refrigerated"`  // یخچالی
	IsSpecial      bool `gorm:"default:false" json:"is_special"`       // خاص
	NeedsPreAuth   bool `gorm:"default:false" json:"needs_pre_auth"`   // نیاز به پیش‌تایید

	// محدودیت‌ها
	MaxQtyPerPrescription int `json:"max_qty_per_prescription"` // حداکثر در هر نسخه
	MaxQtyPerMonth        int `json:"max_qty_per_month"`        // حداکثر در ماه

	// وضعیت
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	EffectiveAt *time.Time `json:"effective_at"`
	ExpiredAt   *time.Time `json:"expired_at"`

	// Relations
	Prices       []DrugPrice       `json:"prices,omitempty"`
	Interactions []DrugInteraction `gorm:"foreignKey:DrugID" json:"interactions,omitempty"`
}

// DrugGroup - گروه دارویی
type DrugGroup struct {
	BaseModel

	Code     string `gorm:"size:20" json:"code"`
	TitleFa  string `gorm:"size:200" json:"title_fa"`
	TitleEn  string `gorm:"size:200" json:"title_en"`
	ParentID *uint  `json:"parent_id"`
	Parent   *DrugGroup `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// DrugPrice - قیمت دارو (تاریخچه)
type DrugPrice struct {
	BaseModel

	DrugID uint `gorm:"index" json:"drug_id"`
	Drug   Drug `gorm:"foreignKey:DrugID" json:"-"`

	// قیمت‌ها
	ConsumerPrice  int64 `json:"consumer_price"`   // قیمت مصرف‌کننده
	InsurancePrice int64 `json:"insurance_price"`  // قیمت بیمه‌ای
	TechnicalFee   int64 `json:"technical_fee"`    // حق فنی

	// تاریخ اعتبار
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	// منبع
	Source    string `gorm:"size:50" json:"source"`  // IRC, Manual, ...
	ApprovedBy *uint `json:"approved_by"`
}

// DrugInteraction - تداخل دارویی
type DrugInteraction struct {
	BaseModel

	DrugID uint `gorm:"index" json:"drug_id"`
	Drug   Drug `gorm:"foreignKey:DrugID" json:"-"`

	InteractingDrugID uint `gorm:"index" json:"interacting_drug_id"`
	InteractingDrug   Drug `gorm:"foreignKey:InteractingDrugID" json:"interacting_drug,omitempty"`

	Severity    uint8  `json:"severity"`     // شدت (1: خفیف, 2: متوسط, 3: شدید)
	Description string `gorm:"size:1000" json:"description"`
	Action      string `gorm:"size:500" json:"action"` // اقدام لازم
}

// DrugAlternative - جایگزین دارویی
type DrugAlternative struct {
	BaseModel

	DrugID uint `gorm:"index" json:"drug_id"`
	Drug   Drug `gorm:"foreignKey:DrugID" json:"-"`

	AlternativeDrugID uint `gorm:"index" json:"alternative_drug_id"`
	AlternativeDrug   Drug `gorm:"foreignKey:AlternativeDrugID" json:"alternative_drug,omitempty"`

	Notes string `gorm:"size:500" json:"notes"`
}

// Manufacturer - تولیدکننده
type Manufacturer struct {
	BaseModel

	Name      string `gorm:"size:200" json:"name"`
	NameEn    string `gorm:"size:200" json:"name_en"`
	CountryID *uint  `json:"country_id"`
	IsActive  bool   `gorm:"default:true" json:"is_active"`
}

// DrugPrescriptionLimit - محدودیت تجویز دارو
type DrugPrescriptionLimit struct {
	TenantModel

	DrugID uint `gorm:"index" json:"drug_id"`
	Drug   Drug `gorm:"foreignKey:DrugID" json:"-"`

	// شرایط
	SpecialtyID  *uint `json:"specialty_id"`  // محدود به تخصص خاص
	DiagnosisID  *uint `json:"diagnosis_id"`  // محدود به تشخیص خاص
	AgeMin       *int  `json:"age_min"`       // حداقل سن
	AgeMax       *int  `json:"age_max"`       // حداکثر سن
	Gender       *Gender `json:"gender"`      // جنسیت

	// محدودیت‌ها
	MaxQtyPerPrescription int     `json:"max_qty_per_prescription"`
	MaxQtyPerMonth        int     `json:"max_qty_per_month"`
	MaxQtyPerYear         int     `json:"max_qty_per_year"`
	CoveragePercent       float32 `json:"coverage_percent"` // درصد پوشش

	// پیش‌تایید
	NeedsPreAuth  bool   `gorm:"default:false" json:"needs_pre_auth"`
	PreAuthReason string `gorm:"size:500" json:"pre_auth_reason"`

	// تاریخ اعتبار
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}
