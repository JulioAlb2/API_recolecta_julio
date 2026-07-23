package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	"github.com/vicpoo/API_recolecta/src/core"
)

type RecorridoRoutes struct {
	engine *gin.Engine
	ctrl   *controllers.RecorridoController
}

func NewRecorridoRoutes(engine *gin.Engine, ctrl *controllers.RecorridoController) *RecorridoRoutes {
	return &RecorridoRoutes{engine: engine, ctrl: ctrl}
}

func (r *RecorridoRoutes) Run() {
	// Prefijo /api para coincidir con API_URL del móvil (.../api + /recorrido/...).
	group := r.engine.Group("/api/recorrido")
	{
		conductor := group.Group("")
		conductor.Use(
			core.JWTAuthMiddleware(),
			core.RequireRole(core.CONDUCTOR),
			core.DeviceValidationMiddleware(),
		)
		conductor.POST("/iniciar", r.ctrl.Iniciar)
		conductor.PUT("/finalizar", r.ctrl.Finalizar)
		conductor.PUT("/avanzar", r.ctrl.Avanzar)
		conductor.GET("/activo", r.ctrl.GetActivo)

		group.GET("/activo/publico", core.JWTAuthMiddleware(), r.ctrl.GetActivoPublic)
	}
}
