package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/bank-melli/tpa/internal/delivery/http/handler"
)

func SetupEmployeeRoutes(api fiber.Router) {
	employeeHandler := handler.NewEmployeeHandler()
	importHandler := handler.NewEmployeeImportHandler()

	employees := api.Group("/employees")
	{
		// CRUD Operations
		employees.Get("/", employeeHandler.GetEmployees)           // لیست کارمندان
		employees.Get("/:id", employeeHandler.GetEmployee)         // جزئیات کارمند
		employees.Post("/", employeeHandler.CreateEmployee)        // ایجاد کارمند
		employees.Put("/:id", employeeHandler.UpdateEmployee)      // ویرایش کارمند
		employees.Delete("/:id", employeeHandler.DeleteEmployee)   // حذف کارمند
		employees.Get("/autocomplete", employeeHandler.AutoCompleteLookup) // جستجوی autocomplete

		// Import & Sync Operations
		employees.Post("/upload", importHandler.UploadCSV)           // آپلود CSV
		employees.Post("/process", importHandler.ProcessCSV)         // پردازش CSV
		employees.Post("/import", importHandler.ImportEmployees)     // Import به جدول اصلی
		employees.Post("/sync", importHandler.SyncFromHR)            // همگام‌سازی با HR
		employees.Post("/validate", importHandler.ValidateCSV)       // اعتبارسنجی CSV

		// Statistics & History
		employees.Get("/stats", importHandler.GetEmployeeStats)         // آمار کارمندان
		employees.Get("/import/history", importHandler.GetImportHistory) // تاریخچه import
		employees.Get("/sample-csv", importHandler.DownloadSampleCSV)    // دانلود CSV نمونه
	}
}
