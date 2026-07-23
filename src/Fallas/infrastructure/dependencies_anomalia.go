// dependencies.go
package infrastructure

import (
	"github.com/vicpoo/API_recolecta/src/Fallas/application"
	alertaDomain "github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

// InitAnomaliaDependencies arma el dominio Fallas/Anomalia. modeloReportesURL
// y clasificadorURL vienen de config.Config (env MODELO_REPORTES_URL /
// CLASIFICADOR_URL) y alimentan al cliente HTTP del pipeline de
// validacion/clasificacion de reportes.
func InitAnomaliaDependencies(alertaRepo alertaDomain.AlertaUsuarioRepository, modeloReportesURL, clasificadorURL string) (
	*CreateAnomaliaController,
	*GetAnomaliaByIdController,
	*UpdateAnomaliaController,
	*DeleteAnomaliaController,
	*GetAllAnomaliasController,
	*GetAnomaliasByPuntoIDController,
	*GetAnomaliasByChoferIDController,
	*GetAnomaliasByCamionIDController,
	*GetAnomaliasByRutaIDController,
	*GetAnomaliasByReferenciaIDController,
	*GetAnomaliasByEstadoController,
	*GetAnomaliasByTipoAnomaliaController,
	*GetAnomaliasByFechaRangeController,
) {
	// Repositorio PostgreSQL
	repo := NewPostgresAnomaliaRepository()

	// Cliente del pipeline modelo_reportes -> clasificador_reportes
	pipelineClient := NewHTTPPipelineClient(modeloReportesURL, clasificadorURL)
	pipelineUseCase := application.NewProcesarPipelineAnomaliaUseCase(repo, pipelineClient)

	// Casos de uso
	createUseCase := application.NewCreateAnomaliaUseCase(repo, alertaRepo, pipelineUseCase)
	getByIDUseCase := application.NewGetAnomaliaByIdUseCase(repo)
	updateUseCase := application.NewUpdateAnomaliaUseCase(repo)
	deleteUseCase := application.NewDeleteAnomaliaUseCase(repo)
	getAllUseCase := application.NewGetAllAnomaliasUseCase(repo)
	getByPuntoIDUseCase := application.NewGetAnomaliasByPuntoIDUseCase(repo)
	getByChoferIDUseCase := application.NewGetAnomaliasByConductorIDUseCase(repo)
	getByCamionIDUseCase := application.NewGetAnomaliasByCamionIDUseCase(repo)
	getByRutaIDUseCase := application.NewGetAnomaliasByRutaIDUseCase(repo)
	getByReferenciaIDUseCase := application.NewGetAnomaliasByReferenciaIDUseCase(repo)
	getByEstadoUseCase := application.NewGetAnomaliasByEstadoUseCase(repo)
	getByTipoAnomaliaUseCase := application.NewGetAnomaliasByTipoAnomaliaUseCase(repo)
	getByFechaRangeUseCase := application.NewGetAnomaliasByFechaRangeUseCase(repo)

	// Controladores
	createController := NewCreateAnomaliaController(createUseCase)
	getByIDController := NewGetAnomaliaByIdController(getByIDUseCase)
	updateController := NewUpdateAnomaliaController(updateUseCase)
	deleteController := NewDeleteAnomaliaController(deleteUseCase)
	getAllController := NewGetAllAnomaliasController(getAllUseCase)
	getByPuntoIDController := NewGetAnomaliasByPuntoIDController(getByPuntoIDUseCase)
	getByChoferIDController := NewGetAnomaliasByChoferIDController(getByChoferIDUseCase)
	getByCamionIDController := NewGetAnomaliasByCamionIDController(getByCamionIDUseCase)
	getByRutaIDController := NewGetAnomaliasByRutaIDController(getByRutaIDUseCase)
	getByReferenciaIDController := NewGetAnomaliasByReferenciaIDController(getByReferenciaIDUseCase)
	getByEstadoController := NewGetAnomaliasByEstadoController(getByEstadoUseCase)
	getByTipoAnomaliaController := NewGetAnomaliasByTipoAnomaliaController(getByTipoAnomaliaUseCase)
	getByFechaRangeController := NewGetAnomaliasByFechaRangeController(getByFechaRangeUseCase)

	return createController, getByIDController, updateController, deleteController, getAllController,
		getByPuntoIDController, getByChoferIDController, getByCamionIDController, getByRutaIDController,
		getByReferenciaIDController, getByEstadoController, getByTipoAnomaliaController, getByFechaRangeController
}