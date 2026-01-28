package gorm

import (
	"context"

	"github.com/bank-melli/tpa/internal/domain/entity"
	"github.com/bank-melli/tpa/internal/domain/repository"
	"gorm.io/gorm"
)

type claimRepository struct {
	db *gorm.DB
}

// NewClaimRepository creates a new claim repository
func NewClaimRepository(db *gorm.DB) repository.ClaimRepository {
	return &claimRepository{db: db}
}

func (r *claimRepository) Create(ctx context.Context, claim *entity.Claim) error {
	return r.db.WithContext(ctx).Create(claim).Error
}

func (r *claimRepository) Update(ctx context.Context, claim *entity.Claim) error {
	return r.db.WithContext(ctx).Save(claim).Error
}

func (r *claimRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entity.Claim{}, id).Error
}

func (r *claimRepository) FindByID(ctx context.Context, id uint, opts ...repository.QueryOptions) (*entity.Claim, error) {
	var claim entity.Claim
	query := r.db.WithContext(ctx)

	if len(opts) > 0 && opts[0].TenantID > 0 {
		query = query.Where("tenant_id = ?", opts[0].TenantID)
	}

	if len(opts) > 0 && len(opts[0].Preloads) > 0 {
		for _, preload := range opts[0].Preloads {
			query = query.Preload(preload)
		}
	}

	err := query.First(&claim, id).Error
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

func (r *claimRepository) FindAll(ctx context.Context, opts ...repository.QueryOptions) ([]entity.Claim, error) {
	var claims []entity.Claim
	query := r.db.WithContext(ctx)

	if len(opts) > 0 {
		query = r.applyOptions(query, opts[0])
	}

	err := query.Find(&claims).Error
	return claims, err
}

func (r *claimRepository) FindWithPagination(ctx context.Context, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	var claims []entity.Claim
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Claim{})
	query = r.applyOptions(query, opts)

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Apply pagination
	if opts.Pagination != nil {
		offset := (opts.Pagination.Page - 1) * opts.Pagination.PageSize
		query = query.Offset(offset).Limit(opts.Pagination.PageSize)

		if opts.Pagination.Sort != "" {
			order := opts.Pagination.Sort
			if opts.Pagination.Order == "desc" {
				order += " DESC"
			}
			query = query.Order(order)
		}
	}

	// Execute query
	if err := query.Find(&claims).Error; err != nil {
		return nil, err
	}

	totalPages := 1
	if opts.Pagination != nil && opts.Pagination.PageSize > 0 {
		totalPages = int((total + int64(opts.Pagination.PageSize) - 1) / int64(opts.Pagination.PageSize))
	}

	return &repository.PaginatedResult[entity.Claim]{
		Items:      claims,
		Total:      total,
		Page:       opts.Pagination.Page,
		PageSize:   opts.Pagination.PageSize,
		TotalPages: totalPages,
	}, nil
}

func (r *claimRepository) Count(ctx context.Context, opts ...repository.QueryOptions) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&entity.Claim{})

	if len(opts) > 0 {
		query = r.applyOptions(query, opts[0])
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *claimRepository) FindByTrackingCode(ctx context.Context, tenantID uint, trackingCode string) (*entity.Claim, error) {
	var claim entity.Claim
	err := r.db.WithContext(ctx).
		Where("tenant_id = ? AND tracking_code = ?", tenantID, trackingCode).
		Preload("PolicyMember").
		Preload("PolicyMember.Person").
		Preload("Center").
		Preload("Items").
		Preload("Diagnoses").
		First(&claim).Error
	if err != nil {
		return nil, err
	}
	return &claim, nil
}

func (r *claimRepository) FindByPolicyMember(ctx context.Context, tenantID uint, policyMemberID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	opts.Filters = append(opts.Filters, repository.Filter{
		Field:    "policy_member_id",
		Operator: "eq",
		Value:    policyMemberID,
	})
	opts.TenantID = tenantID
	return r.FindWithPagination(ctx, opts)
}

func (r *claimRepository) FindByPackage(ctx context.Context, tenantID uint, packageID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	opts.Filters = append(opts.Filters, repository.Filter{
		Field:    "package_id",
		Operator: "eq",
		Value:    packageID,
	})
	opts.TenantID = tenantID
	return r.FindWithPagination(ctx, opts)
}

