// pipeline_client.go
package infrastructure

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
)

// HTTPPipelineClient implementa domain.PipelineReportesClient llamando por
// HTTP a modelo_reportes (/infer) y clasificador_reportes (/clasificar).
// Dentro de la red de Docker (app_internal_net) estos servicios se alcanzan
// por su nombre de contenedor, ej. http://modelo_reportes:8000 -- no hace
// falta pasar por el proxy de nginx para llamadas backend-a-backend, eso es
// solo para acceso externo/testing.
type HTTPPipelineClient struct {
	modeloReportesURL string
	clasificadorURL   string
	http              *http.Client
}

func NewHTTPPipelineClient(modeloReportesURL, clasificadorURL string) *HTTPPipelineClient {
	return &HTTPPipelineClient{
		modeloReportesURL: strings.TrimRight(modeloReportesURL, "/"),
		clasificadorURL:   strings.TrimRight(clasificadorURL, "/"),
		http: &http.Client{
			// Timeout corto: si el microservicio esta caido o tarda demasiado,
			// preferimos marcar el reporte como error y seguir, no colgar la
			// goroutine indefinidamente.
			Timeout: 8 * time.Second,
		},
	}
}

// --- modelo_reportes -------------------------------------------------------

type inferRequestBody struct {
	Reporte  string `json:"reporte"`
	Tiempo   *int   `json:"tiempo,omitempty"`
	TenantID int    `json:"tenant_id"`
}

type inferResponseBody struct {
	ID               int    `json:"id"`
	NivelRiesgo      string `json:"nivel_riesgo"`
	NivelRiesgoFinal string `json:"nivel_riesgo_final"`
}

func (c *HTTPPipelineClient) InferirRiesgo(ctx context.Context, reporte string, tenantID int, tiempoMin *int) (*repositories.InferenciaResultado, error) {
	body := inferRequestBody{Reporte: reporte, Tiempo: tiempoMin, TenantID: tenantID}

	var resp inferResponseBody
	if err := c.postJSON(ctx, c.modeloReportesURL+"/infer", body, &resp); err != nil {
		return nil, fmt.Errorf("modelo_reportes /infer: %w", err)
	}

	return &repositories.InferenciaResultado{
		ID:               resp.ID,
		NivelRiesgo:      resp.NivelRiesgo,
		NivelRiesgoFinal: resp.NivelRiesgoFinal,
	}, nil
}

// --- clasificador_reportes --------------------------------------------------

type clasificarRequestBody struct {
	Reporte      string  `json:"reporte"`
	Tiempo       *int    `json:"tiempo,omitempty"`
	InferenciaID *int    `json:"inferencia_id,omitempty"`
	Origen       *string `json:"origen,omitempty"`
	TenantID     int     `json:"tenant_id"`
}

type clasificarResponseBody struct {
	ID        int     `json:"id"`
	Categoria string  `json:"categoria"`
	Subtipo   *string `json:"subtipo"`
	Confianza float64 `json:"confianza"`
	Accion    string  `json:"accion"`
}

func (c *HTTPPipelineClient) ClasificarReporte(ctx context.Context, reporte string, tenantID int, inferenciaID *int, origen *string, tiempoMin *int) (*repositories.ClasificacionResultado, error) {
	body := clasificarRequestBody{
		Reporte:      reporte,
		Tiempo:       tiempoMin,
		InferenciaID: inferenciaID,
		Origen:       origen,
		TenantID:     tenantID,
	}

	var resp clasificarResponseBody
	if err := c.postJSON(ctx, c.clasificadorURL+"/clasificar", body, &resp); err != nil {
		return nil, fmt.Errorf("clasificador_reportes /clasificar: %w", err)
	}

	return &repositories.ClasificacionResultado{
		ID:        resp.ID,
		Categoria: resp.Categoria,
		Subtipo:   resp.Subtipo,
		Confianza: resp.Confianza,
		Accion:    resp.Accion,
	}, nil
}

// --- helper compartido -------------------------------------------------------

func (c *HTTPPipelineClient) postJSON(ctx context.Context, url string, payload interface{}, out interface{}) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("no se pudo serializar el request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("no se pudo construir el request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("fallo la conexion: %w", err)
	}
	defer res.Body.Close()

	respBody, _ := io.ReadAll(res.Body)

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("respuesta %d: %s", res.StatusCode, strings.TrimSpace(string(respBody)))
	}

	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("no se pudo interpretar la respuesta: %w", err)
	}

	return nil
}
