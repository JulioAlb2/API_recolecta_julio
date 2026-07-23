package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RegistroVaciadoRoutes struct {
	engine *gin.Engine

	createController            *controllers.CreateRegistroVaciadoController
	getAllController            *controllers.GetAllRegistroVaciadoController
	getByIDController           *controllers.GetRegistroVaciadoByIDController
	getByRellenoIDController    *controllers.GetRegistroVaciadoByRellenoIDController
	getByRutaCamionIDController *controllers.GetRegistroVaciadoByRutaCamionIDController
	existsController            *controllers.ExistsRegistroVaciadoController
	deleteController            *controllers.DeleteRegistroVaciadoController
}

func NewRegistroVaciadoRoutes(
	engine *gin.Engine,
	createController *controllers.CreateRegistroVaciadoController,
	getAllController *controllers.GetAllRegistroVaciadoController,
	getByIDController *controllers.GetRegistroVaciadoByIDController,
	getByRellenoIDController *controllers.GetRegistroVaciadoByRellenoIDController,
	getByRutaCamionIDController *controllers.GetRegistroVaciadoByRutaCamionIDController,
	existsController *controllers.ExistsRegistroVaciadoController,
	deleteController *controllers.DeleteRegistroVaciadoController,
) *RegistroVaciadoRoutes {
	return &RegistroVaciadoRoutes{
		engine: engine,

		createController:            createController,
		getAllController:            getAllController,
		getByIDController:           getByIDController,
		getByRellenoIDController:    getByRellenoIDController,
		getByRutaCamionIDController: getByRutaCamionIDController,
		existsController:            existsController,
		deleteController:            deleteController,
	}
}

func (r *RegistroVaciadoRoutes) Run() {
	group := r.engine.Group("/api/registro-vaciado")
	group.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
	{
		group.POST("/", r.createController.Run)
		group.GET("/", r.getAllController.Run)
		group.GET("/:id", r.getByIDController.Run)
		group.GET("/relleno/:relleno_id", r.getByRellenoIDController.Run)
		group.GET("/ruta-camion/:ruta_camion_id", r.getByRutaCamionIDController.Run)
		group.GET("/exists/:id", r.existsController.Run)
		group.DELETE("/:id", r.deleteController.Run)
	}
}
