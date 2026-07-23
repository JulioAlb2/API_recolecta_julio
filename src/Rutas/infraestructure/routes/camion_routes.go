package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type CamionRoutes struct {
	engine *gin.Engine

	createCamionController   *controllers.CreateCamionController
	getAllCamionController   *controllers.GetAllCamionController
	getCamionByIdController  *controllers.GetCamionByIDController
	updateCamionController   *controllers.UpdateCamionController
	deleteCamionController   *controllers.DeleteCamionController
	getCamionByPlaca         *controllers.GetCamionByPlacaController
	getCamionByModelo        *controllers.GetCamionByModeloController
	telemetryController      *controllers.ProcessTelemetryController // Controlador de telemetría
}

func NewCamionRoutes(
	engine *gin.Engine,
	createCamionController *controllers.CreateCamionController,
	getAllCamionController *controllers.GetAllCamionController,
	getCamionByIdController *controllers.GetCamionByIDController,
	updateCamionController *controllers.UpdateCamionController,
	deleteCamionController *controllers.DeleteCamionController,
	getCamionByPlaca       *controllers.GetCamionByPlacaController,
	getCamionByModelo      *controllers.GetCamionByModeloController, 
	telemetryController    *controllers.ProcessTelemetryController, // Inyección de telemetría
) *CamionRoutes {
	return &CamionRoutes{
		engine: engine,

		createCamionController:  createCamionController,
		getAllCamionController:  getAllCamionController,
		getCamionByIdController: getCamionByIdController,
		updateCamionController: updateCamionController,
		deleteCamionController:  deleteCamionController,
		getCamionByPlaca: getCamionByPlaca,
		getCamionByModelo: getCamionByModelo,
		telemetryController:    telemetryController,
	}
}

func (camionRoutes *CamionRoutes) Run() {
	routes := camionRoutes.engine.Group("/api/camion")
	{
		// Endpoint protegido para telemetría (solo Conductores con dispositivo validado)
		routes.POST("/telemetry", core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR), core.DeviceValidationMiddleware(), camionRoutes.telemetryController.Run)

		routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
		routes.POST("/", camionRoutes.createCamionController.Run)
		routes.GET("/", camionRoutes.getAllCamionController.Run)
		routes.GET("/:id", camionRoutes.getCamionByIdController.Run)
		routes.DELETE("/:id", camionRoutes.deleteCamionController.Run)
		routes.PUT("/:id", camionRoutes.updateCamionController.Run)
		routes.GET("/placa/:placa", camionRoutes.getCamionByPlaca.Run)
		routes.GET("/modelo", camionRoutes.getCamionByModelo.Run)
	}
}