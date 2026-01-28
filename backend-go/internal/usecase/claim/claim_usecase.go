package claim

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/bank-melli/tpa/internal/domain/entity"
	"github.com/bank-melli/tpa/internal/domain/repository"
)

var (
	ErrClaimNotFound     = errors.New("claim not found")
	ErrInvalidStatus     = errors.New("invalid claim status for this operation")
	ErrDuplicateTracking = errors.New("duplicate tracking code")
	ErrMemberNotEligible = errors.New("policy member is not eligible")
)

// UseCase defines claim business logic
type UseCase interface {
	// CRUD
	Create(ctx context.Context, claim *entity.Claim) error
	Update(ctx context.Context, claim *entity.Claim) error
	Delete(ctx context.Context, tenantID, id uint) error
	GetByID(ctx context.Context, tenantID, id uint) (*entity.Claim, error)
	GetByTrackingCode(ctx context.Context, tenantID uint, trackingCode string) (*entity.Claim, error)

	// Listing
	List(ctx context.Context, tenantID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error)
	ListByCenter(ctx context.Context, tenantID, centerID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error)
	ListByPackage(ctx context.Context, tenantID, packageID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error)
	ListPendingExamination(ctx context.Context, tenantID uint, examinerID *uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error)

	// Workflow
	Submit(ctx context.Context, tenantID, id uint) error
	StartExamination(ctx context.Context, tenantID, id, examinerID uint) error
	CompleteExamination(ctx context.Context, tenantID, id, examinerID uint, result ExaminationResult) error
	Approve(ctx context.Context, tenantID, id, approverID uint) error
	Reject(ctx context.Context, tenantID, id, userID uint, reason string) error
	Return(ctx context.Context, tenantID, id, userID uint, reason string) error
	Void(ctx context.Context, tenantID, id, userID uint, reason string) error

	// Items
	AddItem(ctx context.Context, tenantID, claimID uint, item *entity.ClaimItem) error
	UpdateItem(ctx context.Context, tenantID uint, item *entity.ClaimItem) error
	RemoveItem(ctx context.Context, tenantID, claimID, itemID uint) error
	GetItems(ctx context.Context, claimID uint) ([]entity.ClaimItem, error)

	// Statistics
	GetStatsByStatus(ctx context.Context, tenantID uint) (map[entity.ClaimStatus]int64, error)
	GetStatsByType(ctx context.Context, tenantID uint) (map[entity.ClaimType]int64, error)
	GetTotalAmounts(ctx context.Context, tenantID uint, filters []repository.Filter) (requested, approved, deduction int64, err error)
}

// ExaminationResult holds examination results
type ExaminationResult struct {
	ApprovedAmount  int64
	Deduction       int64
	DeductionReason string
	Notes           string
	Items           []ItemExaminationResult
}

// ItemExaminationResult holds item examination results
type ItemExaminationResult struct {
	ItemID          uint
	ConfirmedPrice  int64
	Deduction       int64
	ReasonCodeIDs   []uint
	Notes           string
}

type useCase struct {
	claimRepo      repository.ClaimRepository
	claimItemRepo  repository.ClaimItemRepository
	policyMemberRepo repository.PolicyMemberRepository
}

// NewUseCase creates a new claim use case
func NewUseCase(
	claimRepo repository.ClaimRepository,
	claimItemRepo repository.ClaimItemRepository,
	policyMemberRepo repository.PolicyMemberRepository,
) UseCase {
	return &useCase{
		claimRepo:      claimRepo,
		claimItemRepo:  claimItemRepo,
		policyMemberRepo: policyMemberRepo,
	}
}

func (uc *useCase) Create(ctx context.Context, claim *entity.Claim) error {
	// Check for duplicate tracking code
	existing, _ := uc.claimRepo.FindByTrackingCode(ctx, claim.TenantID, claim.TrackingCode)
	if existing != nil {
		return ErrDuplicateTracking
	}

	// Set initial status
	claim.Status = entity.ClaimStatusWaitRegister
	claim.RegisterDate = time.Now()

	return uc.claimRepo.Create(ctx, claim)
}

