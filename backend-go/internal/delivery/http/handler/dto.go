package handler

import (
	"time"

	"github.com/bank-melli/tpa/internal/domain/entity"
)

// CreateClaimRequest represents claim creation request
type CreateClaimRequest struct {
	TrackingCode   string              `json:"tracking_code" validate:"required"`
	HID            string              `json:"hid"`
	MRN            string              `json:"mrn"`
	PolicyMemberID uint                `json:"policy_member_id" validate:"required"`
	CenterID       uint                `json:"center_id" validate:"required"`
	ClaimType      entity.ClaimType    `json:"claim_type" validate:"required"`
	AdmissionType  entity.AdmissionType `json:"admission_type"`
	AdmissionDate  time.Time           `json:"admission_date" validate:"required"`
	DischargeDate  *time.Time          `json:"discharge_date"`
	ServiceDate    time.Time           `json:"service_date" validate:"required"`
	RequestAmount  int64               `json:"request_amount"`
	BasicInsShare  int64               `json:"basic_ins_share"`
}

// UpdateClaimRequest represents claim update request
type UpdateClaimRequest struct {
	RequestAmount *int64 `json:"request_amount"`
	BasicInsShare *int64 `json:"basic_ins_share"`
}

// CompleteExaminationRequest represents examination completion request
type CompleteExaminationRequest struct {
	ApprovedAmount  int64                       `json:"approved_amount"`
	Deduction       int64                       `json:"deduction"`
	DeductionReason string                      `json:"deduction_reason"`
	Notes           string                      `json:"notes"`
	Items           []ItemExaminationRequestDTO `json:"items"`
}

// ItemExaminationRequestDTO represents item examination request
type ItemExaminationRequestDTO struct {
	ItemID         uint   `json:"item_id"`
	ConfirmedPrice int64  `json:"confirmed_price"`
	Deduction      int64  `json:"deduction"`
	ReasonCodeIDs  []uint `json:"reason_code_ids"`
	Notes          string `json:"notes"`
}

// RejectRequest represents rejection request
type RejectRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// CreatePackageRequest represents package creation request
type CreatePackageRequest struct {
	CenterID          uint       `json:"center_id" validate:"required"`
	WorkUnitID        *uint      `json:"work_unit_id"`
	Title             string     `json:"title" validate:"required"`
	LetterNumber      string     `json:"letter_number"`
	LetterDate        *time.Time `json:"letter_date"`
	ReceiveLetterDate *time.Time `json:"receive_letter_date"`
	LetterImageURL    string     `json:"letter_image_url"`
	ClaimIDs          []uint     `json:"claim_ids"`
}

// UpdatePackageRequest represents package update request
type UpdatePackageRequest struct {
	Title             *string    `json:"title"`
	LetterNumber      *string    `json:"letter_number"`
	LetterDate        *time.Time `json:"letter_date"`
	ReceiveLetterDate *time.Time `json:"receive_letter_date"`
}

// CreateCenterRequest represents center creation request
type CreateCenterRequest struct {
	Title          string             `json:"title" validate:"required"`
	SiamID         string             `json:"siam_id" validate:"required"`
	Code           string             `json:"code"`
	Type           entity.CenterType  `json:"type" validate:"required"`
	Level          int                `json:"level"`
	ProvinceID     *uint              `json:"province_id"`
	CityID         *uint              `json:"city_id"`
	Address        string             `json:"address"`
	PostalCode     string             `json:"postal_code"`
	Phone          string             `json:"phone"`
	Fax            string             `json:"fax"`
	Email          string             `json:"email"`
	Website        string             `json:"website"`
	OwnerName      string             `json:"owner_name"`
	ManagerName    string             `json:"manager_name"`
	ManagerPhone   string             `json:"manager_phone"`
	DependencyType uint8              `json:"dependency_type"`
	PaymentID      string             `json:"payment_id"`
	AccountNumber  string             `json:"account_number"`
	ShebaNumber    string             `json:"sheba_number"`
	EconomicCode   string             `json:"economic_code"`
	NationalID     string             `json:"national_id"`
}

// CreateSettlementRequest represents settlement creation request
type CreateSettlementRequest struct {
	CenterID    uint       `json:"center_id" validate:"required"`
	PeriodStart time.Time  `json:"period_start" validate:"required"`
	PeriodEnd   time.Time  `json:"period_end" validate:"required"`
	PackageIDs  []uint     `json:"package_ids"`
	Notes       string     `json:"notes"`
}

// LoginRequest represents login request
type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents login response
type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	ExpiresIn    int64       `json:"expires_in"`
	User         UserInfoDTO `json:"user"`
}

// UserInfoDTO represents user info in responses
type UserInfoDTO struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	RoleName    string `json:"role_name"`
	RoleTitle   string `json:"role_title"`
	CenterID    *uint  `json:"center_id,omitempty"`
	CenterTitle string `json:"center_title,omitempty"`
	TenantID    uint   `json:"tenant_id"`
}

