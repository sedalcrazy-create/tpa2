package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"tpa-backend/internal/domain/entity"
)

type EmployeeHandler struct {
	// TODO: Add usecase dependency
}

func NewEmployeeHandler() *EmployeeHandler {
	return &EmployeeHandler{}
}

// GetEmployees - لیست کارمندان (عین actionAdmin در Yii)
// GET /api/v1/employees
func (h *EmployeeHandler) GetEmployees(c *fiber.Ctx) error {
	// Query parameters
	search := c.Query("search", "")
	status := c.Query("status", "all")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// Mock data برای تست
	employees := []fiber.Map{
		{
			"id":             1,
			"personnel_code": "12345",
			"national_code":  "1234567890",
			"first_name":     "علی",
			"last_name":      "احمدی",
			"parent_id":      nil,
			"relation_type":  "کارمند اصلی",
			"is_active":      true,
			"status":         "active",
		},
		{
			"id":             2,
			"personnel_code": "12346",
			"national_code":  "1234567891",
			"first_name":     "فاطمه",
			"last_name":      "محمدی",
			"parent_id":      1,
			"relation_type":  "همسر",
			"is_active":      true,
			"status":         "active",
		},
		{
			"id":             3,
			"personnel_code": "12347",
			"national_code":  "1234567892",
			"first_name":     "محمد",
			"last_name":      "رضایی",
			"parent_id":      nil,
			"relation_type":  "کارمند اصلی",
			"is_active":      true,
			"status":         "active",
		},
	}

	// Statistics
	stats := fiber.Map{
		"total":          150,
		"active":         142,
		"family_members": 234,
		"retired":        8,
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"employees": employees,
			"stats":     stats,
			"pagination": fiber.Map{
				"page":  page,
				"limit": limit,
				"total": len(employees),
			},
		},
	})
}

// GetEmployee - جزئیات یک کارمند
// GET /api/v1/employees/:id
func (h *EmployeeHandler) GetEmployee(c *fiber.Ctx) error {
	id := c.Params("id")

	employee := fiber.Map{
		"id":               1,
		"personnel_code":   "12345",
		"national_code":    "1234567890",
		"first_name":       "علی",
		"last_name":        "احمدی",
		"father_name":      "حسن",
		"birth_date":       "1360-05-15",
		"gender":           "male",
		"marital_status":   "married",
		"phone":            "02112345678",
		"mobile":           "09121234567",
		"email":            "ali.ahmadi@example.com",
		"recruitment_date": "1395-01-01",
		"is_active":        true,
		"status":           "active",
		"family_members": []fiber.Map{
			{
				"id":            2,
				"first_name":    "فاطمه",
				"last_name":     "محمدی",
				"relation_type": "همسر",
				"national_code": "1234567891",
			},
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    employee,
	})
}

// CreateEmployee - ایجاد کارمند جدید
// POST /api/v1/employees
func (h *EmployeeHandler) CreateEmployee(c *fiber.Ctx) error {
	var req entity.Employee

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "داده‌های ورودی نامعتبر است",
		})
	}

	// TODO: Validate and save to database

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "کارمند با موفقیت ایجاد شد",
		"data": fiber.Map{
			"id": 1,
		},
	})
}

// UpdateEmployee - ویرایش کارمند
// PUT /api/v1/employees/:id
func (h *EmployeeHandler) UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var req entity.Employee

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "داده‌های ورودی نامعتبر است",
		})
	}

	// TODO: Update in database

	return c.JSON(fiber.Map{
		"success": true,
		"message": "کارمند با موفقیت به‌روزرسانی شد",
	})
}

// DeleteEmployee - حذف (soft delete) کارمند
// DELETE /api/v1/employees/:id
func (h *EmployeeHandler) DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")

	// TODO: Soft delete in database

	return c.JSON(fiber.Map{
		"success": true,
		"message": "کارمند با موفقیت حذف شد",
	})
}

// AutoCompleteLookup - جستجوی autocomplete کارمندان
// GET /api/v1/employees/autocomplete
func (h *EmployeeHandler) AutoCompleteLookup(c *fiber.Ctx) error {
	q := c.Query("q", "")

	// Mock results
	results := []fiber.Map{
		{
			"id":    1,
			"label": "علی احمدی - 12345",
			"value": 1,
		},
		{
			"id":    2,
			"label": "محمد رضایی - 12347",
			"value": 2,
		},
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}
