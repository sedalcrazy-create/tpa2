package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"tpa/internal/domain/entity"
	"tpa/internal/infrastructure/database"
)

type PrescriptionHandler struct {
	db *database.Database
}

func NewPrescriptionHandler(db *database.Database) *PrescriptionHandler {
	return &PrescriptionHandler{db: db}
}

// GetAll godoc
// @Summary List prescriptions
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param employee_id query int false "Filter by employee ID"
// @Param status query string false "Filter by status"
// @Param prescription_type query string false "Filter by type"
// @Success 200 {array} entity.Prescription
// @Router /prescriptions [get]
func (h *PrescriptionHandler) GetAll(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	query := h.db.DB.Where("tenant_id = ? AND is_active = ?", tenantID, true)

	if employeeID := c.QueryParam("employee_id"); employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}
	if status := c.QueryParam("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if prescriptionType := c.QueryParam("prescription_type"); prescriptionType != "" {
		query = query.Where("prescription_type = ?", prescriptionType)
	}

	var prescriptions []entity.Prescription
	if err := query.Preload("Employee").Preload("Items").
		Order("prescription_date DESC").
		Find(&prescriptions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, prescriptions)
}

// GetByID godoc
// @Summary Get prescription by ID
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Prescription ID"
// @Success 200 {object} entity.Prescription
// @Router /prescriptions/{id} [get]
func (h *PrescriptionHandler) GetByID(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var prescription entity.Prescription
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Employee").Preload("Dependent").Preload("Center").
		Preload("Items").Preload("Items.Item").Preload("Items.Instruction").
		First(&prescription).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	return c.JSON(http.StatusOK, prescription)
}

// Create godoc
// @Summary Create prescription
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param prescription body entity.Prescription true "Prescription data"
// @Success 201 {object} entity.Prescription
// @Router /prescriptions [post]
func (h *PrescriptionHandler) Create(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)

	var prescription entity.Prescription
	if err := c.Bind(&prescription); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	prescription.TenantID = tenantID

	if err := h.db.DB.Create(&prescription).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, prescription)
}

// Update godoc
// @Summary Update prescription
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Prescription ID"
// @Param prescription body entity.Prescription true "Prescription data"
// @Success 200 {object} entity.Prescription
// @Router /prescriptions/{id} [put]
func (h *PrescriptionHandler) Update(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var prescription entity.Prescription
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		First(&prescription).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if err := c.Bind(&prescription); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.db.DB.Save(&prescription).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, prescription)
}

// ConvertToClaim godoc
// @Summary Convert prescription to claim
// @Tags Prescriptions
// @Accept json
// @Produce json
// @Security Bearer
// @Param id path int true "Prescription ID"
// @Success 200 {object} entity.Claim
// @Router /prescriptions/{id}/convert-to-claim [post]
func (h *PrescriptionHandler) ConvertToClaim(c echo.Context) error {
	tenantID := c.Get("tenant_id").(uint)
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var prescription entity.Prescription
	if err := h.db.DB.Where("tenant_id = ? AND id = ?", tenantID, id).
		Preload("Items").
		First(&prescription).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "not found"})
	}

	if !prescription.CanConvertToClaim() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "prescription cannot be converted"})
	}

	// Create claim from prescription
	// This is a simplified version - you would implement full conversion logic
	claim := entity.Claim{
		// Map prescription fields to claim fields
		// ...
	}

	if err := h.db.DB.Create(&claim).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Update prescription
	prescription.IsConvertedToClaim = true
	claimID := claim.ID
	prescription.ClaimID = &claimID
	h.db.DB.Save(&prescription)

	return c.JSON(http.StatusOK, claim)
}
