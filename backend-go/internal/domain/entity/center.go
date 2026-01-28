package entity

import (
	"time"
)

// Center - مرکز درمانی طرف قرارداد
type Center struct {
	TenantModel

	// اطلاعات پایه
	Title       string     `gorm:"size:200" json:"title"`                                  // نام مرکز
	SiamID      string     `gorm:"size:20;uniqueIndex:idx_tenant_siam" json:"siam_id"`    // کد سیام
	Code        string     `gorm:"size:20" json:"code"`                                    // کد داخلی
	Type        CenterType `json:"type"`                                                   // نوع مرکز
	Level       int        `json:"level"`                                                  // سطح اعتباربخشی

	// موقعیت
	ProvinceID *uint     `json:"province_id"`
	Province   *Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	CityID     *uint     `json:"city_id"`
	City       *Province `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Address    string    `gorm:"size:500" json:"address"`
	PostalCode string    `gorm:"size:10" json:"postal_code"`
	Phone      string    `gorm:"size:20" json:"phone"`
	Fax        string    `gorm:"size:20" json:"fax"`
	Email      string    `gorm:"size:100" json:"email"`
	Website    string    `gorm:"size:200" json:"website"`

	// مالک/مدیر
	OwnerName    string `gorm:"size:100" json:"owner_name"`
	ManagerName  string `gorm:"size:100" json:"manager_name"`
	ManagerPhone string `gorm:"size:20" json:"manager_phone"`

	// نوع مالکیت
	DependencyType uint8 `json:"dependency_type"` // دولتی، خصوصی، ...

	// اطلاعات مالی/پرداخت
	PaymentID            string `gorm:"size:50" json:"payment_id"`              // شناسه پرداخت
	AccountNumber        string `gorm:"size:30" json:"account_number"`          // شماره حساب
	AccountOwnerName     string `gorm:"size:100" json:"account_owner_name"`     // نام صاحب حساب
	ShebaNumber          string `gorm:"size:26" json:"sheba_number"`            // شماره شبا
	EconomicCode         string `gorm:"size:20" json:"economic_code"`           // کد اقتصادی
	NationalID           string `gorm:"size:15" json:"national_id"`             // شناسه ملی

	// وضعیت قرارداد
	ContractStatus ContractStatus `json:"contract_status"`
	ContractNumber string         `gorm:"size:50" json:"contract_number"`
	ContractStart  *time.Time     `json:"contract_start"`
	ContractEnd    *time.Time     `json:"contract_end"`

	// تنظیمات
	IsActive bool `gorm:"default:true" json:"is_active"`

	// Relations
	Claims   []Claim   `json:"claims,omitempty"`
	Packages []Package `json:"packages,omitempty"`
	Users    []User    `json:"users,omitempty"` // کاربران مرکز
}

// CenterContract - قرارداد با مرکز درمانی
type CenterContract struct {
	TenantModel

	CenterID uint   `gorm:"index" json:"center_id"`
	Center   Center `gorm:"foreignKey:CenterID" json:"center,omitempty"`

	ContractNumber string    `gorm:"size:50" json:"contract_number"`
	ContractType   uint8     `json:"contract_type"`
	StartDate      time.Time `json:"start_date"`
	EndDate        time.Time `json:"end_date"`

	// شرایط قرارداد
	DiscountPercent   float32 `json:"discount_percent"`    // درصد تخفیف
	PaymentTermDays   int     `json:"payment_term_days"`   // مهلت پرداخت (روز)
	MinSettlementAmount int64 `json:"min_settlement_amount"` // حداقل مبلغ تسویه

	Terms       string `gorm:"type:text" json:"terms"`        // متن قرارداد
	Attachments string `gorm:"size:1000" json:"attachments"`  // پیوست‌ها

	Status    uint8      `json:"status"`
	SignedAt  *time.Time `json:"signed_at"`
	SignedBy  *uint      `json:"signed_by"`

	Notes string `gorm:"size:1000" json:"notes"`
}

// Package - بسته اسناد (ارسالی از مرکز)
type Package struct {
	TenantModel

	// مرکز
	CenterID uint   `gorm:"index" json:"center_id"`
	Center   Center `gorm:"foreignKey:CenterID" json:"center,omitempty"`

	// واحد کاری (اداره امور)
	WorkUnitID *uint     `json:"work_unit_id"`
	WorkUnit   *WorkUnit `gorm:"foreignKey:WorkUnitID" json:"work_unit,omitempty"`

	// اطلاعات نامه
	Title            string     `gorm:"size:200" json:"title"`
	LetterNumber     string     `gorm:"size:50" json:"letter_number"`      // شماره نامه
	LetterDate       *time.Time `json:"letter_date"`                       // تاریخ نامه
	ReceiveLetterDate *time.Time `json:"receive_letter_date"`              // تاریخ دریافت نامه
	LetterImageURL   string     `gorm:"size:500" json:"letter_image_url"` // تصویر نامه

	// آمار
	ClaimCount   int   `json:"claim_count"`    // تعداد ادعاها
	TotalAmount  int64 `json:"total_amount"`   // مبلغ کل درخواستی
	ApprovedAmount int64 `json:"approved_amount"` // مبلغ کل تایید شده
	DeductionAmount int64 `json:"deduction_amount"` // کسورات کل

	// وضعیت
	Status PackageStatus `json:"status"`

	// ارزیابی
	CheckingUserID *uint      `json:"checking_user_id"`
	CheckingDate   *time.Time `json:"checking_date"`

	// پرداخت
	IsPayment     bool       `gorm:"default:false" json:"is_payment"`
	PaymentDate   *time.Time `json:"payment_date"`
	PaymentAmount *int64     `json:"payment_amount"`
	PaymentRef    string     `gorm:"size:50" json:"payment_ref"` // شماره پیگیری پرداخت

	// بسته والد (برای تکمیلی)
	ParentID *uint    `json:"parent_id"`
	Parent   *Package `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	// ایجاد
	CreateDate time.Time `json:"create_date"`

	// Relations
	Claims []Claim `json:"claims,omitempty"`
}

