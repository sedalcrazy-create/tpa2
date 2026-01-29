package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type ContractHandler struct {
	db *database.Database
}

func NewContractHandler(db *database.Database) *ContractHandler {
	return &ContractHandler{db: db}
}

// GetAll godoc
// @Summary List contracts
// @Tags Contracts
// @Accept json
// @Produce json
// @Security Bearer
// @Param status query string false "Filter by status"
// @Param employer_name query string false "Filter by employer name"
// @Success 200 {array} entity.Contract
// @Router /contracts [get]
func (h *ContractHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	query := h.db.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if status := c.QueryParam("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if employerName := c.QueryParam("employer_name"); employerName != "" {
		query = query.Where("employer_name ILIKE ?", "%"+employerName+"%")
	}

	var contracts []entity.Contract
	if err := query.Preload("ContractType").
		Order("start_date DESC").
		Find(&contracts).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, contracts)
}

// GetByID godoc
// @Summary Get contract by ID
// @Tags Contracts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Contract ID"
// @Success 200 {object} entity.Contract
// @Router /contracts/{id} [get]
func (h *ContractHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var contract entity.Contract
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("ContractType").Preload("Policies").
		First(&contract).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, contract)
}

// Create godoc
// @Summary Create contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Security Bearer
// @Param contract body entity.Contract true "Contract data"
// @Success 201 {object} entity.Contract
// @Router /contracts [post]
func (h *ContractHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	userID := c.Get("user_id").(uint)

	var contract entity.Contract
	if err := c.Bind(&contract); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	contract.TenantID = tenantID
	contract.CreatedBy = &userID

	if err := h.db.DB.Create(&contract).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, contract)
}

// Update godoc
// @Summary Update contract
// @Tags Contracts
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Contract ID"
// @Param contract body entity.Contract true "Contract data"
// @Success 200 {object} entity.Contract
// @Router /contracts/{id} [put]
func (h *ContractHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var contract entity.Contract
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&contract).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&contract); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&contract).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, contract)
}
