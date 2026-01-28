package entity

import (
	"time"
)

// Person - شخص (اطلاعات پایه فردی)
type Person struct {
	TenantModel

	// اطلاعات هویتی
	NationalCode string `gorm:"size:10;uniqueIndex:idx_tenant_national" json:"national_code"` // کد ملی
	FirstName    string `gorm:"size:100" json:"first_name"`                                   // نام
	LastName     string `gorm:"size:100" json:"last_name"`                                    // نام خانوادگی
	FatherName   string `gorm:"size:100" json:"father_name"`                                  // نام پدر

	// اطلاعات شناسنامه
	BirthCertNo  string     `gorm:"size:20" json:"birth_cert_no"`  // شماره شناسنامه
	BirthDate    *time.Time `json:"birth_date"`                    // تاریخ تولد
	BirthPlaceID *uint      `json:"birth_place_id"`                // محل تولد
	IssuePlace   string     `gorm:"size:100" json:"issue_place"`   // محل صدور

	Gender        Gender  `json:"gender"`          // جنسیت
	MaritalStatus uint8   `json:"marital_status"`  // وضعیت تاهل
	NationalityID *uint   `json:"nationality_id"`  // ملیت

	// اطلاعات تماس
	Mobile  string `gorm:"size:15" json:"mobile"`
	Phone   string `gorm:"size:15" json:"phone"`
	Email   string `gorm:"size:100" json:"email"`
	Address string `gorm:"size:500" json:"address"`
	PostalCode string `gorm:"size:10" json:"postal_code"`

	// اطلاعات بانکی
	AccountNumber string `gorm:"size:30" json:"account_number"`  // شماره حساب
	ShebaNumber   string `gorm:"size:26" json:"sheba_number"`    // شماره شبا
	CardNumber    string `gorm:"size:16" json:"card_number"`     // شماره کارت

	// وضعیت
	IsDeceased   bool       `gorm:"default:false" json:"is_deceased"`
	DeceasedDate *time.Time `json:"deceased_date,omitempty"`

	// Relations
	BirthPlace  *Province    `gorm:"foreignKey:BirthPlaceID" json:"birth_place,omitempty"`
	PolicyMembers []PolicyMember `json:"policy_members,omitempty"`
}

// Employee - بیمه‌شده اصلی (کارمند)
type Employee struct {
	TenantModel

	PersonID uint   `gorm:"uniqueIndex:idx_tenant_person" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"person"`

	// اطلاعات شغلی
	EmployeeCode string `gorm:"size:20;uniqueIndex:idx_tenant_emp_code" json:"employee_code"` // کد پرسنلی
	JobTitle     string `gorm:"size:100" json:"job_title"`                                    // عنوان شغلی
	Department   string `gorm:"size:100" json:"department"`                                   // واحد سازمانی
	WorkPlaceID  *uint  `json:"work_place_id"`                                                // محل کار

	// تاریخ‌های استخدامی
	HireDate      *time.Time `json:"hire_date"`       // تاریخ استخدام
	RetirementDate *time.Time `json:"retirement_date"` // تاریخ بازنشستگی

	// وضعیت
	EmploymentStatus uint8 `json:"employment_status"` // شاغل، بازنشسته، ...
	EducationLevel   uint8 `json:"education_level"`   // سطح تحصیلات

	// گروه‌های خاص
	EmployeeGroupID        *uint `json:"employee_group_id"`
	EmployeeSpecialGroupID *uint `json:"employee_special_group_id"`

	// تنظیمات خاص
	HasSpecialDisease bool `gorm:"default:false" json:"has_special_disease"` // بیماری خاص

	// Relations
	WorkPlace      *Province          `gorm:"foreignKey:WorkPlaceID" json:"work_place,omitempty"`
	EmployeeGroup  *EmployeeGroup     `gorm:"foreignKey:EmployeeGroupID" json:"employee_group,omitempty"`
	FamilyMembers  []FamilyMember     `json:"family_members,omitempty"`
}

// FamilyMember - عضو خانواده (تبعی)
type FamilyMember struct {
	TenantModel

	EmployeeID uint     `gorm:"index" json:"employee_id"`
	Employee   Employee `gorm:"foreignKey:EmployeeID" json:"employee,omitempty"`

	PersonID uint   `json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"person"`

	RelationType RelationType `json:"relation_type"` // نسبت

	// تاریخ اعتبار
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	// وضعیت
	IsActive bool `gorm:"default:true" json:"is_active"`

	// اطلاعات تکمیلی برای فرزندان
	IsStudent       bool `gorm:"default:false" json:"is_student"`        // دانشجو
	IsMarried       bool `gorm:"default:false" json:"is_married"`        // متاهل
	HasDisability   bool `gorm:"default:false" json:"has_disability"`    // معلول
}

