package routes

import (
	"github.com/labstack/echo/v4"
	"tpa/internal/delivery/http/handler"
	"tpa/internal/infrastructure/database"
)

// RegisterTPARoutes registers all TPA-related routes
func RegisterTPARoutes(e *echo.Group, db *database.Database) {
	// Initialize handlers
	cecHandler := handler.NewCustomEmployeeCodeHandler(db)
	ipcHandler := handler.NewItemPriceConditionHandler(db)
	instructionHandler := handler.NewInstructionHandler(db)
	insuranceRuleHandler := handler.NewInsuranceRuleHandler(db)
	prescriptionHandler := handler.NewPrescriptionHandler(db)
	employeeIllnessHandler := handler.NewEmployeeIllnessHandler(db)
	contractHandler := handler.NewContractHandler(db)

	// Custom Employee Codes
	cec := e.Group("/custom-employee-codes")
	cec.GET("", cecHandler.GetAll)
	cec.GET("/:id", cecHandler.GetByID)
	cec.POST("", cecHandler.Create)
	cec.PUT("/:id", cecHandler.Update)
	cec.DELETE("/:id", cecHandler.Delete)

	// Item Price Conditions
	ipc := e.Group("/item-price-conditions")
	ipc.GET("", ipcHandler.GetAll)
	ipc.GET("/:id", ipcHandler.GetByID)
	ipc.POST("", ipcHandler.Create)
	ipc.PUT("/:id", ipcHandler.Update)
	ipc.DELETE("/:id", ipcHandler.Delete)
	ipc.POST("/calculate", ipcHandler.Calculate)

	// Instructions
	inst := e.Group("/instructions")
	inst.GET("", instructionHandler.GetAll)
	inst.GET("/:id", instructionHandler.GetByID)
	inst.POST("", instructionHandler.Create)
	inst.PUT("/:id", instructionHandler.Update)
	inst.DELETE("/:id", instructionHandler.Delete)

	// Insurance Rules
	rules := e.Group("/insurance-rules")
	rules.GET("", insuranceRuleHandler.GetAll)
	rules.GET("/:id", insuranceRuleHandler.GetByID)
	rules.POST("", insuranceRuleHandler.Create)
	rules.PUT("/:id", insuranceRuleHandler.Update)
	rules.DELETE("/:id", insuranceRuleHandler.Delete)

	// Prescriptions
	rx := e.Group("/prescriptions")
	rx.GET("", prescriptionHandler.GetAll)
	rx.GET("/:id", prescriptionHandler.GetByID)
	rx.POST("", prescriptionHandler.Create)
	rx.PUT("/:id", prescriptionHandler.Update)
	rx.POST("/:id/convert-to-claim", prescriptionHandler.ConvertToClaim)

	// Employee Illnesses
	illness := e.Group("/employee-illnesses")
	illness.GET("", employeeIllnessHandler.GetAll)
	illness.GET("/:id", employeeIllnessHandler.GetByID)
	illness.POST("", employeeIllnessHandler.Create)
	illness.PUT("/:id", employeeIllnessHandler.Update)

	// Contracts
	contracts := e.Group("/contracts")
	contracts.GET("", contractHandler.GetAll)
	contracts.GET("/:id", contractHandler.GetByID)
	contracts.POST("", contractHandler.Create)
	contracts.PUT("/:id", contractHandler.Update)
}
