package grule

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

// RuleEngine wraps Grule rule engine with versioning and audit
type RuleEngine struct {
	knowledgeLibrary *ast.KnowledgeLibrary
	mu               sync.RWMutex
	auditLogger      AuditLogger
}

// AuditLogger interface for rule execution logging
type AuditLogger interface {
	LogExecution(ctx context.Context, log *ExecutionLog) error
	LogDecision(ctx context.Context, audit *DecisionAudit) error
}

// ExecutionLog represents a single rule execution
type ExecutionLog struct {
	ExecutionID    string
	TenantID       uint
	ClaimID        uint
	ClaimType      string
	RuleID         uint
	RuleCode       string
	RuleVersion    int
	RuleName       string
	Category       string
	Fired          bool
	Condition      string
	ActionTaken    string
	ResultJSON     string
	ExecutionTime  int64
	InputSnapshot  string
	OutputSnapshot string
	HasError       bool
	ErrorMessage   string
	CreatedAt      time.Time
}

// DecisionAudit represents complete claim decision audit
type DecisionAudit struct {
	ExecutionID       string
	TenantID          uint
	ClaimID           uint
	TotalRulesEval    int
	RulesFired        int
	ExecutionTimeMs   int64
	RequestedAmount   int64
	ApprovedAmount    int64
	TotalDeductions   int64
	FranchiseAmount   int64
	BasicInsurShare   int64
	SupplementalShare int64
	DecisionType      string
	RejectionCodes    string
	ClaimSnapshot     string
	RulesSnapshot     string
	ProcessedBy       string
	ProcessingNote    string
	CreatedAt         time.Time
}

// RuleDefinition for loading rules
type RuleDefinition struct {
	ID          uint
	Code        string
	Version     int
	Name        string
	Category    string
	Content     string
	Salience    int
	EffectiveFrom *time.Time
	EffectiveTo   *time.Time
}

// ClaimContext is the data context passed to rules
type ClaimContext struct {
	// Claim info
	ClaimID        uint    `json:"claim_id"`
	ClaimType      string  `json:"claim_type"`
	TenantID       uint    `json:"tenant_id"`

	// Financial
	RequestedAmount int64  `json:"requested_amount"`
	ApprovedAmount  int64  `json:"approved_amount"`
	Deductions      int64  `json:"deductions"`
	Franchise       int64  `json:"franchise"`
	BasicShare      int64  `json:"basic_share"`
	SuppShare       int64  `json:"supp_share"`

	// Service/Drug info
	ServiceCode     string `json:"service_code"`
	ServiceType     string `json:"service_type"`
	DrugCode        string `json:"drug_code"`
	DrugGenericCode string `json:"drug_generic_code"`
	Quantity        int    `json:"quantity"`
	UnitPrice       int64  `json:"unit_price"`

	// Provider info
	ProviderLevel   int    `json:"provider_level"`
	ProviderType    string `json:"provider_type"`

	// Patient info
	PatientAge      int    `json:"patient_age"`
	PatientGender   string `json:"patient_gender"`
	RelationType    string `json:"relation_type"`

	// Coverage info
	CoveragePercent int    `json:"coverage_percent"`
	CoverageLimit   int64  `json:"coverage_limit"`
	UsedAmount      int64  `json:"used_amount"`
	RemainingLimit  int64  `json:"remaining_limit"`

	// Flags
	IsEmergency     bool   `json:"is_emergency"`
	NeedsPreAuth    bool   `json:"needs_pre_auth"`
	HasPreAuth      bool   `json:"has_pre_auth"`

	// Results (set by rules)
	IsApproved      bool     `json:"is_approved"`
	RejectionCodes  []string `json:"rejection_codes"`
	DeductionCodes  []string `json:"deduction_codes"`
	Notes           []string `json:"notes"`
}

