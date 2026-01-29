package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type EmployeeIllnessHandler struct {
	db *database.Database
}

func NewEmployeeIllnessHandler(db *database.Database) *EmployeeIllnessHandler {
	return &EmployeeIllnessHandler{db: db}
}

// GetAll godoc
// @Summary List employee illnesses
// @Tags EmployeeIllnesses
// @Accept json
// @Produce json
// @Security Bearer
// @Param employee_id query int false "Filter by employee ID"
// @Param is_covered query bool false "Filter by coverage status"
// @Success 200 {array} entity.EmployeeIllness
// @Router /employee-illnesses [get]
func (h *EmployeeIllnessHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	query := h.db.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if employeeID := c.QueryParam("employee_id"); employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if isCovered := c.QueryParam("is_covered"); isCovered != "" {
		query = query.Where("is_covered = ?", isCovered)
	}

	var illnesses []entity.EmployeeIllness
	if err := query.Preload("Employee").
		Order("diagnosis_date DESC").
		Find(&illnesses).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, illnesses)
}

// GetByID godoc
// @Summary Get employee illness by ID
// @Tags EmployeeIllnesses
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Illness ID"
// @Success 200 {object} entity.EmployeeIllness
// @Router /employee-illnesses/{id} [get]
func (h *EmployeeIllnessHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var illness entity.EmployeeIllness
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Employee").
		First(&illness).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, illness)
}

// Create godoc
// @Summary Create employee illness
// @Tags EmployeeIllnesses
// @Accept json
// @Produce json
// @Security Bearer
// @Param illness body entity.EmployeeIllness true "Illness data"
// @Success 201 {object} entity.EmployeeIllness
// @Router /employee-illnesses [post]
func (h *EmployeeIllnessHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var illness entity.EmployeeIllness
	if err := c.Bind(&illness); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	illness.TenantID = tenantID

	if err := h.db.DB.Create(&illness).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, illness)
}

// Update godoc
// @Summary Update employee illness
// @Tags EmployeeIllnesses
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Illness ID"
// @Param illness body entity.EmployeeIllness true "Illness data"
// @Success 200 {object} entity.EmployeeIllness
// @Router /employee-illnesses/{id} [put]
func (h *EmployeeIllnessHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var illness entity.EmployeeIllness
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&illness).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&illness); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&illness).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, illness)
}
