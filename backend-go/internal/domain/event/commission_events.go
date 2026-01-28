package event

import (
	"encoding/json"
	"time"
)

// EventVersion defines the current event schema version
const EventVersion = "1.0.0"

// EventType defines event types for commission decisions
type EventType string

const (
	// Commission case events
	EventCommissionCaseCreated    EventType = "commission.case.created"
	EventCommissionCaseAssigned   EventType = "commission.case.assigned"
	EventCommissionCaseReviewed   EventType = "commission.case.reviewed"
	EventCommissionVerdictIssued  EventType = "commission.verdict.issued"
	EventCommissionCaseClosed     EventType = "commission.case.closed"
	EventCommissionCaseAppealed   EventType = "commission.case.appealed"

	// Social work events
	EventSocialWorkCaseCreated    EventType = "socialwork.case.created"
	EventSocialWorkAssessmentDone EventType = "socialwork.assessment.done"
	EventSocialWorkReferralIssued EventType = "socialwork.referral.issued"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	// Event metadata
	EventID     string    `json:"event_id"`      // Unique event identifier (UUID)
	EventType   EventType `json:"event_type"`    // Type of event
	Version     string    `json:"version"`       // Schema version
	Timestamp   time.Time `json:"timestamp"`     // When event occurred
	Source      string    `json:"source"`        // Source service (commission-api, tpa-api)

	// Correlation
	CorrelationID string `json:"correlation_id,omitempty"` // For tracing related events
	CausationID   string `json:"causation_id,omitempty"`   // Event that caused this event

	// Tenant context
	TenantID uint `json:"tenant_id"`
}

// CommissionVerdictEvent is emitted when a medical commission issues a verdict
// This is the primary event consumed by TPA Core for financial decisions
type CommissionVerdictEvent struct {
	BaseEvent

	// Case reference
	CaseID     string `json:"case_id"`      // UUID of commission case
	CaseNumber string `json:"case_number"`  // Human-readable case number

	// Insured person
	InsuredPerson InsuredPersonRef `json:"insured_person"`

	// Verdict details
	Verdict VerdictDetails `json:"verdict"`

	// Financial implications (to be processed by TPA Core)
	FinancialImpact *FinancialImpact `json:"financial_impact,omitempty"`

	// Attachments
	Documents []DocumentRef `json:"documents,omitempty"`
}

// InsuredPersonRef contains reference to insured person
type InsuredPersonRef struct {
	ID            string `json:"id"`             // UUID in commission system
	NationalID    string `json:"national_id"`    // کد ملی
	PersonnelCode string `json:"personnel_code"` // کد پرسنلی
	FullName      string `json:"full_name"`
	Relation      string `json:"relation"`       // main, spouse, child, parent
}

// VerdictDetails contains the commission verdict
type VerdictDetails struct {
	VerdictID       string     `json:"verdict_id"`
	VerdictType     string     `json:"verdict_type"`     // disability, disease, etc.
	VerdictCode     string     `json:"verdict_code"`     // Standard verdict code
	VerdictText     string     `json:"verdict_text"`     // Full verdict text
	DisabilityRate  *int       `json:"disability_rate"`  // درصد از کار افتادگی (0-100)
	EffectiveFrom   *time.Time `json:"effective_from"`
	EffectiveTo     *time.Time `json:"effective_to"`
	IsPermanent     bool       `json:"is_permanent"`
	NeedsReview     bool       `json:"needs_review"`
	ReviewDate      *time.Time `json:"review_date,omitempty"`
	ApprovedBy      string     `json:"approved_by"`      // User who approved
	ApprovedAt      time.Time  `json:"approved_at"`
	CommissionLevel string     `json:"commission_level"` // provincial, central
}

// FinancialImpact describes how this verdict affects coverage
type FinancialImpact struct {
	// Coverage changes
	CoverageType       string  `json:"coverage_type"`        // full, partial, none
	CoveragePercent    int     `json:"coverage_percent"`     // New coverage percentage
	CoverageLimitDelta int64   `json:"coverage_limit_delta"` // Change in coverage limit

	// Benefits
	MonthlyAllowance   int64   `json:"monthly_allowance,omitempty"`
	LumpSumPayment     int64   `json:"lump_sum_payment,omitempty"`

	// Restrictions
	ServiceRestrictions []string `json:"service_restrictions,omitempty"`
	ProviderRestrictions []string `json:"provider_restrictions,omitempty"`

	// Validity
	ValidFrom time.Time  `json:"valid_from"`
	ValidTo   *time.Time `json:"valid_to,omitempty"`
}

// DocumentRef references an attached document
type DocumentRef struct {
	DocumentID   string `json:"document_id"`
	DocumentType string `json:"document_type"` // verdict_letter, medical_report, etc.
	FileName     string `json:"file_name"`
	FileURL      string `json:"file_url"`
	Checksum     string `json:"checksum"` // For integrity verification
}

// SocialWorkReferralEvent is emitted when social work issues a referral
type SocialWorkReferralEvent struct {
	BaseEvent

	// Case reference
	CaseID   string `json:"case_id"`
	CaseType string `json:"case_type"` // financial_aid, loan, equipment, etc.

	// Insured person
	InsuredPerson InsuredPersonRef `json:"insured_person"`

	// Referral details
	Referral ReferralDetails `json:"referral"`
}

// ReferralDetails contains referral information
type ReferralDetails struct {
	ReferralID     string    `json:"referral_id"`
	ReferralType   string    `json:"referral_type"`
	ReferralReason string    `json:"referral_reason"`
	ReferredTo     string    `json:"referred_to"`      // Department/unit
	Priority       string    `json:"priority"`         // normal, urgent, critical
	DueDate        *time.Time `json:"due_date,omitempty"`
	Amount         *int64    `json:"amount,omitempty"` // If financial
	IssuedBy       string    `json:"issued_by"`
	IssuedAt       time.Time `json:"issued_at"`
	Notes          string    `json:"notes,omitempty"`
}

// EventEnvelope wraps events for transport
type EventEnvelope struct {
	Event   interface{} `json:"event"`
	Schema  string      `json:"schema"`  // JSON Schema reference
	Encoded bool        `json:"encoded"` // If event is base64 encoded
}

// ToJSON serializes event to JSON
func (e *CommissionVerdictEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// FromJSON deserializes event from JSON
func (e *CommissionVerdictEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}

// Validate validates event data
func (e *CommissionVerdictEvent) Validate() error {
	if e.EventID == "" {
		return ErrMissingEventID
	}
	if e.TenantID == 0 {
		return ErrMissingTenantID
	}
	if e.CaseID == "" {
		return ErrMissingCaseID
	}
	if e.InsuredPerson.NationalID == "" {
		return ErrMissingNationalID
	}
	return nil
}

// Event validation errors
var (
	ErrMissingEventID    = eventError("event_id is required")
	ErrMissingTenantID   = eventError("tenant_id is required")
	ErrMissingCaseID     = eventError("case_id is required")
	ErrMissingNationalID = eventError("national_id is required")
)

type eventError string

func (e eventError) Error() string { return string(e) }