// NewRuleEngine creates a new rule engine instance
func NewRuleEngine(auditLogger AuditLogger) *RuleEngine {
	return &RuleEngine{
		knowledgeLibrary: ast.NewKnowledgeLibrary(),
		auditLogger:      auditLogger,
	}
}

// LoadRules loads rules into the knowledge base
func (e *RuleEngine) LoadRules(ctx context.Context, tenantID uint, rules []RuleDefinition) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	kb := fmt.Sprintf("tenant_%d", tenantID)
	version := fmt.Sprintf("v%d", time.Now().Unix())

	// Build GRL content
	var grlContent string
	for _, rule := range rules {
		if !e.isRuleEffective(rule) {
			continue
		}
		grlContent += rule.Content + "\n\n"
	}

	if grlContent == "" {
		return nil // No effective rules
	}

	// Build knowledge base
	ruleBuilder := builder.NewRuleBuilder(e.knowledgeLibrary)
	bs := pkg.NewBytesResource([]byte(grlContent))
	err := ruleBuilder.BuildRuleFromResource(kb, version, bs)
	if err != nil {
		return fmt.Errorf("failed to build rules: %w", err)
	}

	return nil
}

// isRuleEffective checks if rule is within effective date range
func (e *RuleEngine) isRuleEffective(rule RuleDefinition) bool {
	now := time.Now()
	if rule.EffectiveFrom != nil && now.Before(*rule.EffectiveFrom) {
		return false
	}
	if rule.EffectiveTo != nil && now.After(*rule.EffectiveTo) {
		return false
	}
	return true
}

// Evaluate runs rules against a claim context
func (e *RuleEngine) Evaluate(ctx context.Context, tenantID uint, claimCtx *ClaimContext, rules []RuleDefinition) (*EvaluationResult, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	executionID := uuid.New().String()
	startTime := time.Now()

	// Create data context
	dataCtx := ast.NewDataContext()
	if err := dataCtx.Add("Claim", claimCtx); err != nil {
		return nil, fmt.Errorf("failed to add claim context: %w", err)
	}

	// Get knowledge base
	kb := fmt.Sprintf("tenant_%d", tenantID)
	version := fmt.Sprintf("v%d", time.Now().Unix())

	knowledgeBase, err := e.knowledgeLibrary.NewKnowledgeBaseInstance(kb, version)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}

	// Execute rules
	eng := engine.NewGruleEngine()
	err = eng.Execute(dataCtx, knowledgeBase)

	executionTime := time.Since(startTime)

	// Build result
	result := &EvaluationResult{
		ExecutionID:     executionID,
		TenantID:        tenantID,
		ClaimID:         claimCtx.ClaimID,
		ExecutionTimeMs: executionTime.Milliseconds(),
		RulesEvaluated:  len(rules),
		Context:         claimCtx,
	}

	if err != nil {
		result.HasError = true
		result.ErrorMessage = err.Error()
	}

	// Log audit
	if e.auditLogger != nil {
		inputJSON, _ := json.Marshal(claimCtx)
		outputJSON, _ := json.Marshal(result)

		audit := &DecisionAudit{
			ExecutionID:       executionID,
			TenantID:          tenantID,
			ClaimID:           claimCtx.ClaimID,
			TotalRulesEval:    len(rules),
			ExecutionTimeMs:   executionTime.Milliseconds(),
			RequestedAmount:   claimCtx.RequestedAmount,
			ApprovedAmount:    claimCtx.ApprovedAmount,
			TotalDeductions:   claimCtx.Deductions,
			FranchiseAmount:   claimCtx.Franchise,
			BasicInsurShare:   claimCtx.BasicShare,
			SupplementalShare: claimCtx.SuppShare,
			ClaimSnapshot:     string(inputJSON),
			RulesSnapshot:     e.buildRulesSnapshot(rules),
			ProcessedBy:       "system",
			CreatedAt:         time.Now(),
		}

		if claimCtx.IsApproved {
			audit.DecisionType = "approved"
		} else if len(claimCtx.RejectionCodes) > 0 {
			audit.DecisionType = "rejected"
			codesJSON, _ := json.Marshal(claimCtx.RejectionCodes)
			audit.RejectionCodes = string(codesJSON)
		} else {
			audit.DecisionType = "partial"
		}

		_ = e.auditLogger.LogDecision(ctx, audit)
	}

	return result, nil
}

