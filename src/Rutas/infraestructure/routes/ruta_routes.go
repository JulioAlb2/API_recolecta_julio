package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RutaRoutes struct {
	engine *gin.Engine

	createController  *controllers.CreateRutaController
	getAllController  *controllers.GetAllRutaController
	getByIdController *controllers.GetRutaByIdController
	updateController  *controllers.UpdateRutaController
	deleteController  *controllers.DeleteRutaController
	getActivas *controllers.GetRutaActivasController
}

func NewRutaRoutes(
	engine *gin.Engine,
	createController *controllers.CreateRutaController,
	getAllController *controllers.GetAllRutaController,
	getByIdController *controllers.GetRutaByIdController,
	updateController *controllers.UpdateRutaController,
	deleteController *controllers.DeleteRutaController,
	getActivasController *controllers.GetRutaActivasController,
) *RutaRoutes {
	return &RutaRoutes{
		engine: engine,

		createController:  createController,
		getAllController:  getAllController,
		getByIdController: getByIdController,
		updateController:  updateController,
		deleteController:  deleteController,
		getActivas: getActivasController,
	}
}

func (r *RutaRoutes) Run() {
	routes := r.engine.Group("/api/rutas")
	routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
	{
		routes.POST("/", r.createController.Run)
		routes.GET("/", r.getAllController.Run)
		routes.GET("/activas", r.getActivas.Run)
		routes.GET("/:id", r.getByIdController.Run)
		routes.PUT("/:id", r.updateController.Run)
		routes.DELETE("/:id", r.deleteController.Run)
	}
}
