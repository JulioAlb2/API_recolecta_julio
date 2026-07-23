package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type EstadoCamionRoutes struct {
	engine *gin.Engine

	createEstadoCamionController *controllers.CreateEstadoCamionController
	getAllEstadoCamionController *controllers.GetAllEstadoCamionController
	getEstadoCamionByIdContrller *controllers.GetEstadoCamionByIdController
	deleteEstadoConntroller      *controllers.DeleteEstadoCamionController
	updateCamionController       *controllers.UpdateEstadoCamionController
}

func NewEstadoCamionRoutes(
	engine *gin.Engine,
	createEstadoCamionController *controllers.CreateEstadoCamionController,
	getAllEstadoCamionController *controllers.GetAllEstadoCamionController,
	getAllEstadocCamionByIdController *controllers.GetEstadoCamionByIdController,
	deleteEstadoController *controllers.DeleteEstadoCamionController,
	updateEstadoCamionController *controllers.UpdateEstadoCamionController,
) *EstadoCamionRoutes {
	return &EstadoCamionRoutes{
		engine:                       engine,
		createEstadoCamionController: createEstadoCamionController,
		getAllEstadoCamionController: getAllEstadoCamionController,
		getEstadoCamionByIdContrller: getAllEstadocCamionByIdController,
		deleteEstadoConntroller:      deleteEstadoController,
		updateCamionController:       updateEstadoCamionController,
	}
}

func (r *EstadoCamionRoutes) Run() {
	routes := r.engine.Group("/api/estado-camion")
	routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
	{
		routes.POST("/", r.createEstadoCamionController.Run)
		routes.GET("/", r.getAllEstadoCamionController.Run)
		routes.GET("/camion/:id", r.getEstadoCamionByIdContrller.Run)
		routes.PUT("/:id", r.updateCamionController.Run)
		routes.DELETE("/:id", r.deleteEstadoConntroller.Run)
	}
}