func (r *claimRepository) FindByStatus(ctx context.Context, tenantID uint, status entity.ClaimStatus, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	opts.Filters = append(opts.Filters, repository.Filter{
		Field:    "status",
		Operator: "eq",
		Value:    status,
	})
	opts.TenantID = tenantID
	return r.FindWithPagination(ctx, opts)
}

func (r *claimRepository) FindByCenter(ctx context.Context, tenantID uint, centerID uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	opts.Filters = append(opts.Filters, repository.Filter{
		Field:    "center_id",
		Operator: "eq",
		Value:    centerID,
	})
	opts.TenantID = tenantID
	return r.FindWithPagination(ctx, opts)
}

func (r *claimRepository) FindPendingExamination(ctx context.Context, tenantID uint, examinerID *uint, opts repository.QueryOptions) (*repository.PaginatedResult[entity.Claim], error) {
	opts.Filters = append(opts.Filters, repository.Filter{
		Field:    "status",
		Operator: "eq",
		Value:    entity.ClaimStatusWaitCheck,
	})
	if examinerID != nil {
		opts.Filters = append(opts.Filters, repository.Filter{
			Field:    "handler_user_id",
			Operator: "eq",
			Value:    *examinerID,
		})
	}
	opts.TenantID = tenantID
	return r.FindWithPagination(ctx, opts)
}

func (r *claimRepository) GetStatsByStatus(ctx context.Context, tenantID uint) (map[entity.ClaimStatus]int64, error) {
	type result struct {
		Status entity.ClaimStatus
		Count  int64
	}

	var results []result
	err := r.db.WithContext(ctx).
		Model(&entity.Claim{}).
		Select("status, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[entity.ClaimStatus]int64)
	for _, r := range results {
		stats[r.Status] = r.Count
	}
	return stats, nil
}

func (r *claimRepository) GetStatsByType(ctx context.Context, tenantID uint) (map[entity.ClaimType]int64, error) {
	type result struct {
		ClaimType entity.ClaimType
		Count     int64
	}

	var results []result
	err := r.db.WithContext(ctx).
		Model(&entity.Claim{}).
		Select("claim_type, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("claim_type").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	stats := make(map[entity.ClaimType]int64)
	for _, r := range results {
		stats[r.ClaimType] = r.Count
	}
	return stats, nil
}

func (r *claimRepository) GetTotalAmounts(ctx context.Context, tenantID uint, filters []repository.Filter) (requested, approved, deduction int64, err error) {
	type result struct {
		RequestAmount  int64
		ApprovedAmount int64
		Deduction      int64
	}

	var res result
	query := r.db.WithContext(ctx).
		Model(&entity.Claim{}).
		Select("COALESCE(SUM(request_amount), 0) as request_amount, COALESCE(SUM(approved_amount), 0) as approved_amount, COALESCE(SUM(deduction), 0) as deduction").
		Where("tenant_id = ?", tenantID)

	for _, f := range filters {
		query = r.applyFilter(query, f)
	}

	err = query.Scan(&res).Error
	return res.RequestAmount, res.ApprovedAmount, res.Deduction, err
}

func (r *claimRepository) applyOptions(query *gorm.DB, opts repository.QueryOptions) *gorm.DB {
	if opts.TenantID > 0 {
		query = query.Where("tenant_id = ?", opts.TenantID)
	}

	for _, filter := range opts.Filters {
		query = r.applyFilter(query, filter)
	}

	for _, preload := range opts.Preloads {
		query = query.Preload(preload)
	}

	return query
}

func (r *claimRepository) applyFilter(query *gorm.DB, filter repository.Filter) *gorm.DB {
	switch filter.Operator {
	case "eq":
		return query.Where(filter.Field+" = ?", filter.Value)
	case "ne":
		return query.Where(filter.Field+" != ?", filter.Value)
	case "gt":
		return query.Where(filter.Field+" > ?", filter.Value)
	case "gte":
		return query.Where(filter.Field+" >= ?", filter.Value)
	case "lt":
		return query.Where(filter.Field+" < ?", filter.Value)
	case "lte":
		return query.Where(filter.Field+" <= ?", filter.Value)
	case "like":
		return query.Where(filter.Field+" LIKE ?", "%"+filter.Value.(string)+"%")
	case "in":
		return query.Where(filter.Field+" IN ?", filter.Value)
	default:
		return query
	}
}
