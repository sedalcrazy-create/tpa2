package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type ItemPriceConditionHandler struct {
	db *database.Database
}

func NewItemPriceConditionHandler(db *database.Database) *ItemPriceConditionHandler {
	return &ItemPriceConditionHandler{db: db}
}

// GetAll godoc
// @Summary List item price conditions
// @Tags ItemPriceConditions
// @Accept json
// @Produce json
// @Security Bearer
// @Param item_id query int false "Filter by item ID"
// @Param insurance_id query int false "Filter by insurance ID"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {array} entity.ItemPriceCondition
// @Router /item-price-conditions [get]
func (h *ItemPriceConditionHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	query := h.db.DB.Where("tenant_id = ?", tenantID)

	// Filters
	if itemID := c.QueryParam("item_id"); itemID != "" {
		query = query.Where("item_id = ?", itemID)
	}
	if insuranceID := c.QueryParam("insurance_id"); insuranceID != "" {
		query = query.Where("insurance_id = ?", insuranceID)
	}
	if isActive := c.QueryParam("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive)
	}

	var conditions []entity.ItemPriceCondition
	if err := query.Preload("Item").Preload("Category").Preload("Insurance").
		Order("priority DESC, created_at DESC").
		Find(&conditions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, conditions)
}

// GetByID godoc
// @Summary Get item price condition by ID
// @Tags ItemPriceConditions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Condition ID"
// @Success 200 {object} entity.ItemPriceCondition
// @Router /item-price-conditions/{id} [get]
func (h *ItemPriceConditionHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var condition entity.ItemPriceCondition
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Item").Preload("Category").Preload("Insurance").
		First(&condition).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, condition)
}

// Create godoc
// @Summary Create item price condition
// @Tags ItemPriceConditions
// @Accept json
// @Produce json
// @Security Bearer
// @Param condition body entity.ItemPriceCondition true "Condition data"
// @Success 201 {object} entity.ItemPriceCondition
// @Router /item-price-conditions [post]
func (h *ItemPriceConditionHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var condition entity.ItemPriceCondition
	if err := c.Bind(&condition); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	condition.TenantID = tenantID

	if err := h.db.DB.Create(&condition).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, condition)
}

// Update godoc
// @Summary Update item price condition
// @Tags ItemPriceConditions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Condition ID"
// @Param condition body entity.ItemPriceCondition true "Condition data"
// @Success 200 {object} entity.ItemPriceCondition
// @Router /item-price-conditions/{id} [put]
func (h *ItemPriceConditionHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var condition entity.ItemPriceCondition
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&condition).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&condition); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&condition).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, condition)
}

// Delete godoc
// @Summary Delete item price condition
// @Tags ItemPriceConditions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Condition ID"
// @Success 204
// @Router /item-price-conditions/{id} [delete]
func (h *ItemPriceConditionHandler) Delete(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&entity.ItemPriceCondition{})

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.NoContent(http.StatusNoContent)
}

// Calculate godoc
// @Summary Calculate pricing for an item
// @Tags ItemPriceConditions
// @Accept json
// @Produce json
// @Security Bearer
// @Param request body map[string]interface{} true "Calculation request"
// @Success 200 {object} map[string]interface{}
// @Router /item-price-conditions/calculate [post]
func (h *ItemPriceConditionHandler) Calculate(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var req struct {
		ItemID      uint   `json:"item_id"`
		InsuranceID uint   `json:"insurance_id"`
		BasePrice   int64  `json:"base_price"`
		Age         int    `json:"age"`
		Gender      string `json:"gender"`
		IsMain      bool   `json:"is_main"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Find applicable conditions
	var conditions []entity.ItemPriceCondition
	query := h.db.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if req.InsuranceID > 0 {
		query = query.Where("insurance_id = ?", req.InsuranceID)
	}

	if err := query.Order("priority DESC").Find(&conditions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Find best matching condition
	var bestCondition *entity.ItemPriceCondition
	for i := range conditions {
		if conditions[i].AppliesTo(&req.ItemID, nil, nil, req.Age, req.Gender, req.IsMain) {
			bestCondition = &conditions[i]
			break
		}
	}

	if bestCondition == nil {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"coverage": 0,
			"franchise": 0,
			"message": "No applicable pricing condition found",
		})
	}

	coverage := bestCondition.CalculateCoverage(req.BasePrice)
	franchise := bestCondition.CalculateFranchise(req.BasePrice)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"condition_id": bestCondition.ID,
		"base_price": req.BasePrice,
		"coverage": coverage,
		"franchise": franchise,
		"patient_share": req.BasePrice - coverage + franchise,
	})
}
