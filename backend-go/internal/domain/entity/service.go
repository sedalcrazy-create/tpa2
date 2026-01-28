package entity

import (
	"time"
)

// Service - خدمت درمانی
type Service struct {
	BaseModel

	// کدها
	Code        string `gorm:"size:20;uniqueIndex" json:"code"`     // کد خدمت
	NationalCode string `gorm:"size:20" json:"national_code"`       // کد ملی
	OldCode     string `gorm:"size:20" json:"old_code"`             // کد قدیم

	// نام‌ها
	TitleFa string `gorm:"size:500" json:"title_fa"`
	TitleEn string `gorm:"size:500" json:"title_en"`

	// گروه‌بندی
	GroupID    *uint         `json:"group_id"`
	Group      *ServiceGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
	CategoryID *uint         `json:"category_id"`
	Category   *ServiceGroup `gorm:"foreignKey:CategoryID" json:"category,omitempty"`

	// نوع خدمت
	ServiceTypeID *uint        `json:"service_type_id"`
	ServiceType   *ServiceType `gorm:"foreignKey:ServiceTypeID" json:"service_type,omitempty"`

	// قیمت‌گذاری پایه
	RelativeValue   float32 `json:"relative_value"`    // ارزش نسبی (K)
	TechnicalValue  float32 `json:"technical_value"`   // ارزش فنی
	ProfessionalValue float32 `json:"professional_value"` // ارزش حرفه‌ای

	BasePrice       int64   `json:"base_price"`        // قیمت پایه
	FranchiseRate   float32 `json:"franchise_rate"`    // نرخ فرانشیز

	// محدودیت‌ها
	MaxPerDay     int `json:"max_per_day"`      // حداکثر در روز
	MaxPerMonth   int `json:"max_per_month"`    // حداکثر در ماه
	MaxPerYear    int `json:"max_per_year"`     // حداکثر در سال
	GenderLimit   *Gender `json:"gender_limit"` // محدودیت جنسیتی
	AgeMin        *int `json:"age_min"`
	AgeMax        *int `json:"age_max"`

	// تنظیمات
	NeedsPreAuth     bool `gorm:"default:false" json:"needs_pre_auth"`
	NeedsBodySite    bool `gorm:"default:false" json:"needs_body_site"`    // نیاز به عضو بدن
	NeedsDiagnosis   bool `gorm:"default:false" json:"needs_diagnosis"`    // نیاز به تشخیص
	IsHospitalOnly   bool `gorm:"default:false" json:"is_hospital_only"`   // فقط بستری
	IsOutpatientOnly bool `gorm:"default:false" json:"is_outpatient_only"` // فقط سرپایی

	// تخصص مجاز
	AllowedSpecialties string `gorm:"size:500" json:"allowed_specialties"` // لیست تخصص‌های مجاز (کاما جدا)

	// وضعیت
	IsActive    bool       `gorm:"default:true" json:"is_active"`
	EffectiveAt *time.Time `json:"effective_at"`
	ExpiredAt   *time.Time `json:"expired_at"`

	// Relations
	Prices []ServicePrice `json:"prices,omitempty"`
}

