package entity

// ConditionGroup represents a group of pricing conditions
// Used to apply multiple conditions together as a logical group (AND/OR)
type ConditionGroup struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_cond_group_tenant"`

	// Group information
	Code        string `gorm:"size:50;uniqueIndex:idx_cond_group_code_tenant;not null"`
	Title       string `gorm:"size:200;not null"`
	Description string `gorm:"type:text"`

	// Logical operator
	LogicOperator string `gorm:"size:10;default:'AND'"` // AND, OR

	// Priority
	Priority int  `gorm:"default:0;index"`
	IsActive bool `gorm:"default:true;index"`

	// Relations
	PriceConditions []ItemPriceCondition `gorm:"many2many:condition_group_mappings;"`
}

// TableName specifies the table name for ConditionGroup
func (ConditionGroup) TableName() string {
	return "condition_groups"
}

// ConditionGroupMapping represents the many-to-many relationship
type ConditionGroupMapping struct {
	BaseModel

	ConditionGroupID     uint                `gorm:"not null;index"`
	ConditionGroup       ConditionGroup      `gorm:"foreignKey:ConditionGroupID"`
	ItemPriceConditionID uint                `gorm:"not null;index"`
	ItemPriceCondition   ItemPriceCondition  `gorm:"foreignKey:ItemPriceConditionID"`

	// Mapping properties
	SortOrder int  `gorm:"default:0"`
	IsActive  bool `gorm:"default:true"`
}

// TableName specifies the table name for ConditionGroupMapping
func (ConditionGroupMapping) TableName() string {
	return "condition_group_mappings"
}
