package infrastructure

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
)

type FailedNotificationsController struct {
	uc *application.QueryFailedNotificationsUseCase
}

func NewFailedNotificationsController(uc *application.QueryFailedNotificationsUseCase) *FailedNotificationsController {
	return &FailedNotificationsController{uc: uc}
}

// GetFailed
// @Summary      Consultar fallas de entrega globales
// @Description  Retorna todas las notificaciones que no pudieron ser entregadas en un rango de tiempo especificado. Soporta formatos RFC3339 (ej. 2026-07-19T18:00:00Z) o solo fecha YYYY-MM-DD (ej. 2026-07-19).
// @Tags         PushNotification
// @Produce      json
// @Param        start_time query string false "Fecha de inicio (ej. 2026-07-19 o 2026-07-19T00:00:00Z)"
// @Param        end_time query string false "Fecha de fin (ej. 2026-07-19 o 2026-07-19T23:59:59Z)"
// @Success      200 {array} application.FailedNotificationRecord
// @Failure      400 {object} map[string]string "error"
// @Failure      500 {object} map[string]string "error"
// @Router       /api/notificaciones-push/fallidas [get]
func (ctrl *FailedNotificationsController) GetFailed(c *gin.Context) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")

	// Valores por defecto: últimas 24 horas si no se especifican params
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	if startStr != "" {
		parsed, err := parseFlexibleTime(startStr, false)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de start_time inválido. Usar YYYY-MM-DD o RFC3339"})
			return
		}
		start = parsed
	}
	if endStr != "" {
		parsed, err := parseFlexibleTime(endStr, true)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Formato de end_time inválido. Usar YYYY-MM-DD o RFC3339"})
			return
		}
		end = parsed
	}

	// Validar que el inicio de la búsqueda no exceda los 30 días naturales hacia atrás
	limite30Dias := time.Now().Add(-30 * 24 * time.Hour)
	if start.Before(limite30Dias) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "El rango de búsqueda es demasiado alto. Solo se permite consultar notificaciones fallidas dentro de los últimos 30 días naturales hacia hoy.",
		})
		return
	}

	records, err := ctrl.uc.Execute(c.Request.Context(), start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, records)
}

func parseFlexibleTime(str string, isEnd bool) (time.Time, error) {
	// 1. Intentar formato completo RFC3339
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		return t, nil
	}
	// 2. Intentar formato simple de fecha YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", str); err == nil {
		if isEnd {
			// Ajustar al final del día (23:59:59)
			return t.Add(23*time.Hour + 59*time.Minute + 59*time.Second), nil
		}
		return t, nil
	}
	return time.Time{}, fmt.Errorf("formato de fecha no reconocido")
}
