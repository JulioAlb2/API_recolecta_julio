// notificacion_routes.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
)

type NotificacionRouter struct {
	engine *gin.Engine
}

func NewNotificacionRouter(engine *gin.Engine) *NotificacionRouter {
	return &NotificacionRouter{
		engine: engine,
	}
}

func (router *NotificacionRouter) Run() {
	// Inicializar todas las dependencias
	countActivasByUsuarioIDController,
		countByCamionIDController,
		countByTipoController,
		countByUsuarioIDController,
		getAllNotificacionesController,
		getNotificacionByIdController,
		getActivasByUsuarioIDController,
		getActivasController,
		getInactivasController,
		getByCamionIDController,
		getByCamionYTipoController,
		getByUsuarioIDController,
		getByUsuarioYTipoController,
		getByTipoController,
		getByFechaRangeController,
		getGlobalesController,
		marcarComoActivaController,
		marcarComoLeidaController,
		marcarTodasComoLeidasController,
		obtenerNoLeidasController := InitNotificacionDependencies()

	// Grupo de rutas para notificaciones con prefijo /api
	notificacionGroup := router.engine.Group("/api/notificaciones")
	{
		// ================== RUTAS DE CONSULTA BÁSICAS ==================
		notificacionGroup.GET("/", getAllNotificacionesController.Run)
		notificacionGroup.GET("/:id", getNotificacionByIdController.Run)
		
		// ================== RUTAS DE CONTEO ==================
		notificacionGroup.GET("/count/usuario/:usuario_id", countByUsuarioIDController.Run)
		notificacionGroup.GET("/count/activas/usuario/:usuario_id", countActivasByUsuarioIDController.Run)
		notificacionGroup.GET("/count/tipo/:tipo", countByTipoController.Run)
		notificacionGroup.GET("/count/camion/:camion_id", countByCamionIDController.Run)
		
		// ================== RUTAS DE ESTADO ==================
		notificacionGroup.PATCH("/:id/marcar-leida", marcarComoLeidaController.Run)
		notificacionGroup.PATCH("/:id/reactivar", marcarComoActivaController.Run)
		notificacionGroup.PATCH("/usuario/:usuario_id/marcar-todas-leidas", marcarTodasComoLeidasController.Run)
		
		// ================== RUTAS DE CONSULTA POR USUARIO ==================
		notificacionGroup.GET("/usuario/:usuario_id", getByUsuarioIDController.Run)
		notificacionGroup.GET("/activas/usuario/:usuario_id", getActivasByUsuarioIDController.Run)
		notificacionGroup.GET("/usuario/:usuario_id/tipo/:tipo", getByUsuarioYTipoController.Run)
		
		// ================== RUTAS DE CONSULTA POR CAMIÓN ==================
		notificacionGroup.GET("/camion/:camion_id", getByCamionIDController.Run)
		notificacionGroup.GET("/camion/:camion_id/tipo/:tipo", getByCamionYTipoController.Run)
		
		// ================== RUTAS DE CONSULTA POR TIPO ==================
		notificacionGroup.GET("/tipo/:tipo", getByTipoController.Run)
		
		// ================== RUTAS DE CONSULTA ESPECIALES ==================
		notificacionGroup.GET("/activas", getActivasController.Run)
		notificacionGroup.GET("/inactivas", getInactivasController.Run)
		notificacionGroup.GET("/globales", getGlobalesController.Run)
		
		// ================== RUTAS DE CONSULTA POR FECHAS ==================
		notificacionGroup.GET("/rango-fecha", getByFechaRangeController.Run)
		
		// ================== RUTAS ALTERNATIVAS (DUPLICADOS) ==================
		notificacionGroup.GET("/no-leidas/usuario/:usuario_id", obtenerNoLeidasController.Run)
	}
}