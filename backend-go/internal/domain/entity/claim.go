package entity

import (
	"time"
)

// Claim - ادعای درمانی (نسخه)
type Claim struct {
	TenantModel

	// شناسه‌ها
	TrackingCode string `gorm:"size:50;uniqueIndex:idx_tenant_tracking" json:"tracking_code"` // کد پیگیری
	HID          string `gorm:"size:50" json:"hid"`                                           // شناسه سلامت
	MRN          string `gorm:"size:50" json:"mrn"`                                           // شماره پرونده پزشکی

	// بیمه‌شده
	PolicyMemberID uint         `gorm:"index" json:"policy_member_id"`
	PolicyMember   PolicyMember `gorm:"foreignKey:PolicyMemberID" json:"policy_member,omitempty"`

	// مرکز درمانی
	CenterID uint    `gorm:"index" json:"center_id"`
	Center   *Center `gorm:"foreignKey:CenterID" json:"center,omitempty"`

	// بسته ارسالی
	PackageID *uint    `gorm:"index" json:"package_id"`
	Package   *Package `gorm:"foreignKey:PackageID" json:"package,omitempty"`

	// نوع و وضعیت
	ClaimType     ClaimType     `json:"claim_type"`     // نوع ادعا (دارو، بستری، ...)
	Status        ClaimStatus   `json:"status"`         // وضعیت
	AdmissionType AdmissionType `json:"admission_type"` // سرپایی/بستری

	// تاریخ‌های مراجعه
	AdmissionDate time.Time  `json:"admission_date"`         // تاریخ پذیرش
	DischargeDate *time.Time `json:"discharge_date"`         // تاریخ ترخیص
	ServiceDate   time.Time  `json:"service_date"`           // تاریخ خدمت

	// وضعیت ترخیص (برای بستری)
	DischargeCondition *DischargeCondition `json:"discharge_condition"`

	// پزشکان
	AttendingDoctorID *uint        `json:"attending_doctor_id"` // پزشک معالج
	AttendingDoctor   *Provider    `gorm:"foreignKey:AttendingDoctorID" json:"attending_doctor,omitempty"`
	ReferrerID        *uint        `json:"referrer_id"`         // پزشک ارجاع‌دهنده
	Referrer          *Provider    `gorm:"foreignKey:ReferrerID" json:"referrer,omitempty"`

	// بیمه پایه
	BasicInsurerID *uint    `json:"basic_insurer_id"`
	BasicInsurer   *Insurer `gorm:"foreignKey:BasicInsurerID" json:"basic_insurer,omitempty"`

	// پیش‌تایید
	PreAuthID *uint    `json:"pre_auth_id"`
	PreAuth   *PreAuth `gorm:"foreignKey:PreAuthID" json:"pre_auth,omitempty"`

	// ادعای والد (برای ادعاهای تکمیلی)
	ParentID *uint  `json:"parent_id"`
	Parent   *Claim `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	// مبالغ
	RequestAmount    int64 `json:"request_amount"`     // مبلغ درخواستی
	ApprovedAmount   int64 `json:"approved_amount"`    // مبلغ تایید شده
	BasicInsShare    int64 `json:"basic_ins_share"`    // سهم بیمه پایه
	SupplementShare  int64 `json:"supplement_share"`   // سهم بیمه تکمیلی
	PatientShare     int64 `json:"patient_share"`      // سهم بیمار
	Deduction        int64 `json:"deduction"`          // کسورات

	// ارزیابی
	HandlerUserID     *uint      `json:"handler_user_id"`     // ارزیاب
	HandlerDate       *time.Time `json:"handler_date"`
	CheckingDesc      string     `gorm:"size:1000" json:"checking_desc"` // توضیحات ارزیابی

	// بررسی دارو
	DrugCheckingUserID *uint      `json:"drug_checking_user_id"`
	DrugCheckingDate   *time.Time `json:"drug_checking_date"`

	// تحویل مدارک
	DocDeliveryDate    *time.Time `json:"doc_delivery_date"`
	DocHandlerUserID   *uint      `json:"doc_handler_user_id"`

	// پرداخت
	PaymentDate *time.Time `json:"payment_date"`

	// باطل شدن
	IsVoid      bool       `gorm:"default:false" json:"is_void"`
	VoidDate    *time.Time `json:"void_date"`
	VoidMessage string     `gorm:"size:500" json:"void_message"`
	VoidUserID  *uint      `json:"void_user_id"`

	// ثبت‌کننده
	RegisterUserID uint      `json:"register_user_id"`
	RegisterDate   time.Time `json:"register_date"`

	// Relations
	Items       []ClaimItem      `json:"items,omitempty"`
	Diagnoses   []ClaimDiagnosis `json:"diagnoses,omitempty"`
	Attachments []ClaimAttachment `json:"attachments,omitempty"`
	Notes       []ClaimNote      `json:"notes,omitempty"`
}

// ClaimItem - قلم ادعا (دارو یا خدمت)
type ClaimItem struct {
	TenantModel

	ClaimID uint  `gorm:"index" json:"claim_id"`
	Claim   Claim `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`

	// نوع قلم (یکی از دو تا)
	DrugID    *uint    `json:"drug_id"`
	Drug      *Drug    `gorm:"foreignKey:DrugID" json:"drug,omitempty"`
	ServiceID *uint    `json:"service_id"`
	Service   *Service `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	ItemID    *uint    `json:"item_id"` // Universal Item reference (new)
	Item      *Item    `gorm:"foreignKey:ItemID" json:"item,omitempty"`

	// Prescription reference (NEW)
	PrescriptionItemID *uint             `json:"prescription_item_id"` // ارجاع به قلم نسخه
	PrescriptionItem   *PrescriptionItem `gorm:"foreignKey:PrescriptionItemID" json:"prescription_item,omitempty"`

	// Instruction reference (NEW)
	InstructionID *uint        `json:"instruction_id"` // دستور مصرف
	Instruction   *Instruction `gorm:"foreignKey:InstructionID" json:"instruction,omitempty"`

	// اطلاعات از مدرک
	DocItemCode  string `gorm:"size:50" json:"doc_item_code"`   // کد روی مدرک
	DocItemTitle string `gorm:"size:255" json:"doc_item_title"` // نام روی مدرک
	DocCount     *int   `json:"doc_count"`                       // تعداد روی مدرک

	// تعداد و مقدار
	Count       int `json:"count"`        // تعداد تایید شده
	BatchNumber int `json:"batch_number"` // شماره بچ

	// Usage instructions (NEW)
	Dosage    *string `gorm:"size:100" json:"dosage"`    // دوز
	Frequency *string `gorm:"size:100" json:"frequency"` // دفعات
	Duration  *string `gorm:"size:100" json:"duration"`  // مدت

	// تاریخ
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	// مبالغ
	RequestPrice    int64 `json:"request_price"`    // مبلغ درخواستی
	ConfirmedPrice  int64 `json:"confirmed_price"`  // مبلغ تایید شده
	BasicInsShare   int64 `json:"basic_ins_share"`  // سهم بیمه پایه
	SupplementShare int64 `json:"supplement_share"` // سهم بیمه تکمیلی
	Franchise       int64 `json:"franchise"`        // فرانشیز
	Deduction       int64 `json:"deduction"`        // کسورات
	OutsideBasic    int64 `json:"outside_basic"`    // خارج از تعهد پایه

	// ویژگی‌های خاص
	Tooth          *uint8   `json:"tooth"`            // شماره دندان
	LeftEyeDiopter *float32 `json:"left_eye_diopter"` // دیوپتر چشم چپ
	RightEyeDiopter *float32 `json:"right_eye_diopter"` // دیوپتر چشم راست

	// Body site reference (NEW - direct reference)
	BodySiteID *uint     `json:"body_site_id"` // محل انجام خدمت
	BodySite   *BodySite `gorm:"foreignKey:BodySiteID" json:"body_site,omitempty"`

	// ارائه‌دهنده
	ProviderID *uint     `json:"provider_id"`
	Provider   *Provider `gorm:"foreignKey:ProviderID" json:"provider,omitempty"`

	// قلم والد
	ParentID *uint      `json:"parent_id"`
	Parent   *ClaimItem `gorm:"foreignKey:ParentID" json:"parent,omitempty"`

	// Relations
	BodySites   []ClaimItemBodySite   `json:"body_sites,omitempty"` // Multiple body sites (legacy)
	ReasonCodes []ClaimItemReasonCode `json:"reason_codes,omitempty"`
}

