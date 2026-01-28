package handler

import (
	"strconv"

	"github.com/bank-melli/tpa/internal/domain/entity"
	"github.com/bank-melli/tpa/internal/domain/repository"
	"github.com/bank-melli/tpa/internal/usecase/claim"
	"github.com/gofiber/fiber/v2"
)

// ClaimHandler handles claim HTTP requests
type ClaimHandler struct {
	useCase claim.UseCase
}

// NewClaimHandler creates a new claim handler
func NewClaimHandler(uc claim.UseCase) *ClaimHandler {
	return &ClaimHandler{useCase: uc}
}

// RegisterRoutes registers claim routes
func (h *ClaimHandler) RegisterRoutes(router fiber.Router) {
	claims := router.Group("/claims")
	{
		claims.Get("/", h.List)
		claims.Post("/", h.Create)
		claims.Get("/:id", h.GetByID)
		claims.Put("/:id", h.Update)
		claims.Delete("/:id", h.Delete)

		// Workflow
		claims.Post("/:id/submit", h.Submit)
		claims.Post("/:id/start-examination", h.StartExamination)
		claims.Post("/:id/complete-examination", h.CompleteExamination)
		claims.Post("/:id/approve", h.Approve)
		claims.Post("/:id/reject", h.Reject)
		claims.Post("/:id/return", h.Return)
		claims.Post("/:id/void", h.Void)

		// Items
		claims.Get("/:id/items", h.GetItems)
		claims.Post("/:id/items", h.AddItem)
		claims.Put("/:id/items/:itemId", h.UpdateItem)
		claims.Delete("/:id/items/:itemId", h.RemoveItem)

		// Stats
		claims.Get("/stats/by-status", h.StatsByStatus)
		claims.Get("/stats/by-type", h.StatsByType)
		claims.Get("/stats/amounts", h.TotalAmounts)

		// Search
		claims.Get("/by-tracking/:trackingCode", h.GetByTrackingCode)
	}
}

// Create handles POST /claims
func (h *ClaimHandler) Create(c *fiber.Ctx) error {
	tenantID := getTenantID(c)

	var req CreateClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	claimEntity := &entity.Claim{
		TenantModel: entity.TenantModel{
			AuditModel: entity.AuditModel{
				CreatedBy: getUserID(c),
			},
			TenantID: tenantID,
		},
		TrackingCode:     req.TrackingCode,
		HID:              req.HID,
		MRN:              req.MRN,
		PolicyMemberID:   req.PolicyMemberID,
		CenterID:         req.CenterID,
		ClaimType:        req.ClaimType,
		AdmissionType:    req.AdmissionType,
		AdmissionDate:    req.AdmissionDate,
		DischargeDate:    req.DischargeDate,
		ServiceDate:      req.ServiceDate,
		RequestAmount:    req.RequestAmount,
		BasicInsShare:    req.BasicInsShare,
		RegisterUserID:   getUserID(c),
	}

	if err := h.useCase.Create(c.Context(), claimEntity); err != nil {
		if err == claim.ErrDuplicateTracking {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Status(fiber.StatusCreated).JSON(claimEntity)
}

// List handles GET /claims
func (h *ClaimHandler) List(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	opts := parseQueryOptions(c)

	result, err := h.useCase.List(c.Context(), tenantID, opts)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(result)
}

// GetByID handles GET /claims/:id
func (h *ClaimHandler) GetByID(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	claimEntity, err := h.useCase.GetByID(c.Context(), tenantID, uint(id))
	if err != nil {
		if err == claim.ErrClaimNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(claimEntity)
}

// GetByTrackingCode handles GET /claims/by-tracking/:trackingCode
func (h *ClaimHandler) GetByTrackingCode(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	trackingCode := c.Params("trackingCode")

	claimEntity, err := h.useCase.GetByTrackingCode(c.Context(), tenantID, trackingCode)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "claim not found")
	}

	return c.JSON(claimEntity)
}

// Update handles PUT /claims/:id
func (h *ClaimHandler) Update(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	var req UpdateClaimRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	existing, err := h.useCase.GetByID(c.Context(), tenantID, uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, "claim not found")
	}

	// Update fields
	if req.RequestAmount != nil {
		existing.RequestAmount = *req.RequestAmount
	}
	if req.BasicInsShare != nil {
		existing.BasicInsShare = *req.BasicInsShare
	}
	existing.UpdatedBy = getUserID(c)

	if err := h.useCase.Update(c.Context(), existing); err != nil {
		if err == claim.ErrInvalidStatus {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(existing)
}

// Delete handles DELETE /claims/:id
func (h *ClaimHandler) Delete(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid id")
	}

	if err := h.useCase.Delete(c.Context(), tenantID, uint(id)); err != nil {
		if err == claim.ErrClaimNotFound {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		if err == claim.ErrInvalidStatus {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Submit handles POST /claims/:id/submit
func (h *ClaimHandler) Submit(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	if err := h.useCase.Submit(c.Context(), tenantID, uint(id)); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "claim submitted successfully"})
}

// StartExamination handles POST /claims/:id/start-examination
func (h *ClaimHandler) StartExamination(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	examinerID := getUserID(c)

	if err := h.useCase.StartExamination(c.Context(), tenantID, uint(id), examinerID); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "examination started"})
}

// CompleteExamination handles POST /claims/:id/complete-examination
func (h *ClaimHandler) CompleteExamination(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	examinerID := getUserID(c)

	var req CompleteExaminationRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	result := claim.ExaminationResult{
		ApprovedAmount:  req.ApprovedAmount,
		Deduction:       req.Deduction,
		DeductionReason: req.DeductionReason,
		Notes:           req.Notes,
	}

	for _, item := range req.Items {
		result.Items = append(result.Items, claim.ItemExaminationResult{
			ItemID:         item.ItemID,
			ConfirmedPrice: item.ConfirmedPrice,
			Deduction:      item.Deduction,
			ReasonCodeIDs:  item.ReasonCodeIDs,
			Notes:          item.Notes,
		})
	}

	if err := h.useCase.CompleteExamination(c.Context(), tenantID, uint(id), examinerID, result); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "examination completed"})
}

