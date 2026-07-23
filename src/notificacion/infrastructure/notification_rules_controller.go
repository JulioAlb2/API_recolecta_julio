package infrastructure

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/notificacion/application"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type NotificationRulesController struct {
	uc *application.ManageNotificationRulesUseCase
}

func NewNotificationRulesController(uc *application.ManageNotificationRulesUseCase) *NotificationRulesController {
	return &NotificationRulesController{uc: uc}
}

// tenantIDFromContext lee el tenant_id que JWTAuthMiddleware ya dejo en el
// contexto de Gin (todas las rutas de reglas estan detras de ese middleware).
func tenantIDFromContext(c *gin.Context) (int, bool) {
	val, exists := c.Get("tenant_id")
	if !exists {
		return 0, false
	}
	tenantID, ok := val.(int)
	return tenantID, ok
}

func (ctrl *NotificationRulesController) List(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}
	rules, err := ctrl.uc.List(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rules)
}

func (ctrl *NotificationRulesController) GetByStateCode(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}
	code := c.Param("codigo_estado")
	rule, err := ctrl.uc.GetByStateCode(c.Request.Context(), tenantID, code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (ctrl *NotificationRulesController) Upsert(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}
	code := c.Param("codigo_estado")
	var rule domain.NotificationRule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rule.StateCode = code
	rule.TenantID = tenantID
	if err := ctrl.uc.Upsert(c.Request.Context(), tenantID, rule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (ctrl *NotificationRulesController) Delete(c *gin.Context) {
	tenantID, ok := tenantIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant no encontrado en token"})
		return
	}
	code := c.Param("codigo_estado")
	if err := ctrl.uc.Delete(c.Request.Context(), tenantID, code); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