// ClaimDiagnosis - تشخیص ادعا
type ClaimDiagnosis struct {
	BaseModel

	ClaimID uint  `gorm:"index" json:"claim_id"`
	Claim   Claim `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`

	DiagnosisID uint      `json:"diagnosis_id"`
	Diagnosis   *Diagnosis `gorm:"foreignKey:DiagnosisID" json:"diagnosis,omitempty"`

	IsPrimary   bool   `gorm:"default:false" json:"is_primary"` // تشخیص اصلی
	Description string `gorm:"size:500" json:"description"`
}

// ClaimAttachment - پیوست ادعا
type ClaimAttachment struct {
	BaseModel

	ClaimID uint  `gorm:"index" json:"claim_id"`
	Claim   Claim `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`

	FileType   string `gorm:"size:50" json:"file_type"`   // نوع فایل
	FileName   string `gorm:"size:255" json:"file_name"`  // نام فایل
	FilePath   string `gorm:"size:500" json:"file_path"`  // مسیر فایل
	FileSize   int64  `json:"file_size"`                  // حجم فایل
	MimeType   string `gorm:"size:100" json:"mime_type"`
	UploadedBy uint   `json:"uploaded_by"`
}

// ClaimNote - یادداشت ادعا
type ClaimNote struct {
	BaseModel

	ClaimID uint  `gorm:"index" json:"claim_id"`
	Claim   Claim `gorm:"foreignKey:ClaimID" json:"claim,omitempty"`

	Note     string `gorm:"size:2000" json:"note"`
	NoteType uint8  `json:"note_type"` // نوع یادداشت
	UserID   uint   `json:"user_id"`
}