func (uc *useCase) Update(ctx context.Context, claim *entity.Claim) error {
	existing, err := uc.claimRepo.FindByID(ctx, claim.ID, repository.QueryOptions{TenantID: claim.TenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	// Only allow updates for certain statuses
	if existing.Status != entity.ClaimStatusWaitRegister && existing.Status != entity.ClaimStatusReturned {
		return ErrInvalidStatus
	}

	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) Delete(ctx context.Context, tenantID, id uint) error {
	existing, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	// Only allow deletion for draft status
	if existing.Status != entity.ClaimStatusWaitRegister {
		return ErrInvalidStatus
	}

	return uc.claimRepo.Delete(ctx, id)
}

func (uc *useCase) GetByID(ctx context.Context, tenantID, id uint) (*entity.Claim, error) {
	return uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{
		TenantID: tenantID,
		Preloads: []string{
			"PolicyMember",
			"PolicyMember.Person",
			"Center",
			"Items",
			"Items.Drug",
			"Items.Service",
			"Diagnoses",
			"Diagnoses.Diagnosis",
			"Attachments",
			"Notes",
		},
	})
}

func (uc *useCase) GetByTrackingCode(ctx context.Context, tenantID uint, trackingCode string) (*entity.Claim, error) {
	return uc.claimRepo.FindByTrackingCode(ctx, tenantID, trackingCode)
}

func (uc *useCase) List(ctx context.Context, tenantID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	opts.TenantID = tenantID
	return uc.claimRepo.FindWithPagination(ctx, opts)
}

func (uc *useCase) ListByCenter(ctx context.Context, tenantID, centerID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	return uc.claimRepo.FindByCenter(ctx, tenantID, centerID, opts)
}

func (uc *useCase) ListByPackage(ctx context.Context, tenantID, packageID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	return uc.claimRepo.FindByPackage(ctx, tenantID, packageID, opts)
}

func (uc *useCase) ListPendingExamination(ctx context.Context, tenantID uint, examinerID *uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	return uc.claimRepo.FindPendingExamination(ctx, tenantID, examinerID, opts)
}

func (uc *useCase) Submit(ctx context.Context, tenantID, id uint) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitRegister && claim.Status != entity.ClaimStatusReturned {
		return ErrInvalidStatus
	}

	claim.Status = entity.ClaimStatusWaitCheck
	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) StartExamination(ctx context.Context, tenantID, id, examinerID uint) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitCheck {
		return ErrInvalidStatus
	}

	claim.HandlerUserID = &examinerID
	now := time.Now()
	claim.HandlerDate = &now

	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) CompleteExamination(ctx context.Context, tenantID, id, examinerID uint, result ExaminationResult) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitCheck && claim.Status != entity.ClaimStatusWaitCheckAgain {
		return ErrInvalidStatus
	}

	// Update claim amounts
	claim.ApprovedAmount = result.ApprovedAmount
	claim.Deduction = result.Deduction
	claim.CheckingDesc = result.Notes
	claim.Status = entity.ClaimStatusWaitCheckConfirm

	// Update items
	for _, itemResult := range result.Items {
		item, err := uc.findItemByID(ctx, claim.ID, itemResult.ItemID)
		if err != nil {
			continue
		}
		item.ConfirmedPrice = itemResult.ConfirmedPrice
		item.Deduction = itemResult.Deduction
		if err := uc.claimItemRepo.Update(ctx, item); err != nil {
			return fmt.Errorf("failed to update item %d: %w", itemResult.ItemID, err)
		}
	}

	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) Approve(ctx context.Context, tenantID, id, approverID uint) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitCheckConfirm {
		return ErrInvalidStatus
	}

	claim.Status = entity.ClaimStatusWaitSendFinancial
	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) Reject(ctx context.Context, tenantID, id, userID uint, reason string) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	claim.ApprovedAmount = 0
	claim.CheckingDesc = reason
	claim.Status = entity.ClaimStatusArchived

	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) Return(ctx context.Context, tenantID, id, userID uint, reason string) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	claim.Status = entity.ClaimStatusReturned
	claim.CheckingDesc = reason

	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) Void(ctx context.Context, tenantID, id, userID uint, reason string) error {
	claim, err := uc.claimRepo.FindByID(ctx, id, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	claim.IsVoid = true
	now := time.Now()
	claim.VoidDate = &now
	claim.VoidMessage = reason
	claim.VoidUserID = &userID

	return uc.claimRepo.Update(ctx, claim)
}

func (uc *useCase) AddItem(ctx context.Context, tenantID, claimID uint, item *entity.ClaimItem) error {
	claim, err := uc.claimRepo.FindByID(ctx, claimID, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitRegister && claim.Status != entity.ClaimStatusReturned {
		return ErrInvalidStatus
	}

	item.ClaimID = claimID
	item.TenantID = tenantID
	return uc.claimItemRepo.Create(ctx, item)
}

func (uc *useCase) UpdateItem(ctx context.Context, tenantID uint, item *entity.ClaimItem) error {
	claim, err := uc.claimRepo.FindByID(ctx, item.ClaimID, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitRegister &&
	   claim.Status != entity.ClaimStatusReturned &&
	   claim.Status != entity.ClaimStatusWaitCheck {
		return ErrInvalidStatus
	}

	return uc.claimItemRepo.Update(ctx, item)
}

func (uc *useCase) RemoveItem(ctx context.Context, tenantID, claimID, itemID uint) error {
	claim, err := uc.claimRepo.FindByID(ctx, claimID, repository.QueryOptions{TenantID: tenantID})
	if err != nil {
		return ErrClaimNotFound
	}

	if claim.Status != entity.ClaimStatusWaitRegister && claim.Status != entity.ClaimStatusReturned {
		return ErrInvalidStatus
	}

	return uc.claimItemRepo.Delete(ctx, itemID)
}

func (uc *useCase) GetItems(ctx context.Context, claimID uint) ([]entity.ClaimItem, error) {
	return uc.claimItemRepo.FindByClaim(ctx, claimID)
}

func (uc *useCase) GetStatsByStatus(ctx context.Context, tenantID uint) (map[entity.ClaimStatus]int64, error) {
	return uc.claimRepo.GetStatsByStatus(ctx, tenantID)
}

func (uc *useCase) GetStatsByType(ctx context.Context, tenantID uint) (map[entity.ClaimType]int64, error) {
	return uc.claimRepo.GetStatsByType(ctx, tenantID)
}

func (uc *useCase) GetTotalAmounts(ctx context.Context, tenantID uint, filters []repository.Filter) (requested, approved, deduction int64, err error) {
	return uc.claimRepo.GetTotalAmounts(ctx, tenantID, filters)
}

func (uc *useCase) findItemByID(ctx context.Context, claimID, itemID uint) (*entity.ClaimItem, error) {
	items, err := uc.claimItemRepo.FindByClaim(ctx, claimID)
	if err != nil {
		return nil, err
	}
	for _, item := range items {
		if item.ID == itemID {
			return &item, nil
		}
	}
	return nil, errors.New("item not found")
}
