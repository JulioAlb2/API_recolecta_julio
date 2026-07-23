// dependencies.go
package infrastructure

import (
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

func InitNotificacionDependencies() (
	// Controladores de Conteo
	*CountNotificacionesActivasByUsuarioIDController,
	*CountNotificacionesByCamionIDController,
	*CountNotificacionesByTipoController,
	*CountNotificacionesByUsuarioIDController,
	
	// Controladores CRUD Básicos (Lectura)
	*GetAllNotificacionesController,
	*GetNotificacionByIdController,
	
	// Controladores de Consultas Activas/Inactivas
	*GetNotificacionesActivasByUsuarioIDController,
	*GetNotificacionesActivasController,
	*GetNotificacionesInactivasController,
	
	// Controladores por Camión
	*GetNotificacionesByCamionIDController,
	*GetNotificacionesByCamionYTipoController,
	
	// Controladores por Usuario
	*GetNotificacionesByUsuarioIDController,
	*GetNotificacionesByUsuarioYTipoController,
	
	// Controladores por Tipo y Rango
	*GetNotificacionesByTipoController,
	*GetNotificacionesByFechaRangeController,
	
	// Controladores de Notificaciones Especiales
	*GetNotificacionesGlobalesController,
	
	// Controladores de Estado
	*MarcarNotificacionComoActivaController,
	*MarcarNotificacionComoLeidaController,
	*MarcarTodasNotificacionesComoLeidasController,
	
	// Controladores de Notificación
	*ObtenerNumeroNotificacionesNoLeidasController,
) {
	// Repositorio PostgreSQL
	repo := NewPostgresNotificacionRepository()

	// ================== CASOS DE USO DE CONTEO ==================
	countActivasByUsuarioIDUseCase := application.NewCountNotificacionesActivasByUsuarioIDUseCase(repo)
	countByCamionIDUseCase := application.NewCountNotificacionesByCamionIDUseCase(repo)
	countByTipoUseCase := application.NewCountNotificacionesByTipoUseCase(repo)
	countByUsuarioIDUseCase := application.NewCountNotificacionesByUsuarioIDUseCase(repo)
	
	// ================== CASOS DE USO CRUD BÁSICOS (LECTURA) ==================
	getAllNotificacionesUseCase := application.NewGetAllNotificacionesUseCase(repo)
	getNotificacionByIdUseCase := application.NewGetNotificacionByIdUseCase(repo)
	
	// ================== CASOS DE USO DE CONSULTAS ==================
	getActivasByUsuarioIDUseCase := application.NewGetNotificacionesActivasByUsuarioIDUseCase(repo)
	getActivasUseCase := application.NewGetNotificacionesActivasUseCase(repo)
	getInactivasUseCase := application.NewGetNotificacionesInactivasUseCase(repo)
	
	getByCamionIDUseCase := application.NewGetNotificacionesByCamionIDUseCase(repo)
	getByCamionYTipoUseCase := application.NewGetNotificacionesByCamionYTipoUseCase(repo)
	
	getByUsuarioIDUseCase := application.NewGetNotificacionesByUsuarioIDUseCase(repo)
	getByUsuarioYTipoUseCase := application.NewGetNotificacionesByUsuarioYTipoUseCase(repo)
	
	getByTipoUseCase := application.NewGetNotificacionesByTipoUseCase(repo)
	getByFechaRangeUseCase := application.NewGetNotificacionesByFechaRangeUseCase(repo)
	
	getGlobalesUseCase := application.NewGetNotificacionesGlobalesUseCase(repo)
	
	// ================== CASOS DE USO DE ESTADO ==================
	marcarComoActivaUseCase := application.NewMarcarNotificacionComoActivaUseCase(repo)
	marcarComoLeidaUseCase := application.NewMarcarNotificacionComoLeidaUseCase(repo)
	marcarTodasComoLeidasUseCase := application.NewMarcarTodasNotificacionesComoLeidasUseCase(repo)
	
	// ================== CASOS DE USO DE NOTIFICACIÓN ==================
	obtenerNoLeidasUseCase := application.NewObtenerNumeroNotificacionesNoLeidasUseCase(repo)

	// ================== CONTROLADORES DE CONTEO ==================
	countActivasByUsuarioIDController := NewCountNotificacionesActivasByUsuarioIDController(countActivasByUsuarioIDUseCase)
	countByCamionIDController := NewCountNotificacionesByCamionIDController(countByCamionIDUseCase)
	countByTipoController := NewCountNotificacionesByTipoController(countByTipoUseCase)
	countByUsuarioIDController := NewCountNotificacionesByUsuarioIDController(countByUsuarioIDUseCase)
	
	// ================== CONTROLADORES CRUD BÁSICOS (LECTURA) ==================
	getAllNotificacionesController := NewGetAllNotificacionesController(getAllNotificacionesUseCase)
	getNotificacionByIdController := NewGetNotificacionByIdController(getNotificacionByIdUseCase)
	
	// ================== CONTROLADORES DE CONSULTAS ACTIVAS/INACTIVAS ==================
	getActivasByUsuarioIDController := NewGetNotificacionesActivasByUsuarioIDController(getActivasByUsuarioIDUseCase)
	getActivasController := NewGetNotificacionesActivasController(getActivasUseCase)
	getInactivasController := NewGetNotificacionesInactivasController(getInactivasUseCase)
	
	// ================== CONTROLADORES POR CAMIÓN ==================
	getByCamionIDController := NewGetNotificacionesByCamionIDController(getByCamionIDUseCase)
	getByCamionYTipoController := NewGetNotificacionesByCamionYTipoController(getByCamionYTipoUseCase)
	
	// ================== CONTROLADORES POR USUARIO ==================
	getByUsuarioIDController := NewGetNotificacionesByUsuarioIDController(getByUsuarioIDUseCase)
	getByUsuarioYTipoController := NewGetNotificacionesByUsuarioYTipoController(getByUsuarioYTipoUseCase)
	
	// ================== CONTROLADORES POR TIPO Y RANGO ==================
	getByTipoController := NewGetNotificacionesByTipoController(getByTipoUseCase)
	getByFechaRangeController := NewGetNotificacionesByFechaRangeController(getByFechaRangeUseCase)
	
	// ================== CONTROLADORES DE NOTIFICACIONES ESPECIALES ==================
	getGlobalesController := NewGetNotificacionesGlobalesController(getGlobalesUseCase)
	
	// ================== CONTROLADORES DE ESTADO ==================
	marcarComoActivaController := NewMarcarNotificacionComoActivaController(marcarComoActivaUseCase)
	marcarComoLeidaController := NewMarcarNotificacionComoLeidaController(marcarComoLeidaUseCase)
	marcarTodasComoLeidasController := NewMarcarTodasNotificacionesComoLeidasController(marcarTodasComoLeidasUseCase)
	
	// ================== CONTROLADORES DE NOTIFICACIÓN ==================
	obtenerNoLeidasController := NewObtenerNumeroNotificacionesNoLeidasController(obtenerNoLeidasUseCase)
	
	return countActivasByUsuarioIDController,
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
	       
	       obtenerNoLeidasController
}