// ClaimItemBodySite - عضو بدن مرتبط با قلم
type ClaimItemBodySite struct {
	BaseModel

	ClaimItemID uint       `gorm:"index" json:"claim_item_id"`
	ClaimItem   ClaimItem  `gorm:"foreignKey:ClaimItemID" json:"-"`

	BodySiteID uint      `json:"body_site_id"`
	BodySite   *BodySite `gorm:"foreignKey:BodySiteID" json:"body_site,omitempty"`
}

// ClaimItemReasonCode - کد دلیل کسورات
type ClaimItemReasonCode struct {
	BaseModel

	ClaimItemID uint       `gorm:"index" json:"claim_item_id"`
	ClaimItem   ClaimItem  `gorm:"foreignKey:ClaimItemID" json:"-"`

	ReasonCodeID uint        `json:"reason_code_id"`
	ReasonCode   *ReasonCode `gorm:"foreignKey:ReasonCodeID" json:"reason_code,omitempty"`
}

// PreAuth - پیش‌تایید (علی‌الحساب)
type PreAuth struct {
	TenantModel

	PersonID uint   `gorm:"index" json:"person_id"`
	Person   Person `gorm:"foreignKey:PersonID" json:"person,omitempty"`

	Subject string `gorm:"size:500" json:"subject"` // موضوع

	Amount  int64 `json:"amount"`   // مبلغ
	Type    uint8 `json:"type"`     // نوع

	PaymentDate *time.Time `json:"payment_date"`

	RegisterUserID uint      `json:"register_user_id"`
	RegisterDate   time.Time `json:"register_date"`

	// ادعای مرتبط (constraint:false to avoid circular dependency with Claim.PreAuthID)
	ClaimID *uint  `gorm:"index" json:"claim_id"`
	Claim   *Claim `gorm:"-" json:"claim,omitempty"` // Loaded manually, no FK constraint
}

// Diagnosis - تشخیص (ICD)
type Diagnosis struct {
	BaseModel

	Code        string `gorm:"size:20;uniqueIndex" json:"code"` // کد ICD
	TitleFa     string `gorm:"size:500" json:"title_fa"`
	TitleEn     string `gorm:"size:500" json:"title_en"`
	Description string `gorm:"size:1000" json:"description"`

	GroupID *uint           `json:"group_id"`
	Group   *DiagnosisGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`

	IsActive bool `gorm:"default:true" json:"is_active"`
}

// DiagnosisGroup - گروه تشخیص
type DiagnosisGroup struct {
	BaseModel

	Code     string `gorm:"size:20" json:"code"`
	TitleFa  string `gorm:"size:200" json:"title_fa"`
	TitleEn  string `gorm:"size:200" json:"title_en"`
	IsActive bool   `gorm:"default:true" json:"is_active"`
}

// BodySite - عضو بدن (anatomical body sites)
type BodySite struct {
	BaseModel

	Code    string `gorm:"size:20;uniqueIndex" json:"code"` // Unique code
	TitleFa string `gorm:"size:200;index" json:"title_fa"`
	TitleEn string `gorm:"size:200" json:"title_en"`

	// Hierarchy (NEW)
	ParentID *uint      `gorm:"index" json:"parent_id"`
	Parent   *BodySite  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children []BodySite `gorm:"foreignKey:ParentID" json:"children,omitempty"`

	// Category (NEW)
	Category string `gorm:"size:50;index" json:"category"` // HEAD, TRUNK, UPPER_LIMB, LOWER_LIMB, etc.

	// Laterality (NEW)
	Side string `gorm:"size:10" json:"side"` // LEFT, RIGHT, BILATERAL, NONE

	// Medical codes (NEW)
	ICD10Code   *string `gorm:"size:20;index" json:"icd10_code"`  // ICD-10 anatomical code
	SNOMEDCode  *string `gorm:"size:20;index" json:"snomed_code"` // SNOMED CT code
	CPTModifier *string `gorm:"size:10" json:"cpt_modifier"`      // CPT modifier for laterality

	// Description (NEW)
	Description string `gorm:"type:text" json:"description"`

	// Status (NEW)
	IsActive  bool `gorm:"default:true;index" json:"is_active"`
	SortOrder int  `gorm:"default:0" json:"sort_order"`

	// Legacy group reference
	GroupID *uint          `json:"group_id"`
	Group   *BodySiteGroup `gorm:"foreignKey:GroupID" json:"group,omitempty"`
}

// BodySiteGroup - گروه اعضای بدن
type BodySiteGroup struct {
	BaseModel

	TitleFa string `gorm:"size:200" json:"title_fa"`
	TitleEn string `gorm:"size:200" json:"title_en"`
}

// ReasonCode - کد دلیل کسورات
type ReasonCode struct {
	BaseModel

	Code        string `gorm:"size:20;uniqueIndex" json:"code"`
	TitleFa     string `gorm:"size:500" json:"title_fa"`
	Description string `gorm:"size:1000" json:"description"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}
