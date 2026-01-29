package entity

// SpecialEmployeeType - گروه‌های ایثارگری (از stored procedure)
type SpecialEmployeeType struct {
	BaseModel

	Code               string  `gorm:"size:50;uniqueIndex;not null"`
	Title              string  `gorm:"size:255;not null"`
	TitleEn            string  `gorm:"size:255"`
	Description        string  `gorm:"type:text"`
	Priority           int     `gorm:"default:0"`
	DiscountPercentage *float64
	IsActive           bool `gorm:"default:true"`

	// Relations
	Employees []Employee `gorm:"foreignKey:SpecialEmployeeTypeID"`
}

// SpecialEmployeeType codes (from stored procedure id_set logic)
const (
	SpecialTypeJanbaazCombined = "JANBAAZ_COMBINED" // جانباز / رزمنده / ترکیبی - ID: 1
	SpecialTypeAzadeh          = "AZADEH"           // آزاده - ID: 2
	SpecialTypeShahidChild50   = "SHAHID_CHILD_50"  // فرزند شاهد (50% جانبازی) - ID: 3
)

// TableName specifies the table name
func (SpecialEmployeeType) TableName() string {
	return "special_employee_types"
}

// GenerateTypeCode generates employee type code based on Refah/Yii logic
// Formula: id_set*1000 + (isRetired?100:200) + id_cec
func (set *SpecialEmployeeType) GenerateTypeCode(isRetired bool, cecID uint) int {
	code := set.ID * 1000

	if isRetired {
		code += 100
	} else {
		code += 200
	}

	code += int(cecID)
	return code
}
