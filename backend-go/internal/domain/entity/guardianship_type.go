package entity

// GuardianshipType - نوع کفالت
type GuardianshipType struct {
	BaseModel

	Code    string `gorm:"size:50;uniqueIndex;not null"`
	Title   string `gorm:"size:255;not null"`
	TitleEn string `gorm:"size:255"`

	// Relations
	Employees []Employee `gorm:"foreignKey:GuardianshipTypeID"`
}

// TableName specifies the table name
func (GuardianshipType) TableName() string {
	return "guardianship_types"
}
