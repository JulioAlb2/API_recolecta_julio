// anomalia_routes.go
package infrastructure

import (
	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/core"
	alertaDomain "github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

type AnomaliaRouter struct {
	engine            *gin.Engine
	alertaRepo        alertaDomain.AlertaUsuarioRepository
	modeloReportesURL string
	clasificadorURL   string
}

func NewAnomaliaRouter(engine *gin.Engine, alertaRepo alertaDomain.AlertaUsuarioRepository, modeloReportesURL, clasificadorURL string) *AnomaliaRouter {
	return &AnomaliaRouter{
		engine:            engine,
		alertaRepo:        alertaRepo,
		modeloReportesURL: modeloReportesURL,
		clasificadorURL:   clasificadorURL,
	}
}

func (router *AnomaliaRouter) Run() {
	// Inicializar dependencias
	createController, getByIdController, updateController, deleteController,
		getAllController, getByPuntoIDController, getByChoferIDController,
		getByCamionIDController, getByRutaIDController, getByReferenciaIDController,
		getByEstadoController, getByTipoAnomaliaController, getByFechaRangeController := InitAnomaliaDependencies(router.alertaRepo, router.modeloReportesURL, router.clasificadorURL)

	// Crear anomalia: un solo endpoint, abierto a cualquier usuario
	// autenticado. Los ciudadanos no tienen role_id en el esquema de roles
	// de empleados (su JWT trae role_id: 0, ver login_ciudadano.go) y
	// CONDUCTOR (4) tampoco es staff, asi que este POST va en su propio
	// grupo con solo JWTAuthMiddleware (sin RequireRole) -- en Gin el
	// middleware se aplica por grupo, por eso no puede ir junto al resto
	// del CRUD de abajo, que si debe quedar restringido a staff.
	crearGroup := router.engine.Group("/api/anomalias")
	crearGroup.Use(core.JWTAuthMiddleware())
	{
		crearGroup.POST("/", createController.Run)
	}

	// Resto del CRUD: solo staff (ADMIN/SUPERVISOR/COORDINADOR)
	anomaliaGroup := router.engine.Group("/api/anomalias")
	anomaliaGroup.Use(core.JWTAuthMiddleware(), core.RequireRole(core.ADMIN, core.SUPERVISOR, core.COORDINADOR))
	{
		anomaliaGroup.GET("/:id", getByIdController.Run)
		anomaliaGroup.PUT("/:id", updateController.Run)
		anomaliaGroup.DELETE("/:id", deleteController.Run)
		anomaliaGroup.GET("/", getAllController.Run)

		// Rutas específicas
		anomaliaGroup.GET("/punto/:puntoId", getByPuntoIDController.Run)
		anomaliaGroup.GET("/chofer/:choferId", getByChoferIDController.Run)
		anomaliaGroup.GET("/camion/:camionId", getByCamionIDController.Run)
		anomaliaGroup.GET("/ruta/:rutaId", getByRutaIDController.Run)
		anomaliaGroup.GET("/referencia/:referenciaId", getByReferenciaIDController.Run)
		anomaliaGroup.GET("/estado", getByEstadoController.Run)        // Query param: ?estado=
		anomaliaGroup.GET("/tipo", getByTipoAnomaliaController.Run)    // Query param: ?tipo_anomalia=
		anomaliaGroup.GET("/por-fecha", getByFechaRangeController.Run) // Query params: ?fecha_inicio=&fecha_fin=
	}
}