// EvaluationResult contains the result of rule evaluation
type EvaluationResult struct {
	ExecutionID     string        `json:"execution_id"`
	TenantID        uint          `json:"tenant_id"`
	ClaimID         uint          `json:"claim_id"`
	ExecutionTimeMs int64         `json:"execution_time_ms"`
	RulesEvaluated  int           `json:"rules_evaluated"`
	RulesFired      int           `json:"rules_fired"`
	HasError        bool          `json:"has_error"`
	ErrorMessage    string        `json:"error_message,omitempty"`
	Context         *ClaimContext `json:"context"`
}

// buildRulesSnapshot creates a snapshot of rules used
func (e *RuleEngine) buildRulesSnapshot(rules []RuleDefinition) string {
	snapshot := make([]map[string]interface{}, len(rules))
	for i, r := range rules {
		snapshot[i] = map[string]interface{}{
			"id":       r.ID,
			"code":     r.Code,
			"version":  r.Version,
			"checksum": e.hashContent(r.Content),
		}
	}
	data, _ := json.Marshal(snapshot)
	return string(data)
}

// hashContent generates SHA256 hash of content
func (e *RuleEngine) hashContent(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}

// ReplayDecision replays a claim decision using historical snapshot
func (e *RuleEngine) ReplayDecision(ctx context.Context, audit *DecisionAudit) (*EvaluationResult, error) {
	// Restore claim context from snapshot
	var claimCtx ClaimContext
	if err := json.Unmarshal([]byte(audit.ClaimSnapshot), &claimCtx); err != nil {
		return nil, fmt.Errorf("failed to restore claim snapshot: %w", err)
	}

	// Restore rules from snapshot
	var rulesSnapshot []map[string]interface{}
	if err := json.Unmarshal([]byte(audit.RulesSnapshot), &rulesSnapshot); err != nil {
		return nil, fmt.Errorf("failed to restore rules snapshot: %w", err)
	}

	// Note: Full replay would require loading historical rule versions from DB
	// This is a simplified version that re-evaluates with current rules
	return e.Evaluate(ctx, audit.TenantID, &claimCtx, nil)
}

// Example GRL rules for reference
const ExampleCoverageRule = `
rule CheckCoverageLimit "بررسی سقف پوشش" salience 10 {
    when
        Claim.RequestedAmount > Claim.RemainingLimit
    then
        Claim.ApprovedAmount = Claim.RemainingLimit;
        Claim.Deductions = Claim.RequestedAmount - Claim.RemainingLimit;
        Claim.DeductionCodes = Append(Claim.DeductionCodes, "LIMIT_EXCEEDED");
        Retract("CheckCoverageLimit");
}
`

const ExampleFranchiseRule = `
rule CalculateFranchise "محاسبه فرانشیز" salience 5 {
    when
        Claim.ApprovedAmount > 0 && Claim.CoveragePercent > 0
    then
        Claim.Franchise = Claim.ApprovedAmount * (100 - Claim.CoveragePercent) / 100;
        Claim.SuppShare = Claim.ApprovedAmount - Claim.Franchise;
        Retract("CalculateFranchise");
}
`

const ExampleEmergencyRule = `
rule EmergencyBonus "پوشش اورژانس" salience 15 {
    when
        Claim.IsEmergency == true && Claim.CoveragePercent < 100
    then
        Claim.CoveragePercent = Claim.CoveragePercent + 10;
        Claim.Notes = Append(Claim.Notes, "Emergency coverage bonus applied");
        Retract("EmergencyBonus");
}
`
