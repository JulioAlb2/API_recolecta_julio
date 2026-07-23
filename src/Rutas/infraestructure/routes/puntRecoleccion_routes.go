package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PuntoRecoleccionRoutes struct {
	engine *gin.Engine

	createController    *controllers.CreatePuntoRecoleccionController
	getAllController    *controllers.GetAllPuntoRecoleccionController
	getByIdController   *controllers.GetPuntoRecoleccionByIdController
	getByRutaController *controllers.GetPuntoRecoleccionByRutaController
	updateController    *controllers.UpdatePuntoRecoleccionController
	deleteController    *controllers.DeletePuntoRecoleccionController
}

func NewPuntoRecoleccionRoutes(
	engine *gin.Engine,
	createController *controllers.CreatePuntoRecoleccionController,
	getAllController *controllers.GetAllPuntoRecoleccionController,
	getByIdController *controllers.GetPuntoRecoleccionByIdController,
	getByRutaController *controllers.GetPuntoRecoleccionByRutaController,
	updateController *controllers.UpdatePuntoRecoleccionController,
	deleteController *controllers.DeletePuntoRecoleccionController,
) *PuntoRecoleccionRoutes {
	return &PuntoRecoleccionRoutes{
		engine:              engine,
		createController:    createController,
		getAllController:    getAllController,
		getByIdController:   getByIdController,
		getByRutaController: getByRutaController,
		updateController:    updateController,
		deleteController:    deleteController,
	}
}

func (r *PuntoRecoleccionRoutes) Run() {
	routes := r.engine.Group("/api/puntos-recoleccion")
	routes.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
	{
		routes.POST("/", r.createController.Run)
		routes.GET("/", r.getAllController.Run)
		routes.GET("/:id", r.getByIdController.Run)
		routes.GET("/ruta/:rutaId", r.getByRutaController.Run)
		routes.PUT("/:id", r.updateController.Run)
		routes.DELETE("/:id", r.deleteController.Run)
	}
}
