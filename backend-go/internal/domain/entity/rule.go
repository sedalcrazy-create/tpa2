package entity

import (
	"time"

	"gorm.io/gorm"
)

// RuleCategory defines rule categories
type RuleCategory string

const (
	RuleCategoryCoverage   RuleCategory = "coverage"   // پوشش بیمه‌ای
	RuleCategoryFranchise  RuleCategory = "franchise"  // فرانشیز
	RuleCategoryDeduction  RuleCategory = "deduction"  // کسورات
	RuleCategoryLimit      RuleCategory = "limit"      // سقف تعهد
	RuleCategoryValidation RuleCategory = "validation" // اعتبارسنجی
)

// RuleStatus defines rule lifecycle status
type RuleStatus string

const (
	RuleStatusDraft     RuleStatus = "draft"     // پیش‌نویس
	RuleStatusActive    RuleStatus = "active"    // فعال
	RuleStatusInactive  RuleStatus = "inactive"  // غیرفعال
	RuleStatusArchived  RuleStatus = "archived"  // بایگانی
	RuleStatusDeprecated RuleStatus = "deprecated" // منسوخ
)

// Rule represents a business rule definition
type Rule struct {
	BaseModel
	TenantID uint `gorm:"index;not null" json:"tenant_id"`

	// Rule identification
	Code        string       `gorm:"size:50;uniqueIndex:idx_rule_code_version" json:"code"`
	Version     int          `gorm:"uniqueIndex:idx_rule_code_version;default:1" json:"version"`
	Name        string       `gorm:"size:200" json:"name"`
	NameFa      string       `gorm:"size:200" json:"name_fa"` // نام فارسی
	Description string       `gorm:"size:2000" json:"description"`
	Category    RuleCategory `gorm:"size:50;index" json:"category"`
	Status      RuleStatus   `gorm:"size:50;index;default:'draft'" json:"status"`

	// Rule content (GRL format for Grule)
	RuleContent string `gorm:"type:text;not null" json:"rule_content"`
	Salience    int    `gorm:"default:0" json:"salience"` // Priority (higher = first)

	// Applicability
	ClaimTypes     string `gorm:"size:500" json:"claim_types"`      // Comma-separated claim types
	ServiceTypes   string `gorm:"size:500" json:"service_types"`    // Comma-separated service types
	ProviderLevels string `gorm:"size:200" json:"provider_levels"`  // Comma-separated provider levels

	// Validity period
	EffectiveFrom *time.Time `gorm:"index" json:"effective_from"`
	EffectiveTo   *time.Time `gorm:"index" json:"effective_to"`

	// Versioning
	PreviousVersionID *uint `json:"previous_version_id"`
	IsLatest          bool  `gorm:"default:true" json:"is_latest"`

	// Audit
	CreatedBy   uint       `json:"created_by"`
	ApprovedBy  *uint      `json:"approved_by"`
	ApprovedAt  *time.Time `json:"approved_at"`
	PublishedAt *time.Time `json:"published_at"`

	// Metadata
	Tags     string `gorm:"size:500" json:"tags"` // Comma-separated tags
	Notes    string `gorm:"size:2000" json:"notes"`
	Checksum string `gorm:"size:64" json:"checksum"` // SHA256 of rule content
}

// RuleExecutionLog records each rule execution for audit
type RuleExecutionLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	// Context
	TenantID  uint   `gorm:"index;not null" json:"tenant_id"`
	ClaimID   uint   `gorm:"index" json:"claim_id"`
	ClaimType string `gorm:"size:50" json:"claim_type"`

	// Execution batch (groups all rules fired for one claim evaluation)
	ExecutionID string `gorm:"size:36;index" json:"execution_id"` // UUID

	// Rule info (denormalized for immutability)
	RuleID      uint         `gorm:"index" json:"rule_id"`
	RuleCode    string       `gorm:"size:50" json:"rule_code"`
	RuleVersion int          `json:"rule_version"`
	RuleName    string       `gorm:"size:200" json:"rule_name"`
	Category    RuleCategory `gorm:"size:50" json:"category"`

	// Execution result
	Fired         bool   `json:"fired"`              // Rule was triggered
	Condition     string `gorm:"size:1000" json:"condition"` // Rule condition that matched
	ActionTaken   string `gorm:"size:500" json:"action_taken"`
	ResultJSON    string `gorm:"type:text" json:"result_json"` // Structured result
	ExecutionTime int64  `json:"execution_time_ns"`            // Execution time in nanoseconds

	// Input/Output snapshots for replay
	InputSnapshot  string `gorm:"type:text" json:"input_snapshot"`  // JSON of input data
	OutputSnapshot string `gorm:"type:text" json:"output_snapshot"` // JSON of output data

	// Error handling
	HasError     bool   `json:"has_error"`
	ErrorMessage string `gorm:"size:2000" json:"error_message"`
}

