package entity

import "time"

// InsuranceHistory represents the history of insurance coverage for an employee
// Tracks all insurance policies, changes, suspensions, and renewals over time
type InsuranceHistory struct {
	BaseModel
	TenantID uint `gorm:"not null;index:idx_ins_history_tenant"`

	// Employee reference
	EmployeeID uint     `gorm:"not null;index:idx_ins_history_employee"`
	Employee   Employee `gorm:"foreignKey:EmployeeID"`

	// Insurance reference
	InsuranceID uint      `gorm:"not null;index:idx_ins_history_insurance"`
	Insurance   Insurance `gorm:"foreignKey:InsuranceID"`

	// Policy reference
	PolicyID *uint   `gorm:"index"`
	Policy   *Policy `gorm:"foreignKey:PolicyID"`

	// History event type
	EventType string `gorm:"size:50;not null;index"` // START, RENEWAL, CHANGE, SUSPENSION, TERMINATION, REINSTATEMENT

	// Period
	StartDate time.Time `gorm:"not null;index"` // تاریخ شروع
	EndDate   *time.Time `gorm:"index"`          // تاریخ پایان

	// Coverage details snapshot
	CoverageType       *string `gorm:"size:50"`        // INDIVIDUAL, FAMILY, GROUP
	PlanName           *string `gorm:"size:200"`       // نام طرح
	PlanCode           *string `gorm:"size:50;index"`  // کد طرح
	AnnualLimit        *int64  `gorm:"type:bigint"`    // سقف سالانه در آن زمان
	FranchiseAmount    *int64  `gorm:"type:bigint"`    // فرانشیز
	FranchisePercentage *float64 `gorm:"type:decimal(5,2)"` // درصد فرانشیز
	PremiumAmount      *int64  `gorm:"type:bigint"`    // حق بیمه

	// Change details (for CHANGE events)
	ChangeReason       *string `gorm:"type:text"`  // دلیل تغییر
	ChangedFields      *string `gorm:"type:text"`  // فیلدهای تغییر یافته (JSON)
	PreviousValues     *string `gorm:"type:text"`  // مقادیر قبلی (JSON)
	NewValues          *string `gorm:"type:text"`  // مقادیر جدید (JSON)

	// Suspension details (for SUSPENSION events)
	SuspensionReason   *string `gorm:"type:text"`  // دلیل تعلیق
	SuspendedAt        *time.Time                  // تاریخ تعلیق
	SuspendedBy        *uint                       // کاربر تعلیق کننده

	// Termination details (for TERMINATION events)
	TerminationReason  *string `gorm:"type:text"`  // دلیل خاتمه
	TerminatedAt       *time.Time                  // تاریخ خاتمه
	TerminatedBy       *uint                       // کاربر خاتمه دهنده
	RefundAmount       *int64  `gorm:"type:bigint"` // مبلغ استرداد

	// Renewal details (for RENEWAL events)
	RenewalDate        *time.Time                  // تاریخ تمدید
	RenewalType        *string `gorm:"size:50"`    // AUTO, MANUAL, FORCED
	PremiumChangePercent *float64 `gorm:"type:decimal(5,2)"` // درصد تغییر حق بیمه

	// Usage statistics during this period
	TotalClaims        int   `gorm:"default:0"`          // تعداد ادعاها
	TotalClaimAmount   int64 `gorm:"type:bigint;default:0"` // مجموع مبلغ ادعاها
	TotalPaidAmount    int64 `gorm:"type:bigint;default:0"` // مجموع پرداختی
	LossRatio          *float64 `gorm:"type:decimal(5,2)"` // ضریب خسارت

	// Documents
	ContractPath       *string `gorm:"size:500"`   // مسیر قرارداد
	AddendumPath       *string `gorm:"size:500"`   // مسیر الحاقیه
	TerminationLetterPath *string `gorm:"size:500"` // مسیر نامه خاتمه
	DocumentsPath      *string `gorm:"type:text"`  // سایر مدارک (JSON)

	// Approval workflow
	RequiresApproval   bool       `gorm:"default:false"`
	ApprovedBy         *uint                      // کاربر تایید کننده
	ApprovedAt         *time.Time                 // تاریخ تایید
	ApprovalNotes      *string `gorm:"type:text"` // یادداشت تایید
	RejectedBy         *uint                      // کاربر رد کننده
	RejectedAt         *time.Time                 // تاریخ رد
	RejectionReason    *string `gorm:"type:text"` // دلیل رد

	// Notes
	Description        string  `gorm:"type:text"`     // توضیحات
	Notes              *string `gorm:"type:text"`     // یادداشت‌ها
	InternalNotes      *string `gorm:"type:text"`     // یادداشت داخلی

	// System tracking
	CreatedBy          *uint                      // کاربر ثبت کننده
	IsSystemGenerated  bool `gorm:"default:false"` // ثبت خودکار توسط سیستم
	SourceSystem       *string `gorm:"size:50"`   // منبع ثبت (MANUAL, AUTO, IMPORT, etc.)

	// Status
	IsActive           bool `gorm:"default:true;index"`
	IsCurrent          bool `gorm:"default:false;index"` // آیا فعلی است؟
}

// TableName specifies the table name for InsuranceHistory
func (InsuranceHistory) TableName() string {
	return "insurance_histories"
}

// IsPeriodActive checks if this history period is currently active
func (ih *InsuranceHistory) IsPeriodActive() bool {
	if !ih.IsActive {
		return false
	}

	now := time.Now()

	if now.Before(ih.StartDate) {
		return false
	}

	if ih.EndDate != nil && now.After(*ih.EndDate) {
		return false
	}

	return true
}

// GetDuration returns the duration of this history period in days
func (ih *InsuranceHistory) GetDuration() int {
	if ih.EndDate == nil {
		// If no end date, calculate from start to now
		return int(time.Since(ih.StartDate).Hours() / 24)
	}

	return int(ih.EndDate.Sub(ih.StartDate).Hours() / 24)
}

// CalculateLossRatio calculates the loss ratio (if not already set)
func (ih *InsuranceHistory) CalculateLossRatio() float64 {
	if ih.PremiumAmount == nil || *ih.PremiumAmount == 0 {
		return 0
	}

	return (float64(ih.TotalPaidAmount) / float64(*ih.PremiumAmount)) * 100
}
