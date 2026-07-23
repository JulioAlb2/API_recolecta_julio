package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/dispositivo/infrastructure/controller"
)

func RegisterDispositivoRoutes(r *gin.Engine, ctrl *controller.DispositivoController) {
	group := r.Group("/api/dispositivos")
	{
		// Rutas del conductor autenticado
		group.POST("/solicitar", core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR), ctrl.Solicitar)
		group.GET("/mi-estado", core.JWTAuthMiddleware(), core.RequireRole(core.CONDUCTOR), ctrl.MiEstado)

		// Rutas protegidas para supervisor, admin o coordinador
		group.Use(core.JWTAuthMiddleware(), core.RequireRole(core.SUPERVISOR, core.ADMIN, core.COORDINADOR))
		group.PUT("/aprobar/:conductor_id", ctrl.Aprobar)
		group.DELETE("/desvincular/:conductor_id", ctrl.Desvincular)
		group.GET("/pendientes", ctrl.ListarPendientes)
	}
}