// RuleSet represents a collection of rules that work together
type RuleSet struct {
	BaseModel
	TenantID uint `gorm:"index;not null" json:"tenant_id"`

	Name        string `gorm:"size:200" json:"name"`
	NameFa      string `gorm:"size:200" json:"name_fa"`
	Description string `gorm:"size:2000" json:"description"`
	Version     int    `gorm:"default:1" json:"version"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`

	// Rules in this set
	Rules []Rule `gorm:"many2many:rule_set_rules;" json:"rules"`
}

// ClaimDecisionAudit provides complete audit trail for claim decisions
type ClaimDecisionAudit struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`

	// Claim reference
	TenantID uint `gorm:"index;not null" json:"tenant_id"`
	ClaimID  uint `gorm:"index;not null" json:"claim_id"`

	// Decision summary
	ExecutionID      string `gorm:"size:36;uniqueIndex" json:"execution_id"` // Links to RuleExecutionLog
	TotalRulesEval   int    `json:"total_rules_evaluated"`
	RulesFired       int    `json:"rules_fired"`
	ExecutionTimeMs  int64  `json:"execution_time_ms"`

	// Financial impact
	RequestedAmount    int64 `json:"requested_amount"`
	ApprovedAmount     int64 `json:"approved_amount"`
	TotalDeductions    int64 `json:"total_deductions"`
	FranchiseAmount    int64 `json:"franchise_amount"`
	BasicInsurShare    int64 `json:"basic_insur_share"`
	SupplementalShare  int64 `json:"supplemental_share"`

	// Decision details
	DecisionType   string `gorm:"size:50" json:"decision_type"` // approved, rejected, partial
	RejectionCodes string `gorm:"size:500" json:"rejection_codes"`

	// Snapshots for replay
	ClaimSnapshot string `gorm:"type:text" json:"claim_snapshot"` // Full claim JSON at decision time
	RulesSnapshot string `gorm:"type:text" json:"rules_snapshot"` // Rules used (IDs + versions)

	// Audit metadata
	ProcessedBy    string `gorm:"size:100" json:"processed_by"` // system or user_id
	ProcessingNote string `gorm:"size:1000" json:"processing_note"`
}

// TableName returns table name
func (Rule) TableName() string {
	return "rules"
}

func (RuleExecutionLog) TableName() string {
	return "rule_execution_logs"
}

func (RuleSet) TableName() string {
	return "rule_sets"
}

func (ClaimDecisionAudit) TableName() string {
	return "claim_decision_audits"
}

// BeforeCreate hook for Rule
func (r *Rule) BeforeCreate(tx *gorm.DB) error {
	// Generate checksum
	if r.RuleContent != "" {
		r.Checksum = generateSHA256(r.RuleContent)
	}
	return nil
}

// generateSHA256 creates checksum for rule content
func generateSHA256(content string) string {
	// Implementation would use crypto/sha256
	// Simplified here - actual implementation in service layer
	return ""
}

// IsEffective checks if rule is currently effective
func (r *Rule) IsEffective() bool {
	now := time.Now()
	if r.Status != RuleStatusActive {
		return false
	}
	if r.EffectiveFrom != nil && now.Before(*r.EffectiveFrom) {
		return false
	}
	if r.EffectiveTo != nil && now.After(*r.EffectiveTo) {
		return false
	}
	return true
}
