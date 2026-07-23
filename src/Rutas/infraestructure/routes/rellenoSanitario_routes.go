package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RellenoSanitarioRoutes struct {
	engine *gin.Engine

	createRellenoSanitarioController  *controllers.CreateRellenoSanitarioController
	getAllRellenoSanitarioController  *controllers.GetAllRellenoSanitarioController
	getRellenoSanitarioByIdController *controllers.GetRellenoSanitarioByIDController
	updateRellenoSanitarioController  *controllers.UpdateRellenoSanitarioController
	deleteRellenoSanitarioController  *controllers.DeleteRellenoSanitarioController
	getRellenoByNombreController      *controllers.GetRellenoSanitarioByNombreController
	existsRellenoSanitarioController  *controllers.ExistsRellenoSanitarioByIdController
}

func NewRellenoSanitarioRoutes(
	engine *gin.Engine,

	createRellenoSanitarioController  *controllers.CreateRellenoSanitarioController,
	getAllRellenoSanitarioController  *controllers.GetAllRellenoSanitarioController,
	getRellenoSanitarioByIdController *controllers.GetRellenoSanitarioByIDController,
	updateRellenoSanitarioController  *controllers.UpdateRellenoSanitarioController,
	deleteRellenoSanitarioController  *controllers.DeleteRellenoSanitarioController,
	getRellenoByNombreController      *controllers.GetRellenoSanitarioByNombreController,
	existsRellenoSanitarioController  *controllers.ExistsRellenoSanitarioByIdController,
) *RellenoSanitarioRoutes {
	return &RellenoSanitarioRoutes{
		engine: engine,

		createRellenoSanitarioController:  createRellenoSanitarioController,
		getAllRellenoSanitarioController:  getAllRellenoSanitarioController,
		getRellenoSanitarioByIdController: getRellenoSanitarioByIdController,
		updateRellenoSanitarioController:  updateRellenoSanitarioController,
		deleteRellenoSanitarioController:  deleteRellenoSanitarioController,
		getRellenoByNombreController:      getRellenoByNombreController,
		existsRellenoSanitarioController:  existsRellenoSanitarioController,
	}
}

func (r *RellenoSanitarioRoutes) Run() {
	routes := r.engine.Group("/api/relleno-sanitario")
	routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.SUPERVISOR, core.COORDINADOR, core.CONDUCTOR))
	{
		routes.POST("/", r.createRellenoSanitarioController.Execute)
		routes.GET("/", r.getAllRellenoSanitarioController.Execute)
		routes.GET("/:id", r.getRellenoSanitarioByIdController.Execute)
		routes.PUT("/:id", r.updateRellenoSanitarioController.Execute)
		routes.DELETE("/:id", r.deleteRellenoSanitarioController.Execute)

		// Extra
		routes.GET("/buscar", r.getRellenoByNombreController.Execute)
		routes.GET("/exists/:id", r.existsRellenoSanitarioController.Execute)
	}
}
