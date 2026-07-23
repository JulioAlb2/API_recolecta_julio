// pipeline_reportes.go
package domain

import "context"

// InferenciaResultado es la respuesta relevante de modelo_reportes (POST /infer).
type InferenciaResultado struct {
	ID               int
	NivelRiesgo      string
	NivelRiesgoFinal string
}

// ClasificacionResultado es la respuesta relevante de clasificador_reportes
// (POST /clasificar). Accion es el campo que en el futuro alimentara al
// algoritmo genetico de rutas (valores: block_edge, inflate_weight,
// marcar_mantenimiento, none).
type ClasificacionResultado struct {
	ID        int
	Categoria string
	Subtipo   *string
	Confianza float64
	Accion    string
}

// PipelineReportesClient es el puerto (en el sentido de arquitectura
// hexagonal) que expone el pipeline de validacion/clasificacion de reportes.
// La implementacion real (HTTP contra los microservicios Python) vive en
// infrastructure/pipeline_client.go; la capa application solo conoce esta
// interfaz, igual que ya hace con IAnomalia para persistencia.
type PipelineReportesClient interface {
	// InferirRiesgo llama a modelo_reportes para filtrar fraude/anomalias.
	InferirRiesgo(ctx context.Context, reporte string, tenantID int, tiempoMin *int) (*InferenciaResultado, error)

	// ClasificarReporte llama a clasificador_reportes para obtener la
	// categoria simbolica y la accion sugerida sobre el grafo de rutas.
	ClasificarReporte(ctx context.Context, reporte string, tenantID int, inferenciaID *int, origen *string, tiempoMin *int) (*ClasificacionResultado, error)
}
