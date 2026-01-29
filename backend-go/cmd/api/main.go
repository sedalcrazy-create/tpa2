package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bank-melli/tpa/internal/config"
	"github.com/bank-melli/tpa/internal/delivery/http/handler"
	"github.com/bank-melli/tpa/internal/infrastructure/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	// TEMPORARILY DISABLED for testing with mock data
	// if err := db.AutoMigrate(); err != nil {
	// 	log.Fatalf("Failed to run migrations: %v", err)
	// }

	// Seed initial data
	// if err := db.Seed(); err != nil {
	// 	log.Printf("Warning: Failed to seed data: %v", err)
	// }

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:       cfg.App.Name,
		ServerHeader:  "TPA-Server",
		StrictRouting: true,
		CaseSensitive: true,
		BodyLimit:     50 * 1024 * 1024, // 50MB
		ErrorHandler:  errorHandler,
	})

	// Global middlewares
	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))
	app.Use(helmet.New())
	app.Use(compress.New())

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Tenant-ID,X-Request-ID",
		AllowCredentials: false,
		MaxAge:           86400,
	}))

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		if err := db.HealthCheck(); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "unhealthy",
				"message": "database connection failed",
			})
		}
		return c.JSON(fiber.Map{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// API routes
	api := app.Group(cfg.App.APIPrefix)

	// Public routes (no auth required)
	setupPublicRoutes(api, db, cfg)

	// Protected routes
	setupProtectedRoutes(api, db, cfg)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.App.Port)
		log.Printf("Starting server on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Printf("Error during shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func setupPublicRoutes(api fiber.Router, db *database.Database, cfg *config.Config) {
	// Auth routes
	auth := api.Group("/auth")
	{
		auth.Post("/login", func(c *fiber.Ctx) error {
			// TODO: Implement login handler
			return c.JSON(fiber.Map{"message": "login endpoint"})
		})

		auth.Post("/refresh", func(c *fiber.Ctx) error {
			// TODO: Implement refresh token handler
			return c.JSON(fiber.Map{"message": "refresh endpoint"})
		})
	}

	// Public lookups
	lookup := api.Group("/lookup")
	{
		lookup.Get("/provinces", func(c *fiber.Ctx) error {
			// TODO: Implement provinces list
			return c.JSON(fiber.Map{"message": "provinces endpoint"})
		})

		lookup.Get("/center-types", func(c *fiber.Ctx) error {
			// TODO: Implement center types list
			return c.JSON(fiber.Map{"message": "center-types endpoint"})
		})

		lookup.Get("/claim-types", func(c *fiber.Ctx) error {
			// TODO: Implement claim types list
			return c.JSON(fiber.Map{"message": "claim-types endpoint"})
		})
	}
}

func setupProtectedRoutes(api fiber.Router, db *database.Database, cfg *config.Config) {
	// TODO: Add auth middleware when implemented
	// protected := api.Use(middleware.AuthMiddleware(&cfg.JWT))

	// For now, use api directly without auth
	protected := api

	// Dashboard
	protected.Get("/dashboard", func(c *fiber.Ctx) error {
		// TODO: Implement dashboard stats
		return c.JSON(fiber.Map{"message": "dashboard endpoint"})
	})

	// Users
	users := protected.Group("/users")
	{
		users.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list users"})
		})
		users.Post("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "create user"})
		})
		users.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get user"})
		})
		users.Put("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "update user"})
		})
		users.Delete("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "delete user"})
		})
	}

	// Claims
	claims := protected.Group("/claims")
	{
		claims.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list claims"})
		})
		claims.Post("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "create claim"})
		})
		claims.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get claim"})
		})
		claims.Put("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "update claim"})
		})
		claims.Delete("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "delete claim"})
		})
		claims.Post("/:id/submit", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "submit claim"})
		})
		claims.Post("/:id/examine", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "examine claim"})
		})
		claims.Post("/:id/approve", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "approve claim"})
		})
		claims.Get("/stats", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "claim stats"})
		})
	}

	// Packages
	packages := protected.Group("/packages")
	{
		packages.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list packages"})
		})
		packages.Post("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "create package"})
		})
		packages.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get package"})
		})
		packages.Put("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "update package"})
		})
		packages.Delete("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "delete package"})
		})
		packages.Post("/:id/submit", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "submit package"})
		})
		packages.Post("/:id/approve", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "approve package"})
		})
	}

	// Centers
	centers := protected.Group("/centers")
	{
		centers.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list centers"})
		})
		centers.Post("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "create center"})
		})
		centers.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get center"})
		})
		centers.Put("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "update center"})
		})
		centers.Delete("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "delete center"})
		})
	}

	// Settlements
	settlements := protected.Group("/settlements")
	{
		settlements.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list settlements"})
		})
		settlements.Post("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "create settlement"})
		})
		settlements.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get settlement"})
		})
		settlements.Post("/:id/approve", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "approve settlement"})
		})
	}

	// Personnel / Policy Members
	members := protected.Group("/members")
	{
		members.Get("/inquiry", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "member inquiry"})
		})
		members.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get member"})
		})
		members.Get("/:id/claims", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "member claims"})
		})
	}

	// Employees - سیستم کارمندان (عین Yii)
	employeeHandler := handler.NewEmployeeHandler()
	importHandler := handler.NewEmployeeImportHandler()

	employees := protected.Group("/employees")
	{
		// CRUD Operations
		employees.Get("/", employeeHandler.GetEmployees)               // لیست کارمندان
		employees.Get("/autocomplete", employeeHandler.AutoCompleteLookup) // جستجوی autocomplete
		employees.Get("/:id", employeeHandler.GetEmployee)             // جزئیات کارمند
		employees.Post("/", employeeHandler.CreateEmployee)            // ایجاد کارمند
		employees.Put("/:id", employeeHandler.UpdateEmployee)          // ویرایش کارمند
		employees.Delete("/:id", employeeHandler.DeleteEmployee)       // حذف کارمند

		// Import & Sync Operations (عین Yii)
		employees.Post("/upload", importHandler.UploadCSV)             // آپلود CSV
		employees.Post("/process", importHandler.ProcessCSV)           // پردازش CSV
		employees.Post("/import", importHandler.ImportEmployees)       // Import به جدول اصلی
		employees.Post("/sync", importHandler.SyncFromHR)              // همگام‌سازی با HR
		employees.Post("/validate", importHandler.ValidateCSV)         // اعتبارسنجی CSV

		// Statistics & History
		employees.Get("/stats", importHandler.GetEmployeeStats)         // آمار کارمندان
		employees.Get("/import/history", importHandler.GetImportHistory) // تاریخچه import
		employees.Get("/sample-csv", importHandler.DownloadSampleCSV)    // دانلود CSV نمونه
	}

	// Drugs
	drugs := protected.Group("/drugs")
	{
		drugs.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list drugs"})
		})
		drugs.Get("/search", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "search drugs"})
		})
		drugs.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get drug"})
		})
		drugs.Get("/:id/price", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get drug price"})
		})
	}

	// Services
	services := protected.Group("/services")
	{
		services.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list services"})
		})
		services.Get("/search", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "search services"})
		})
		services.Get("/:id", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get service"})
		})
		services.Get("/:id/price", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "get service price"})
		})
	}

	// Diagnoses (ICD)
	diagnoses := protected.Group("/diagnoses")
	{
		diagnoses.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list diagnoses"})
		})
		diagnoses.Get("/search", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "search diagnoses"})
		})
	}

	// Reports
	reports := protected.Group("/reports")
	{
		reports.Get("/claims", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "claims report"})
		})
		reports.Get("/payments", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "payments report"})
		})
		reports.Get("/centers", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "centers report"})
		})
		reports.Get("/performance", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "performance report"})
		})
	}

	// Tariffs
	tariffs := protected.Group("/tariffs")
	{
		tariffs.Get("/", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "list tariffs"})
		})
		tariffs.Get("/current", func(c *fiber.Ctx) error {
			return c.JSON(fiber.Map{"message": "current tariff"})
		})
	}
}

func errorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
		message = e.Message
	}

	return c.Status(code).JSON(fiber.Map{
		"error": fiber.Map{
			"code":    code,
			"message": message,
		},
	})
}
