package entity

// RelationType - نسبت خانوادگی
type RelationType struct {
	BaseModel

	Code        string `gorm:"size:50;uniqueIndex;not null"`
	Title       string `gorm:"size:255;not null"`
	TitleEn     string `gorm:"size:255"`
	Description string `gorm:"type:text"`
	CodeNumber  *int
	IsActive    bool `gorm:"default:true"`

	// Relations
	Employees []Employee `gorm:"foreignKey:RelationTypeID"`
}

// RelationType codes (from stored procedure)
const (
	RelationTypeSpouseFemale = "SPOUSE_FEMALE" // همسر (زن) - ID: 1
	RelationTypeSpouseMale   = "SPOUSE_MALE"   // همسر (مرد) - ID: 2
	RelationTypeChild        = "CHILD"         // فرزند - ID: 3
	RelationTypeDaughter     = "DAUGHTER"      // دختر - ID: 4
	RelationTypeSon          = "SON"           // پسر - ID: 5
	RelationTypeMother       = "MOTHER"        // مادر - ID: 6
	RelationTypeFather       = "FATHER"        // پدر - ID: 7
	RelationTypeSelf         = "SELF"          // خود فرد (کارمند اصلی) - ID: 8
	RelationTypeOther        = "OTHER"         // سایر - ID: 11
)

// TableName specifies the table name
func (RelationType) TableName() string {
	return "relation_types"
}

// IsMainEmployee checks if this is the main employee (self)
func (rt *RelationType) IsMainEmployee() bool {
	return rt.Code == RelationTypeSelf
}

// IsFamilyMember checks if this is a family member
func (rt *RelationType) IsFamilyMember() bool {
	return rt.Code != RelationTypeSelf
}
