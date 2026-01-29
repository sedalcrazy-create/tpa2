package entity

import "time"

// CustomEmployeeCode represents special employee codes with custom pricing/discount rules
// Used for retired employees, special groups, VIP personnel, etc.
type CustomEmployeeCode struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_cec_tenant"`

	// Code information
	Code  string `gorm:"size:50;not null;uniqueIndex:idx_cec_code_tenant"` // e.g., "RETIRED", "VIP", "VETERAN"
	Title string `gorm:"size:200;not null"`                                 // فارسی: "بازنشسته", "VIP", "جانباز"

	// Discount settings
	DiscountPercentage *float64 `gorm:"type:decimal(5,2)"` // درصد تخفیف (0-100)
	DiscountAmount     *int64   `gorm:"type:bigint"`       // مبلغ ثابت تخفیف (ریال)

	// Price limit settings
	MaxPricePercentage *float64 `gorm:"type:decimal(5,2)"` // حداکثر درصد از قیمت پایه

	// Special flags
	IsRetired      bool `gorm:"default:false;index"` // بازنشسته
	NoLimitation   bool `gorm:"default:false"`       // بدون محدودیت سقف
	SpecialGroup   bool `gorm:"default:false"`       // گروه ویژه
	PriorityAccess bool `gorm:"default:false"`       // دسترسی اولویت‌دار

	// Validity period
	StartDate *time.Time `gorm:"index"`
	EndDate   *time.Time `gorm:"index"`

	// Relations
	Employees []Employee `gorm:"foreignKey:CustomEmployeeCodeID"`
}

// TableName specifies the table name for CustomEmployeeCode
func (CustomEmployeeCode) TableName() string {
	return "custom_employee_codes"
}

// IsActive checks if the code is currently valid
func (c *CustomEmployeeCode) IsActive() bool {
	now := time.Now()

	if c.StartDate != nil && now.Before(*c.StartDate) {
		return false
	}

	if c.EndDate != nil && now.After(*c.EndDate) {
		return false
	}

	return true
}

// GetDiscount calculates the discount amount for a given price
func (c *CustomEmployeeCode) GetDiscount(price int64) int64 {
	if c.DiscountAmount != nil {
		return *c.DiscountAmount
	}

	if c.DiscountPercentage != nil {
		return int64(float64(price) * (*c.DiscountPercentage / 100.0))
	}

	return 0
}

// GetMaxPrice calculates the maximum allowed price
func (c *CustomEmployeeCode) GetMaxPrice(basePrice int64) int64 {
	if c.MaxPricePercentage == nil {
		return basePrice
	}

	return int64(float64(basePrice) * (*c.MaxPricePercentage / 100.0))
}
