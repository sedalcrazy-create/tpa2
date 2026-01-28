package repository

import (
	"context"

	"github.com/bank-melli/tpa/internal/domain/entity"
)

// Pagination represents pagination parameters
type Pagination struct {
	Page     int
	PageSize int
	Sort     string
	Order    string // asc, desc
}

// PaginatedResult represents a paginated result
type PaginatedResult[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// Filter represents query filters
type Filter struct {
	Field    string
	Operator string // eq, ne, gt, gte, lt, lte, like, in, between
	Value    interface{}
}

// QueryOptions represents query options
type QueryOptions struct {
	Pagination *Pagination
	Filters    []Filter
	Preloads   []string
	TenantID   uint
}

// BaseRepository defines common repository methods
type BaseRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	FindByID(ctx context.Context, id uint, opts ...QueryOptions) (*T, error)
	FindAll(ctx context.Context, opts ...QueryOptions) ([]T, error)
	FindWithPagination(ctx context.Context, opts QueryOptions) (*PaginatedResult[T], error)
	Count(ctx context.Context, opts ...QueryOptions) (int64, error)
}

// ClaimRepository defines claim-specific repository methods
type ClaimRepository interface {
	BaseRepository[entity.Claim]

	// Find by tracking code
	FindByTrackingCode(ctx context.Context, tenantID uint, trackingCode string) (*entity.Claim, error)

	// Find by policy member
	FindByPolicyMember(ctx context.Context, tenantID uint, policyMemberID uint, opts QueryOptions) (*PaginatedResult[entity.Claim], error)

	// Find by package
	FindByPackage(ctx context.Context, tenantID uint, packageID uint, opts QueryOptions) (*PaginatedResult[entity.Claim], error)

	// Find by status
	FindByStatus(ctx context.Context, tenantID uint, status entity.ClaimStatus, opts QueryOptions) (*PaginatedResult[entity.Claim], error)

	// Find by center
	FindByCenter(ctx context.Context, tenantID uint, centerID uint, opts QueryOptions) (*PaginatedResult[entity.Claim], error)

	// Find pending for examination
	FindPendingExamination(ctx context.Context, tenantID uint, examinerID *uint, opts QueryOptions) (*PaginatedResult[entity.Claim], error)

	// Statistics
	GetStatsByStatus(ctx context.Context, tenantID uint) (map[entity.ClaimStatus]int64, error)
	GetStatsByType(ctx context.Context, tenantID uint) (map[entity.ClaimType]int64, error)
	GetTotalAmounts(ctx context.Context, tenantID uint, filters []Filter) (requested, approved, deduction int64, err error)
}

// ClaimItemRepository defines claim item repository methods
type ClaimItemRepository interface {
	BaseRepository[entity.ClaimItem]

	FindByClaim(ctx context.Context, claimID uint) ([]entity.ClaimItem, error)
	BatchCreate(ctx context.Context, items []entity.ClaimItem) error
	BatchUpdate(ctx context.Context, items []entity.ClaimItem) error
}

// PackageRepository defines package-specific repository methods
type PackageRepository interface {
	BaseRepository[entity.Package]

	// Find by center
	FindByCenter(ctx context.Context, tenantID uint, centerID uint, opts QueryOptions) (*PaginatedResult[entity.Package], error)

	// Find by status
	FindByStatus(ctx context.Context, tenantID uint, status entity.PackageStatus, opts QueryOptions) (*PaginatedResult[entity.Package], error)

	// Find pending payment
	FindPendingPayment(ctx context.Context, tenantID uint, opts QueryOptions) (*PaginatedResult[entity.Package], error)

	// Statistics
	GetStatsByStatus(ctx context.Context, tenantID uint) (map[entity.PackageStatus]int64, error)
	GetTotalAmounts(ctx context.Context, tenantID uint, centerID *uint) (total, approved, paid int64, err error)
}

// CenterRepository defines center-specific repository methods
type CenterRepository interface {
	BaseRepository[entity.Center]

	// Find by SIAM ID
	FindBySiamID(ctx context.Context, tenantID uint, siamID string) (*entity.Center, error)

	// Find by type
	FindByType(ctx context.Context, tenantID uint, centerType entity.CenterType, opts QueryOptions) (*PaginatedResult[entity.Center], error)

	// Find active centers
	FindActive(ctx context.Context, tenantID uint, opts QueryOptions) (*PaginatedResult[entity.Center], error)

	// Find with contract status
	FindByContractStatus(ctx context.Context, tenantID uint, status entity.ContractStatus, opts QueryOptions) (*PaginatedResult[entity.Center], error)

	// Search
	Search(ctx context.Context, tenantID uint, query string, opts QueryOptions) (*PaginatedResult[entity.Center], error)
}