// Settlement - تسویه حساب با مرکز
type Settlement struct {
	TenantModel

	CenterID uint   `gorm:"index" json:"center_id"`
	Center   Center `gorm:"foreignKey:CenterID" json:"center,omitempty"`

	// دوره تسویه
	PeriodStart time.Time `json:"period_start"`
	PeriodEnd   time.Time `json:"period_end"`

	// مبالغ
	TotalAmount     int64 `json:"total_amount"`      // مبلغ کل
	DeductionAmount int64 `json:"deduction_amount"`  // کسورات
	PayableAmount   int64 `json:"payable_amount"`    // مبلغ قابل پرداخت
	PaidAmount      int64 `json:"paid_amount"`       // مبلغ پرداخت شده

	// وضعیت
	Status uint8 `json:"status"` // در انتظار، تایید، پرداخت شده

	// پرداخت
	PaymentDate    *time.Time `json:"payment_date"`
	PaymentRef     string     `gorm:"size:50" json:"payment_ref"`
	BankReference  string     `gorm:"size:50" json:"bank_reference"`  // شماره پیگیری بانکی

	// تایید
	ApprovedBy   *uint      `json:"approved_by"`
	ApprovedDate *time.Time `json:"approved_date"`

	Notes string `gorm:"size:1000" json:"notes"`

	// Relations
	Packages []Package `gorm:"many2many:settlement_packages;" json:"packages,omitempty"`
}

// WorkUnit - واحد کاری (اداره امور)
type WorkUnit struct {
	TenantModel

	Title      string `gorm:"size:200" json:"title"`
	Code       string `gorm:"size:20" json:"code"`
	ProvinceID *uint  `json:"province_id"`
	Province   *Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"`
	IsActive   bool   `gorm:"default:true" json:"is_active"`
}

// Provider - ارائه‌دهنده خدمت (پزشک/متخصص)
type Provider struct {
	BaseModel

	// اطلاعات پایه
	NationalCode string `gorm:"size:10" json:"national_code"`
	FirstName    string `gorm:"size:100" json:"first_name"`
	LastName     string `gorm:"size:100" json:"last_name"`
	MedicalCode  string `gorm:"size:20;uniqueIndex" json:"medical_code"` // کد نظام پزشکی

	// تخصص
	SpecialtyID *uint      `json:"specialty_id"`
	Specialty   *Specialty `gorm:"foreignKey:SpecialtyID" json:"specialty,omitempty"`
	DegreeID    *uint      `json:"degree_id"`

	// تماس
	Phone  string `gorm:"size:20" json:"phone"`
	Mobile string `gorm:"size:15" json:"mobile"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// Specialty - تخصص پزشکی
type Specialty struct {
	BaseModel

	Code    string `gorm:"size:20" json:"code"`
	TitleFa string `gorm:"size:200" json:"title_fa"`
	TitleEn string `gorm:"size:200" json:"title_en"`

	GroupID *uint           `json:"group_id"`
	Group   *SpecialtyGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// SpecialtyGroup - گروه تخصص
type SpecialtyGroup struct {
	BaseModel

	TitleFa  string `gorm:"size:200" json:"title_fa"`
	TitleEn  string `gorm:"size:200" json:"title_en"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