// RefreshTokenRequest represents refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// CreateUserRequest represents user creation request
type CreateUserRequest struct {
	Username    string `json:"username" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	FirstName   string `json:"first_name" validate:"required"`
	LastName    string `json:"last_name" validate:"required"`
	Mobile      string `json:"mobile"`
	RoleID      uint   `json:"role_id" validate:"required"`
	CenterID    *uint  `json:"center_id"`
	WorkUnitID  *uint  `json:"work_unit_id"`
	ProvinceID  *uint  `json:"province_id"`
}

// UpdateUserRequest represents user update request
type UpdateUserRequest struct {
	Email      *string `json:"email"`
	FirstName  *string `json:"first_name"`
	LastName   *string `json:"last_name"`
	Mobile     *string `json:"mobile"`
	RoleID     *uint   `json:"role_id"`
	CenterID   *uint   `json:"center_id"`
	WorkUnitID *uint   `json:"work_unit_id"`
	ProvinceID *uint   `json:"province_id"`
	IsActive   *bool   `json:"is_active"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// CreatePersonRequest represents person creation request
type CreatePersonRequest struct {
	NationalCode  string     `json:"national_code" validate:"required,len=10"`
	FirstName     string     `json:"first_name" validate:"required"`
	LastName      string     `json:"last_name" validate:"required"`
	FatherName    string     `json:"father_name"`
	BirthCertNo   string     `json:"birth_cert_no"`
	BirthDate     *time.Time `json:"birth_date"`
	BirthPlaceID  *uint      `json:"birth_place_id"`
	IssuePlace    string     `json:"issue_place"`
	Gender        entity.Gender `json:"gender"`
	MaritalStatus uint8      `json:"marital_status"`
	Mobile        string     `json:"mobile"`
	Phone         string     `json:"phone"`
	Email         string     `json:"email"`
	Address       string     `json:"address"`
	PostalCode    string     `json:"postal_code"`
	AccountNumber string     `json:"account_number"`
	ShebaNumber   string     `json:"sheba_number"`
	CardNumber    string     `json:"card_number"`
}

// PolicyMemberInquiryRequest represents policy member inquiry request
type PolicyMemberInquiryRequest struct {
	NationalCode string `json:"national_code" validate:"required,len=10"`
	ServiceDate  string `json:"service_date"` // YYYY-MM-DD
}

// PolicyMemberInquiryResponse represents policy member inquiry response
type PolicyMemberInquiryResponse struct {
	IsEligible       bool               `json:"is_eligible"`
	Message          string             `json:"message,omitempty"`
	PolicyMember     *PolicyMemberDTO   `json:"policy_member,omitempty"`
	Supervisor       *PolicyMemberDTO   `json:"supervisor,omitempty"`
	ActivePolicies   []PolicyDTO        `json:"active_policies,omitempty"`
}

// PolicyMemberDTO represents policy member data
type PolicyMemberDTO struct {
	ID             uint   `json:"id"`
	NationalCode   string `json:"national_code"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FatherName     string `json:"father_name"`
	BirthDate      string `json:"birth_date,omitempty"`
	Gender         string `json:"gender"`
	DependencyType string `json:"dependency_type"`
}

// PolicyDTO represents policy data
type PolicyDTO struct {
	ID           uint   `json:"id"`
	PolicyNumber string `json:"policy_number"`
	InsurerName  string `json:"insurer_name"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Status       string `json:"status"`
}

// DrugSearchRequest represents drug search request
type DrugSearchRequest struct {
	Query string `json:"query" validate:"required,min=2"`
	Limit int    `json:"limit"`
}

// ServiceSearchRequest represents service search request
type ServiceSearchRequest struct {
	Query string `json:"query" validate:"required,min=2"`
	Limit int    `json:"limit"`
}

// DiagnosisSearchRequest represents diagnosis search request
type DiagnosisSearchRequest struct {
	Query string `json:"query" validate:"required,min=2"`
	Limit int    `json:"limit"`
}

// DashboardStatsResponse represents dashboard statistics
type DashboardStatsResponse struct {
	Claims    ClaimStatsDTO    `json:"claims"`
	Packages  PackageStatsDTO  `json:"packages"`
	Amounts   AmountStatsDTO   `json:"amounts"`
	Centers   CenterStatsDTO   `json:"centers"`
}

// ClaimStatsDTO represents claim statistics
type ClaimStatsDTO struct {
	Total           int64            `json:"total"`
	ByStatus        map[string]int64 `json:"by_status"`
	ByType          map[string]int64 `json:"by_type"`
	PendingReview   int64            `json:"pending_review"`
	PendingApproval int64            `json:"pending_approval"`
}

// PackageStatsDTO represents package statistics
type PackageStatsDTO struct {
	Total          int64            `json:"total"`
	ByStatus       map[string]int64 `json:"by_status"`
	PendingPayment int64            `json:"pending_payment"`
}

// AmountStatsDTO represents amount statistics
type AmountStatsDTO struct {
	TotalRequested int64 `json:"total_requested"`
	TotalApproved  int64 `json:"total_approved"`
	TotalDeduction int64 `json:"total_deduction"`
	TotalPaid      int64 `json:"total_paid"`
}

// CenterStatsDTO represents center statistics
type CenterStatsDTO struct {
	Total       int64            `json:"total"`
	Active      int64            `json:"active"`
	ByType      map[string]int64 `json:"by_type"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// SuccessResponse represents success response
type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}
