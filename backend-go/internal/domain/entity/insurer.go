package entity

// Insurer - بیمه‌گر / سازمان (Tenant)
type Insurer struct {
	BaseModel

	// اطلاعات اساسی
	Title      string `gorm:"size:200;not null" json:"title"`        // نام بیمه‌گر
	TitleEn    string `gorm:"size:200" json:"title_en"`              // نام انگلیسی
	Code       string `gorm:"size:20;uniqueIndex" json:"code"`       // کد بیمه‌گر
	NationalID string `gorm:"size:11" json:"national_id"`            // شناسه ملی
	EconomicCode string `gorm:"size:14" json:"economic_code"`        // کد اقتصادی

	// اطلاعات تماس
	Phone      string `gorm:"size:15" json:"phone"`
	Mobile     string `gorm:"size:15" json:"mobile"`
	Email      string `gorm:"size:100" json:"email"`
	Website    string `gorm:"size:200" json:"website"`
	Address    string `gorm:"type:text" json:"address"`
	PostalCode string `gorm:"size:10" json:"postal_code"`

	// وضعیت
	IsActive bool `gorm:"default:true" json:"is_active"`

	// تنظیمات
	Logo       string `gorm:"size:255" json:"logo"`
	PrimaryColor string `gorm:"size:7" json:"primary_color"` // #RRGGBB
}

// TableName specifies the table name
func (Insurer) TableName() string {
	return "insurers"
}
