package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type InsuranceRuleHandler struct {
	db *database.Database
}

func NewInsuranceRuleHandler(db *database.Database) *InsuranceRuleHandler {
	return &InsuranceRuleHandler{db: db}
}

// GetAll godoc
// @Summary List insurance rules
// @Tags InsuranceRules
// @Accept json
// @Produce json
// @Security Bearer
// @Param insurance_id query int false "Filter by insurance ID"
// @Param rule_type query string false "Filter by rule type"
// @Success 200 {array} entity.InsuranceRule
// @Router /insurance-rules [get]
func (h *InsuranceRuleHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	query := h.db.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if insuranceID := c.QueryParam("insurance_id"); insuranceID != "" {
		query = query.Where("insurance_id = ?", insuranceID)
	}
	if ruleType := c.QueryParam("rule_type"); ruleType != "" {
		query = query.Where("rule_type = ?", ruleType)
	}

	var rules []entity.InsuranceRule
	if err := query.Preload("Insurance").
		Order("priority DESC, created_at DESC").
		Find(&rules).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rules)
}

// GetByID godoc
// @Summary Get insurance rule by ID
// @Tags InsuranceRules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Rule ID"
// @Success 200 {object} entity.InsuranceRule
// @Router /insurance-rules/{id} [get]
func (h *InsuranceRuleHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var rule entity.InsuranceRule
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Insurance").
		First(&rule).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, rule)
}

// Create godoc
// @Summary Create insurance rule
// @Tags InsuranceRules
// @Accept json
// @Produce json
// @Security Bearer
// @Param rule body entity.InsuranceRule true "Rule data"
// @Success 201 {object} entity.InsuranceRule
// @Router /insurance-rules [post]
func (h *InsuranceRuleHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var rule entity.InsuranceRule
	if err := c.Bind(&rule); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	rule.TenantID = tenantID

	if err := h.db.DB.Create(&rule).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, rule)
}

// Update godoc
// @Summary Update insurance rule
// @Tags InsuranceRules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Rule ID"
// @Param rule body entity.InsuranceRule true "Rule data"
// @Success 200 {object} entity.InsuranceRule
// @Router /insurance-rules/{id} [put]
func (h *InsuranceRuleHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var rule entity.InsuranceRule
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&rule).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&rule); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&rule).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, rule)
}

// Delete godoc
// @Summary Delete insurance rule
// @Tags InsuranceRules
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Rule ID"
// @Success 204
// @Router /insurance-rules/{id} [delete]
func (h *InsuranceRuleHandler) Delete(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&entity.InsuranceRule{})

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
