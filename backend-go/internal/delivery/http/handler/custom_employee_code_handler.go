package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type CustomEmployeeCodeHandler struct {
	db *database.Database
}

func NewCustomEmployeeCodeHandler(db *database.Database) *CustomEmployeeCodeHandler {
	return &CustomEmployeeCodeHandler{db: db}
}

// GetAll godoc
// @Summary List custom employee codes
// @Tags CustomEmployeeCodes
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {array} entity.CustomEmployeeCode
// @Router /custom-employee-codes [get]
func (h *CustomEmployeeCodeHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var codes []entity.CustomEmployeeCode
	if err := h.db.DB.Where("tenant_id = ?", tenantID).
		Order("sort_order, title").
		Find(&codes).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, codes)
}

// GetByID godoc
// @Summary Get custom employee code by ID
// @Tags CustomEmployeeCodes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Code ID"
// @Success 200 {object} entity.CustomEmployeeCode
// @Router /custom-employee-codes/{id} [get]
func (h *CustomEmployeeCodeHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var code entity.CustomEmployeeCode
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&code).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, code)
}

// Create godoc
// @Summary Create custom employee code
// @Tags CustomEmployeeCodes
// @Accept json
// @Produce json
// @Security Bearer
// @Param code body entity.CustomEmployeeCode true "Code data"
// @Success 201 {object} entity.CustomEmployeeCode
// @Router /custom-employee-codes [post]
func (h *CustomEmployeeCodeHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var code entity.CustomEmployeeCode
	if err := c.Bind(&code); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	code.TenantID = tenantID

	if err := h.db.DB.Create(&code).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, code)
}

// Update godoc
// @Summary Update custom employee code
// @Tags CustomEmployeeCodes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Code ID"
// @Param code body entity.CustomEmployeeCode true "Code data"
// @Success 200 {object} entity.CustomEmployeeCode
// @Router /custom-employee-codes/{id} [put]
func (h *CustomEmployeeCodeHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var code entity.CustomEmployeeCode
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&code).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&code); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&code).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, code)
}

// Delete godoc
// @Summary Delete custom employee code
// @Tags CustomEmployeeCodes
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Code ID"
// @Success 204
// @Router /custom-employee-codes/{id} [delete]
func (h *CustomEmployeeCodeHandler) Delete(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&entity.CustomEmployeeCode{})

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
