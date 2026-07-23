// ProcesarPipelineAnomaliaUseCase.go
package application

import (
	"context"
	"log"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
)

// ProcesarPipelineAnomaliaUseCase orquesta la validacion/clasificacion de un
// reporte ya guardado como Anomalia:
//
//  1. modelo_reportes (/infer) filtra fraude/anomalias de texto.
//  2. Si el reporte no se rechaza, clasificador_reportes (/clasificar) lo
//     categoriza y sugiere una accion sobre el grafo de rutas.
//  3. El resultado (o el error, si algun microservicio falla) se persiste en
//     la propia fila de anomalia via ActualizarPipeline.
//
// Pensado para correr en background (goroutine) disparado desde
// CreateAnomaliaUseCase, para no hacer esperar al ciudadano/conductor que
// esta reportando mientras corren los modelos.
type ProcesarPipelineAnomaliaUseCase struct {
	repo     repositories.IAnomalia
	pipeline repositories.PipelineReportesClient
}

func NewProcesarPipelineAnomaliaUseCase(repo repositories.IAnomalia, pipeline repositories.PipelineReportesClient) *ProcesarPipelineAnomaliaUseCase {
	return &ProcesarPipelineAnomaliaUseCase{repo: repo, pipeline: pipeline}
}

// Run no devuelve error al caller: como corre en background despues de que
// el reporte ya fue aceptado y persistido, el unico lugar donde puede dejar
// constancia de como termino es en la propia fila (estado_pipeline).
func (uc *ProcesarPipelineAnomaliaUseCase) Run(anomaliaID int32, descripcion string, tenantID int, origen *string) {
	ctx := context.Background()

	inferencia, err := uc.pipeline.InferirRiesgo(ctx, descripcion, tenantID, nil)
	if err != nil {
		log.Println("pipeline reportes: fallo modelo_reportes para anomalia", anomaliaID, ":", err)
		uc.marcarError(anomaliaID, "modelo_reportes no disponible: "+err.Error())
		return
	}

	nivelRiesgo := inferencia.NivelRiesgoFinal
	inferenciaID := int32(inferencia.ID)

	if nivelRiesgo == "alto" {
		// Reporte con alta probabilidad de ser fraudulento/spam: no tiene
		// sentido gastar el clasificador en el. Se marca y se corta aqui.
		if err := uc.repo.ActualizarPipeline(anomaliaID, "rechazado", &nivelRiesgo, &inferenciaID, nil, nil, nil, nil); err != nil {
			log.Println("pipeline reportes: no se pudo guardar el rechazo de la anomalia", anomaliaID, ":", err)
		}
		return
	}

	clasificacion, err := uc.pipeline.ClasificarReporte(ctx, descripcion, tenantID, &inferencia.ID, origen, nil)
	if err != nil {
		log.Println("pipeline reportes: fallo clasificador_reportes para anomalia", anomaliaID, ":", err)
		uc.marcarError(anomaliaID, "clasificador_reportes no disponible: "+err.Error())
		return
	}

	if err := uc.repo.ActualizarPipeline(
		anomaliaID,
		"clasificado",
		&nivelRiesgo,
		&inferenciaID,
		&clasificacion.Categoria,
		clasificacion.Subtipo,
		&clasificacion.Accion,
		nil,
	); err != nil {
		log.Println("pipeline reportes: no se pudo guardar la clasificacion de la anomalia", anomaliaID, ":", err)
	}

	// TODO(algoritmo genetico de rutas): cuando exista ese componente, aqui
	// es donde se dispara la redireccion si clasificacion.Categoria ==
	// "calle_tapada" (equivalente: clasificacion.Accion en
	// {"block_edge", "inflate_weight"}). Por ahora solo queda persistida la
	// accion_sugerida en la anomalia para que se pueda consultar/disparar
	// manualmente o desde un endpoint aparte.
}

func (uc *ProcesarPipelineAnomaliaUseCase) marcarError(anomaliaID int32, mensaje string) {
	if err := uc.repo.ActualizarPipeline(anomaliaID, "error", nil, nil, nil, nil, nil, &mensaje); err != nil {
		log.Println("pipeline reportes: no se pudo guardar el error de la anomalia", anomaliaID, ":", err)
	}
}
