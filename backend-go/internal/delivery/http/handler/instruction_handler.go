package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type InstructionHandler struct {
	db *database.Database
}

func NewInstructionHandler(db *database.Database) *InstructionHandler {
	return &InstructionHandler{db: db}
}

// GetAll godoc
// @Summary List instructions
// @Tags Instructions
// @Accept json
// @Produce json
// @Security Bearer
// @Param category query string false "Filter by category"
// @Param is_template query bool false "Filter by template status"
// @Success 200 {array} entity.Instruction
// @Router /instructions [get]
func (h *InstructionHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	query := h.db.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if category := c.QueryParam("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if isTemplate := c.QueryParam("is_template"); isTemplate != "" {
		query = query.Where("is_template = ?", isTemplate)
	}

	var instructions []entity.Instruction
	if err := query.Order("sort_order, title_fa").Find(&instructions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, instructions)
}

// GetByID godoc
// @Summary Get instruction by ID
// @Tags Instructions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Instruction ID"
// @Success 200 {object} entity.Instruction
// @Router /instructions/{id} [get]
func (h *InstructionHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var instruction entity.Instruction
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&instruction).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, instruction)
}

// Create godoc
// @Summary Create instruction
// @Tags Instructions
// @Accept json
// @Produce json
// @Security Bearer
// @Param instruction body entity.Instruction true "Instruction data"
// @Success 201 {object} entity.Instruction
// @Router /instructions [post]
func (h *InstructionHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var instruction entity.Instruction
	if err := c.Bind(&instruction); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	instruction.TenantID = tenantID

	if err := h.db.DB.Create(&instruction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, instruction)
}

// Update godoc
// @Summary Update instruction
// @Tags Instructions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Instruction ID"
// @Param instruction body entity.Instruction true "Instruction data"
// @Success 200 {object} entity.Instruction
// @Router /instructions/{id} [put]
func (h *InstructionHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var instruction entity.Instruction
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&instruction).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&instruction); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&instruction).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, instruction)
}

// Delete godoc
// @Summary Delete instruction
// @Tags Instructions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Instruction ID"
// @Success 204
// @Router /instructions/{id} [delete]
func (h *InstructionHandler) Delete(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	result := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Delete(&entity.Instruction{})

	if result.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
	}

	if result.RowsAffected == 0 {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.NoContent(http.StatusNoContent)
}
