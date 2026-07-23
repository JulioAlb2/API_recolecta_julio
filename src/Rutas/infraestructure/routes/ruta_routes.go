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
	getActivas        *controllers.GetRutaActivasController
	arrivalController *controllers.ProcessArrivalController // Controlador de arribos
}

func NewRutaRoutes(
	engine *gin.Engine,
	createController *controllers.CreateRutaController,
	getAllController *controllers.GetAllRutaController,
	getByIdController *controllers.GetRutaByIdController,
	updateController *controllers.UpdateRutaController,
	deleteController *controllers.DeleteRutaController,
	getActivasController *controllers.GetRutaActivasController,
	arrivalController *controllers.ProcessArrivalController, // Inyección de arribos
) *RutaRoutes {
	return &RutaRoutes{
		engine: engine,

		createController:  createController,
		getAllController:  getAllController,
		getByIdController: getByIdController,
		updateController:  updateController,
		deleteController:  deleteController,
		getActivas:        getActivasController,
		arrivalController: arrivalController,
	}
}

func (r *RutaRoutes) Run() {
	routes := r.engine.Group("/api/rutas")
	{
		// Endpoint protegido para arribo a puntos (solo Conductores con dispositivo validado)
		routes.POST("/arrival", core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR), core.DeviceValidationMiddleware(), r.arrivalController.Run)

		// Lectura: cualquier JWT válido (ciudadano role_id=0 e empleados).
		// Sin esto el mapa móvil del ciudadano no recibe geometría ni puede pintar el camión.
		routes.GET("/", core.JWTAuthMiddleware(), r.getAllController.Run)
		routes.GET("/activas", core.JWTAuthMiddleware(), r.getActivas.Run)
		routes.GET("/:id", core.JWTAuthMiddleware(), r.getByIdController.Run)

		// Escritura: solo personal operativo.
		write := routes.Group("")
		write.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.CONDUCTOR, core.SUPERVISOR, core.COORDINADOR))
		write.POST("/", r.createController.Run)
		write.PUT("/:id", r.updateController.Run)
		write.DELETE("/:id", r.deleteController.Run)
	}
}