// EmployeeGroup - گروه کارمندان (برای تخفیفات/پوشش خاص)
type EmployeeGroup struct {
	BaseModel
	TenantID uint `gorm:"index" json:"tenant_id"`

	Name        string `gorm:"size:100" json:"name"`
	Description string `gorm:"size:500" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}

// Policy - بیمه‌نامه
type Policy struct {
	TenantModel

	PolicyNumber string `gorm:"size:50;uniqueIndex:idx_tenant_policy_no" json:"policy_number"` // شماره بیمه‌نامه

	// بیمه‌گر و رشته بیمه
	InsurerID        uint `json:"insurer_id"`          // بیمه‌گر
	InsuranceFieldID uint `json:"insurance_field_id"`  // رشته بیمه (درمان گروهی، ...)
	InsuranceTypeID  uint `json:"insurance_type_id"`   // نوع بیمه

	// تاریخ اعتبار
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	IssueDate  time.Time `json:"issue_date"`

	// مالی
	Premium            int64   `json:"premium"`              // حق بیمه
	StandardPriceRate  float32 `gorm:"default:1" json:"standard_price_rate"` // ضریب تعهد

	// تنظیمات
	MustDocDelivery       bool `gorm:"default:true" json:"must_doc_delivery"`        // الزام تحویل مدارک
	DocDeliveryDeadline   int  `gorm:"default:30" json:"doc_delivery_deadline"`      // مهلت تحویل (روز)

	// وضعیت
	Status uint8 `json:"status"` // فعال، غیرفعال، منقضی

	Description string `gorm:"size:500" json:"description"`

	// Relations
	Insurer       *Insurer        `gorm:"foreignKey:InsurerID" json:"insurer,omitempty"`
	PolicyMembers []PolicyMember  `json:"policy_members,omitempty"`
}

// PolicyMember - عضو بیمه‌نامه (ارتباط شخص با بیمه‌نامه)
type PolicyMember struct {
	TenantModel

	PolicyID uint   `gorm:"index" json:"policy_id"`
	Policy   Policy `gorm:"foreignKey:PolicyID" json:"policy,omitempty"`

	PersonID uint   `gorm:"index" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"person,omitempty"`

	// سرپرست (برای افراد تبعی)
	SupervisorID *uint         `json:"supervisor_id"`
	Supervisor   *PolicyMember `gorm:"foreignKey:SupervisorID" json:"supervisor,omitempty"`

	// کدهای ملی (برای استعلام سریع)
	MemberNationalCode     string `gorm:"size:10;index" json:"member_national_code"`
	SupervisorNationalCode string `gorm:"size:10;index" json:"supervisor_national_code"`

	// نوع عضویت
	DependencyType uint8 `json:"dependency_type"` // اصلی، همسر، فرزند، ...
	ContractType   uint8 `json:"contract_type"`   // رسمی، قراردادی، ...

	// واحد کاری
	WorkUnitID *uint `json:"work_unit_id"`

	// تاریخ اعتبار عضویت
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	// Relations
	Claims      []Claim        `json:"claims,omitempty"`
	FamilyMembers []PolicyMember `gorm:"foreignKey:SupervisorID" json:"family_members,omitempty"`
}

// Insurer - بیمه‌گر
type Insurer struct {
	BaseModel

	Name        string `gorm:"size:100" json:"name"`
	Code        string `gorm:"size:20;uniqueIndex" json:"code"`
	Logo        string `gorm:"size:255" json:"logo"`
	Description string `gorm:"size:500" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// تنظیمات API
	APIEndpoint string `gorm:"size:255" json:"api_endpoint"`
	APIKey      string `gorm:"size:255" json:"api_key"`
}

// Province - استان/شهر
type Province struct {
	BaseModel

	Name     string `gorm:"size:100" json:"name"`
	Code     string `gorm:"size:10" json:"code"`
	ParentID *uint  `json:"parent_id"` // برای شهرها
	IsActive bool   `gorm:"default:true" json:"is_active"`
}