// Approve handles POST /claims/:id/approve
func (h *ClaimHandler) Approve(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	approverID := getUserID(c)

	if err := h.useCase.Approve(c.Context(), tenantID, uint(id), approverID); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "claim approved"})
}

// Reject handles POST /claims/:id/reject
func (h *ClaimHandler) Reject(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	userID := getUserID(c)

	var req RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.useCase.Reject(c.Context(), tenantID, uint(id), userID, req.Reason); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "claim rejected"})
}

// Return handles POST /claims/:id/return
func (h *ClaimHandler) Return(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	userID := getUserID(c)

	var req RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.useCase.Return(c.Context(), tenantID, uint(id), userID, req.Reason); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "claim returned"})
}

// Void handles POST /claims/:id/void
func (h *ClaimHandler) Void(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	userID := getUserID(c)

	var req RejectRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.useCase.Void(c.Context(), tenantID, uint(id), userID, req.Reason); err != nil {
		return handleError(err)
	}

	return c.JSON(fiber.Map{"message": "claim voided"})
}

// GetItems handles GET /claims/:id/items
func (h *ClaimHandler) GetItems(c *fiber.Ctx) error {
	id, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	items, err := h.useCase.GetItems(c.Context(), uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(items)
}

// AddItem handles POST /claims/:id/items
func (h *ClaimHandler) AddItem(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	claimID, _ := strconv.ParseUint(c.Params("id"), 10, 32)

	var item entity.ClaimItem
	if err := c.BodyParser(&item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.useCase.AddItem(c.Context(), tenantID, uint(claimID), &item); err != nil {
		return handleError(err)
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// UpdateItem handles PUT /claims/:id/items/:itemId
func (h *ClaimHandler) UpdateItem(c *fiber.Ctx) error {
	tenantID := getTenantID(c)

	var item entity.ClaimItem
	if err := c.BodyParser(&item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	if err := h.useCase.UpdateItem(c.Context(), tenantID, &item); err != nil {
		return handleError(err)
	}

	return c.JSON(item)
}

// RemoveItem handles DELETE /claims/:id/items/:itemId
func (h *ClaimHandler) RemoveItem(c *fiber.Ctx) error {
	tenantID := getTenantID(c)
	claimID, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	itemID, _ := strconv.ParseUint(c.Params("itemId"), 10, 32)

	if err := h.useCase.RemoveItem(c.Context(), tenantID, uint(claimID), uint(itemID)); err != nil {
		return handleError(err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// StatsByStatus handles GET /claims/stats/by-status
func (h *ClaimHandler) StatsByStatus(c *fiber.Ctx) error {
	tenantID := getTenantID(c)

	stats, err := h.useCase.GetStatsByStatus(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(stats)
}

// StatsByType handles GET /claims/stats/by-type
func (h *ClaimHandler) StatsByType(c *fiber.Ctx) error {
	tenantID := getTenantID(c)

	stats, err := h.useCase.GetStatsByType(c.Context(), tenantID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(stats)
}

// TotalAmounts handles GET /claims/stats/amounts
func (h *ClaimHandler) TotalAmounts(c *fiber.Ctx) error {
	tenantID := getTenantID(c)

	requested, approved, deduction, err := h.useCase.GetTotalAmounts(c.Context(), tenantID, nil)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"requested": requested,
		"approved":  approved,
		"deduction": deduction,
	})
}

// Helper functions

func getTenantID(c *fiber.Ctx) uint {
	tenantID := c.Locals("tenant_id")
	if tenantID == nil {
		return 1 // Default tenant
	}
	return tenantID.(uint)
}

func getUserID(c *fiber.Ctx) uint {
	userID := c.Locals("user_id")
	if userID == nil {
		return 0
	}
	return userID.(uint)
}

func parseQueryOptions(c *fiber.Ctx) repository.QueryOptions {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	sort := c.Query("sort", "created_at")
	order := c.Query("order", "desc")

	opts := repository.QueryOptions{
		Pagination: &repository.Pagination{
			Page:     page,
			PageSize: pageSize,
			Sort:     sort,
			Order:    order,
		},
		Preloads: []string{"PolicyMember", "PolicyMember.Person", "Center"},
	}

	// Parse filters from query params
	if status := c.Query("status"); status != "" {
		statusInt, _ := strconv.Atoi(status)
		opts.Filters = append(opts.Filters, repository.Filter{
			Field:    "status",
			Operator: "eq",
			Value:    statusInt,
		})
	}

	if claimType := c.Query("claim_type"); claimType != "" {
		typeInt, _ := strconv.Atoi(claimType)
		opts.Filters = append(opts.Filters, repository.Filter{
			Field:    "claim_type",
			Operator: "eq",
			Value:    typeInt,
		})
	}

	if centerID := c.Query("center_id"); centerID != "" {
		id, _ := strconv.ParseUint(centerID, 10, 32)
		opts.Filters = append(opts.Filters, repository.Filter{
			Field:    "center_id",
			Operator: "eq",
			Value:    uint(id),
		})
	}

	return opts
}

func handleError(err error) error {
	switch err {
	case claim.ErrClaimNotFound:
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case claim.ErrInvalidStatus:
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	case claim.ErrMemberNotEligible:
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
}