// ServiceGroup - گروه خدمات
type ServiceGroup struct {
	BaseModel

	Code     string `gorm:"size:20" json:"code"`
	TitleFa  string `gorm:"size:200" json:"title_fa"`
	TitleEn  string `gorm:"size:200" json:"title_en"`
	ParentID *uint  `json:"parent_id"`
	Parent   *ServiceGroup `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	// نوع قیمت‌گذاری
	PricingType uint8 `json:"pricing_type"` // K, ثابت، درصدی

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// ServiceType - نوع خدمت
type ServiceType struct {
	BaseModel

	Code    string `gorm:"size:20" json:"code"`
	TitleFa string `gorm:"size:100" json:"title_fa"`
	TitleEn string `gorm:"size:100" json:"title_en"`

	// تنظیمات پیش‌فرض
	DefaultFranchiseRate float32 `json:"default_franchise_rate"`
	IsActive             bool    `gorm:"default:true" json:"is_active"`
}

// ServicePrice - قیمت خدمت (تاریخچه)
type ServicePrice struct {
	BaseModel

	ServiceID uint    `gorm:"index" json:"service_id"`
	Service   Service `gorm:"foreignKey:ServiceID" json:"-"`

	// ضرایب K
	KPublic      int64 `json:"k_public"`       // K دولتی
	KPrivate     int64 `json:"k_private"`      // K خصوصی
	KCharity     int64 `json:"k_charity"`      // K خیریه

	// قیمت‌ها
	TechnicalFee     int64 `json:"technical_fee"`      // حق فنی
	ProfessionalFee  int64 `json:"professional_fee"`   // حق تخصص
	AnesthesiaFee    int64 `json:"anesthesia_fee"`     // حق بیهوشی

	// تاریخ اعتبار
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	// منبع
	Source     string `gorm:"size:50" json:"source"`
	ApprovedBy *uint  `json:"approved_by"`
}

// ServiceRelation - ارتباط خدمات (شامل/مستثنی)
type ServiceRelation struct {
	BaseModel

	ServiceID uint    `gorm:"index" json:"service_id"`
	Service   Service `gorm:"foreignKey:ServiceID" json:"-"`

	RelatedServiceID uint    `gorm:"index" json:"related_service_id"`
	RelatedService   Service `gorm:"foreignKey:RelatedServiceID" json:"related_service,omitempty"`

	RelationType uint8  `json:"relation_type"` // 1: شامل, 2: مستثنی, 3: الزامی
	Notes        string `gorm:"size:500" json:"notes"`
}

// ServiceCoverageLimit - محدودیت پوشش خدمت (برای هر بیمه‌گر)
type ServiceCoverageLimit struct {
	TenantModel

	ServiceID uint    `gorm:"index" json:"service_id"`
	Service   Service `gorm:"foreignKey:ServiceID" json:"-"`

	// شرایط
	SpecialtyID   *uint  `json:"specialty_id"`
	DiagnosisID   *uint  `json:"diagnosis_id"`
	CenterTypeID  *uint  `json:"center_type_id"`
	AgeMin        *int   `json:"age_min"`
	AgeMax        *int   `json:"age_max"`
	Gender        *Gender `json:"gender"`

	// محدودیت‌ها
	MaxPerDay       int     `json:"max_per_day"`
	MaxPerMonth     int     `json:"max_per_month"`
	MaxPerYear      int     `json:"max_per_year"`
	CoveragePercent float32 `json:"coverage_percent"`
	MaxCoverage     int64   `json:"max_coverage"` // سقف ریالی

	// پیش‌تایید
	NeedsPreAuth  bool   `gorm:"default:false" json:"needs_pre_auth"`
	PreAuthReason string `gorm:"size:500" json:"pre_auth_reason"`

	// تاریخ اعتبار
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// Tariff - تعرفه (K ضریب)
type Tariff struct {
	TenantModel

	Year int    `gorm:"index" json:"year"` // سال
	Name string `gorm:"size:100" json:"name"`

	// ضرایب
	KPublic      int64   `json:"k_public"`       // K دولتی
	KPrivate     int64   `json:"k_private"`      // K خصوصی
	KCharity     int64   `json:"k_charity"`      // K خیریه
	KUniversity  int64   `json:"k_university"`   // K دانشگاهی

	// نرخ فرانشیز پیش‌فرض
	DefaultFranchiseOutpatient float32 `json:"default_franchise_outpatient"` // سرپایی
	DefaultFranchiseInpatient  float32 `json:"default_franchise_inpatient"`  // بستری

	// تاریخ اعتبار
	EffectiveFrom time.Time  `json:"effective_from"`
	EffectiveTo   *time.Time `json:"effective_to"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}
