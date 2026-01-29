package entity

import "time"

// ItemPriceCondition represents pricing rules and conditions for drugs/services
// This is the core pricing engine entity based on Refah system
type ItemPriceCondition struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_ipc_tenant"`

	// Item reference
	ItemID *uint `gorm:"index:idx_ipc_item"` // NULL = applies to all items
	Item   *Item `gorm:"foreignKey:ItemID"`

	// Category filters (if ItemID is null)
	CategoryID    *uint         `gorm:"index:idx_ipc_category"`
	Category      *ItemCategory `gorm:"foreignKey:CategoryID"`
	SubCategoryID *uint         `gorm:"index:idx_ipc_subcategory"`
	SubCategory   *ItemCategory `gorm:"foreignKey:SubCategoryID"`
	GroupID       *uint         `gorm:"index:idx_ipc_group"`
	Group         *ItemGroup    `gorm:"foreignKey:GroupID"`

	// Insurance reference
	InsuranceID *uint      `gorm:"index:idx_ipc_insurance"`
	Insurance   *Insurance `gorm:"foreignKey:InsuranceID"`

	// Pricing rules
	CoveragePercentage *float64 `gorm:"type:decimal(5,2)"` // درصد پوشش (0-100)
	MaxCoverageAmount  *int64   `gorm:"type:bigint"`       // حداکثر مبلغ پوشش (ریال)
	MinCoverageAmount  *int64   `gorm:"type:bigint"`       // حداقل مبلغ پوشش
	FixedAmount        *int64   `gorm:"type:bigint"`       // مبلغ ثابت پوشش

	// Franchise (فرانشیز)
	FranchisePercentage *float64 `gorm:"type:decimal(5,2)"` // درصد فرانشیز
	FranchiseAmount     *int64   `gorm:"type:bigint"`       // مبلغ ثابت فرانشیز
	MaxFranchise        *int64   `gorm:"type:bigint"`       // حداکثر فرانشیز

	// Deductible (کسر بیمه پایه)
	BaseInsuranceShare *float64 `gorm:"type:decimal(5,2)"` // سهم بیمه پایه (درصد)

	// Quantity limits
	MaxQuantityPerDay   *int `gorm:"type:int"`
	MaxQuantityPerMonth *int `gorm:"type:int"`
	MaxQuantityPerYear  *int `gorm:"type:int"`

	// Age/Gender filters
	MinAge    *int    `gorm:"type:int"`
	MaxAge    *int    `gorm:"type:int"`
	Gender    *string `gorm:"size:10"` // MALE, FEMALE, ALL
	IsForMain bool    `gorm:"default:true"`
	IsForDep  bool    `gorm:"default:true"`

	// Waiting period (دوره انتظار)
	WaitingPeriodDays *int `gorm:"type:int"` // تعداد روز دوره انتظار

	// Special conditions
	NeedsPrescription    bool `gorm:"default:false"` // نیاز به نسخه
	NeedsPreApproval     bool `gorm:"default:false"` // نیاز به تایید قبلی
	NeedsMedicalOpinion  bool `gorm:"default:false"` // نیاز به نظر پزشک
	RequiresDiagnosis    bool `gorm:"default:false"` // نیاز به تشخیص
	AllowedForChronicIll bool `gorm:"default:true"`  // مجاز برای بیماران خاص

	// Priority and status
	Priority int  `gorm:"default:0;index"` // Higher priority = applied first
	IsActive bool `gorm:"default:true;index"`

	// Validity period
	StartDate *time.Time `gorm:"index"`
	EndDate   *time.Time `gorm:"index"`

	// Description
	Title       string `gorm:"size:200"`
	Description string `gorm:"type:text"`
	Note        string `gorm:"type:text"`
}

// TableName specifies the table name for ItemPriceCondition
func (ItemPriceCondition) TableName() string {
	return "item_price_conditions"
}

// IsValid checks if the condition is currently valid
func (ipc *ItemPriceCondition) IsValid() bool {
	if !ipc.IsActive {
		return false
	}

	now := time.Now()

	if ipc.StartDate != nil && now.Before(*ipc.StartDate) {
		return false
	}

	if ipc.EndDate != nil && now.After(*ipc.EndDate) {
		return false
	}

	return true
}

// CalculateCoverage calculates the coverage amount based on the base price
func (ipc *ItemPriceCondition) CalculateCoverage(basePrice int64) int64 {
	// Fixed amount has priority
	if ipc.FixedAmount != nil {
		return *ipc.FixedAmount
	}

	// Calculate percentage-based coverage
	if ipc.CoveragePercentage != nil {
		coverage := int64(float64(basePrice) * (*ipc.CoveragePercentage / 100.0))

		// Apply max limit
		if ipc.MaxCoverageAmount != nil && coverage > *ipc.MaxCoverageAmount {
			coverage = *ipc.MaxCoverageAmount
		}

		// Apply min limit
		if ipc.MinCoverageAmount != nil && coverage < *ipc.MinCoverageAmount {
			coverage = *ipc.MinCoverageAmount
		}

		return coverage
	}

	return 0
}

// CalculateFranchise calculates the franchise amount
func (ipc *ItemPriceCondition) CalculateFranchise(basePrice int64) int64 {
	// Fixed franchise
	if ipc.FranchiseAmount != nil {
		return *ipc.FranchiseAmount
	}

	// Percentage-based franchise
	if ipc.FranchisePercentage != nil {
		franchise := int64(float64(basePrice) * (*ipc.FranchisePercentage / 100.0))

		// Apply max limit
		if ipc.MaxFranchise != nil && franchise > *ipc.MaxFranchise {
			franchise = *ipc.MaxFranchise
		}

		return franchise
	}

	return 0
}

// AppliesTo checks if this condition applies to the given filters
func (ipc *ItemPriceCondition) AppliesTo(itemID *uint, categoryID *uint, groupID *uint, age int, gender string, isMain bool) bool {
	// Check item
	if ipc.ItemID != nil && itemID != nil && *ipc.ItemID != *itemID {
		return false
	}

	// Check category
	if ipc.CategoryID != nil && categoryID != nil && *ipc.CategoryID != *categoryID {
		return false
	}

	// Check group
	if ipc.GroupID != nil && groupID != nil && *ipc.GroupID != *groupID {
		return false
	}

	// Check age
	if ipc.MinAge != nil && age < *ipc.MinAge {
		return false
	}
	if ipc.MaxAge != nil && age > *ipc.MaxAge {
		return false
	}

	// Check gender
	if ipc.Gender != nil && *ipc.Gender != "ALL" && *ipc.Gender != gender {
		return false
	}

	// Check main/dependent
	if isMain && !ipc.IsForMain {
		return false
	}
	if !isMain && !ipc.IsForDep {
		return false
	}

	return true
}
