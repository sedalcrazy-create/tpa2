package entity

import "time"

// EmployeeSpecialDiscount represents individual discount rules for specific employees
// Different from CustomEmployeeCode which is group-based
type EmployeeSpecialDiscount struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_emp_discount_tenant"`

	// Employee reference
	EmployeeID uint     `gorm:"not null;index:idx_emp_discount_employee"`
	Employee   Employee `gorm:"foreignKey:EmployeeID"`

	// Discount information
	Code        string `gorm:"size:50;not null;index"`
	Title       string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`
	DiscountType string `gorm:"size:50;not null;index"` // PERCENTAGE, FIXED_AMOUNT, COMBINED

	// Discount values
	DiscountPercentage *float64 `gorm:"type:decimal(5,2)"` // درصد تخفیف
	DiscountAmount     *int64   `gorm:"type:bigint"`       // مبلغ ثابت تخفیف
	MaxDiscountAmount  *int64   `gorm:"type:bigint"`       // حداکثر مبلغ تخفیف

	// Scope filters (what this discount applies to)
	ApplyToItemID       *uint `gorm:"index"` // تخفیف روی یک قلم خاص
	Item                *Item `gorm:"foreignKey:ApplyToItemID"`
	ApplyToCategoryID   *uint `gorm:"index"` // تخفیف روی یک دسته
	Category            *ItemCategory `gorm:"foreignKey:ApplyToCategoryID"`
	ApplyToGroupID      *uint `gorm:"index"` // تخفیف روی یک گروه
	Group               *ItemGroup `gorm:"foreignKey:ApplyToGroupID"`
	ApplyToServiceType  *string `gorm:"size:50"` // DRUG, DENTAL, OPTICAL, etc.

	// Application rules
	ApplyOnBasePrice    bool `gorm:"default:true"`  // اعمال روی قیمت پایه
	ApplyAfterInsurance bool `gorm:"default:false"` // اعمال بعد از بیمه
	CombineWithOthers   bool `gorm:"default:true"`  // قابل ترکیب با تخفیف‌های دیگر
	Priority            int  `gorm:"default:0;index"` // اولویت اعمال

	// Limits
	MaxUsagePerDay      *int `gorm:"type:int"` // حداکثر استفاده در روز
	MaxUsagePerMonth    *int `gorm:"type:int"` // حداکثر استفاده در ماه
	MaxUsagePerYear     *int `gorm:"type:int"` // حداکثر استفاده در سال
	MaxTotalUsage       *int `gorm:"type:int"` // حداکثر استفاده کل
	MaxDiscountPerYear  *int64 `gorm:"type:bigint"` // حداکثر مبلغ تخفیف سالانه
	UsageCount          int  `gorm:"default:0"` // تعداد استفاده شده
	TotalDiscountGiven  int64 `gorm:"type:bigint;default:0"` // مجموع تخفیف داده شده

	// Reason and approval
	Reason            string  `gorm:"size:200;not null"` // دلیل تخفیف
	GrantedBy         *uint   // کاربر ثبت کننده
	GrantedAt         *time.Time
	ApprovedBy        *uint   // کاربر تایید کننده
	ApprovedAt        *time.Time
	ApprovalNotes     *string `gorm:"type:text"`

	// Validity period
	StartDate *time.Time `gorm:"not null;index"` // شروع اعتبار
	EndDate   *time.Time `gorm:"index"`          // پایان اعتبار

	// Special conditions
	RequiresApproval bool    `gorm:"default:false"` // نیاز به تایید هر بار
	RequiresDocument bool    `gorm:"default:false"` // نیاز به مدرک
	DocumentPath     *string `gorm:"size:500"`      // مسیر مدرک
	Conditions       *string `gorm:"type:text"`     // شرایط خاص (JSON)

	// Status
	IsActive     bool `gorm:"default:true;index"`
	IsSuspended  bool `gorm:"default:false;index"`
	SuspendedAt  *time.Time
	SuspendReason *string `gorm:"type:text"`

	// Notes
	Notes         *string `gorm:"type:text"`
	InternalNotes *string `gorm:"type:text"`
}

// TableName specifies the table name for EmployeeSpecialDiscount
func (EmployeeSpecialDiscount) TableName() string {
	return "employee_special_discounts"
}

// IsValid checks if discount is currently valid
func (esd *EmployeeSpecialDiscount) IsValid() bool {
	if !esd.IsActive || esd.IsSuspended {
		return false
	}

	now := time.Now()

	if esd.StartDate != nil && now.Before(*esd.StartDate) {
		return false
	}

	if esd.EndDate != nil && now.After(*esd.EndDate) {
		return false
	}

	// Check usage limits
	if esd.MaxTotalUsage != nil && esd.UsageCount >= *esd.MaxTotalUsage {
		return false
	}

	return true
}

// CalculateDiscount calculates discount amount for given price
func (esd *EmployeeSpecialDiscount) CalculateDiscount(basePrice int64) int64 {
	if !esd.IsValid() {
		return 0
	}

	discount := int64(0)

	// Calculate based on type
	switch esd.DiscountType {
	case "PERCENTAGE":
		if esd.DiscountPercentage != nil {
			discount = int64(float64(basePrice) * (*esd.DiscountPercentage / 100.0))
		}
	case "FIXED_AMOUNT":
		if esd.DiscountAmount != nil {
			discount = *esd.DiscountAmount
		}
	case "COMBINED":
		// Apply both percentage and fixed
		if esd.DiscountPercentage != nil {
			discount += int64(float64(basePrice) * (*esd.DiscountPercentage / 100.0))
		}
		if esd.DiscountAmount != nil {
			discount += *esd.DiscountAmount
		}
	}

	// Apply max limit
	if esd.MaxDiscountAmount != nil && discount > *esd.MaxDiscountAmount {
		discount = *esd.MaxDiscountAmount
	}

	// Don't exceed base price
	if discount > basePrice {
		discount = basePrice
	}

	return discount
}

// CanUse checks if discount can be used now
func (esd *EmployeeSpecialDiscount) CanUse() bool {
	if !esd.IsValid() {
		return false
	}

	// Check total usage limit
	if esd.MaxTotalUsage != nil && esd.UsageCount >= *esd.MaxTotalUsage {
		return false
	}

	// Additional time-based checks would go here (daily, monthly, yearly)
	// Would require querying usage history

	return true
}

// IncrementUsage increments the usage counter
func (esd *EmployeeSpecialDiscount) IncrementUsage(discountAmount int64) {
	esd.UsageCount++
	esd.TotalDiscountGiven += discountAmount
}