// PersonRepository defines person-specific repository methods
type PersonRepository interface {
	BaseRepository[entity.Person]

	// Find by national code
	FindByNationalCode(ctx context.Context, tenantID uint, nationalCode string) (*entity.Person, error)

	// Search
	Search(ctx context.Context, tenantID uint, query string, opts QueryOptions) (*PaginatedResult[entity.Person], error)
}

// PolicyRepository defines policy-specific repository methods
type PolicyRepository interface {
	BaseRepository[entity.Policy]

	// Find by policy number
	FindByPolicyNumber(ctx context.Context, tenantID uint, policyNumber string) (*entity.Policy, error)

	// Find active policies
	FindActive(ctx context.Context, tenantID uint, opts QueryOptions) (*PaginatedResult[entity.Policy], error)
}

// PolicyMemberRepository defines policy member repository methods
type PolicyMemberRepository interface {
	BaseRepository[entity.PolicyMember]

	// Find by national code
	FindByNationalCode(ctx context.Context, tenantID uint, nationalCode string) ([]entity.PolicyMember, error)

	// Find by supervisor national code
	FindBySupervisorNationalCode(ctx context.Context, tenantID uint, supervisorNationalCode string) ([]entity.PolicyMember, error)

	// Find by policy
	FindByPolicy(ctx context.Context, tenantID uint, policyID uint, opts QueryOptions) (*PaginatedResult[entity.PolicyMember], error)

	// Check eligibility
	CheckEligibility(ctx context.Context, tenantID uint, memberID uint, serviceDate string) (bool, string, error)
}

// DrugRepository defines drug-specific repository methods
type DrugRepository interface {
	BaseRepository[entity.Drug]

	// Find by IRC code
	FindByIRCCode(ctx context.Context, ircCode string) (*entity.Drug, error)

	// Find by generic code
	FindByGenericCode(ctx context.Context, genericCode string) ([]entity.Drug, error)

	// Search
	Search(ctx context.Context, query string, opts QueryOptions) (*PaginatedResult[entity.Drug], error)

	// Get current price
	GetCurrentPrice(ctx context.Context, drugID uint) (*entity.DrugPrice, error)

	// Get interactions
	GetInteractions(ctx context.Context, drugID uint) ([]entity.DrugInteraction, error)
}

// ServiceRepository defines service-specific repository methods
type ServiceRepository interface {
	BaseRepository[entity.Service]

	// Find by code
	FindByCode(ctx context.Context, code string) (*entity.Service, error)

	// Search
	Search(ctx context.Context, query string, opts QueryOptions) (*PaginatedResult[entity.Service], error)

	// Get current price
	GetCurrentPrice(ctx context.Context, serviceID uint) (*entity.ServicePrice, error)
}

// UserRepository defines user-specific repository methods
type UserRepository interface {
	BaseRepository[entity.User]

	// Find by username
	FindByUsername(ctx context.Context, username string) (*entity.User, error)

	// Find by email
	FindByEmail(ctx context.Context, email string) (*entity.User, error)

	// Find by center
	FindByCenter(ctx context.Context, centerID uint, opts QueryOptions) (*PaginatedResult[entity.User], error)

	// Update password
	UpdatePassword(ctx context.Context, userID uint, passwordHash string) error

	// Update last login
	UpdateLastLogin(ctx context.Context, userID uint) error
}

// SettlementRepository defines settlement-specific repository methods
type SettlementRepository interface {
	BaseRepository[entity.Settlement]

	// Find by center
	FindByCenter(ctx context.Context, tenantID uint, centerID uint, opts QueryOptions) (*PaginatedResult[entity.Settlement], error)

	// Find by status
	FindByStatus(ctx context.Context, tenantID uint, status uint8, opts QueryOptions) (*PaginatedResult[entity.Settlement], error)

	// Find pending
	FindPending(ctx context.Context, tenantID uint, opts QueryOptions) (*PaginatedResult[entity.Settlement], error)
}

// DiagnosisRepository defines diagnosis repository methods
type DiagnosisRepository interface {
	BaseRepository[entity.Diagnosis]

	// Find by code (ICD)
	FindByCode(ctx context.Context, code string) (*entity.Diagnosis, error)

	// Search
	Search(ctx context.Context, query string, opts QueryOptions) (*PaginatedResult[entity.Diagnosis], error)
}

// ProviderRepository defines provider repository methods
type ProviderRepository interface {
	BaseRepository[entity.Provider]

	// Find by medical code
	FindByMedicalCode(ctx context.Context, medicalCode string) (*entity.Provider, error)

	// Find by national code
	FindByNationalCode(ctx context.Context, nationalCode string) (*entity.Provider, error)

	// Search
	Search(ctx context.Context, query string, opts QueryOptions) (*PaginatedResult[entity.Provider], error)
}
