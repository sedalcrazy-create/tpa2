package entity

// Instruction represents drug/service usage instructions
// Used for prescription items to specify how to use the medication/service
type Instruction struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_instruction_tenant"`

	// Code and naming
	Code    string `gorm:"size:50;uniqueIndex:idx_instruction_code_tenant;not null"`
	Title   string `gorm:"size:200;not null"`
	TitleFa string `gorm:"size:200;not null;index"`

	// Instruction details
	Frequency   string `gorm:"size:100"` // e.g., "روزی 3 بار", "هر 8 ساعت یکبار"
	Duration    string `gorm:"size:100"` // e.g., "7 روز", "تا بهبودی"
	Timing      string `gorm:"size:100"` // e.g., "قبل از غذا", "بعد از غذا", "با غذا"
	Route       string `gorm:"size:100"` // e.g., "خوراکی", "تزریقی", "موضعی"
	Dosage      string `gorm:"size:100"` // e.g., "1 قرص", "2 قاشق چایخوری"
	Description string `gorm:"type:text"`

	// Templates for common instructions
	IsTemplate bool   `gorm:"default:false;index"`
	Template   string `gorm:"type:text"` // JSON template for complex instructions

	// Category
	Category string `gorm:"size:50;index"` // DRUG, THERAPY, TEST, etc.

	// Special flags
	RequiresSupervision bool `gorm:"default:false"` // نیاز به نظارت
	IsEmergency         bool `gorm:"default:false"` // اورژانسی
	IsChronic           bool `gorm:"default:false"` // برای بیماری مزمن

	// Status
	IsActive  bool `gorm:"default:true;index"`
	SortOrder int  `gorm:"default:0"`

	// Relations
	PrescriptionItems []PrescriptionItem `gorm:"foreignKey:InstructionID"`
}

// TableName specifies the table name for Instruction
func (Instruction) TableName() string {
	return "instructions"
}

// GetFullInstruction returns the complete instruction text
func (i *Instruction) GetFullInstruction() string {
	instruction := ""

	if i.Dosage != "" {
		instruction += i.Dosage
	}

	if i.Frequency != "" {
		if instruction != "" {
			instruction += " - "
		}
		instruction += i.Frequency
	}

	if i.Timing != "" {
		if instruction != "" {
			instruction += " - "
		}
		instruction += i.Timing
	}

	if i.Duration != "" {
		if instruction != "" {
			instruction += " - "
		}
		instruction += "به مدت " + i.Duration
	}

	if i.Route != "" {
		if instruction != "" {
			instruction += " - "
		}
		instruction += i.Route
	}

	return instruction
}
