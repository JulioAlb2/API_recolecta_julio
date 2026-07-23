package controller_ciudadano

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateFCMTokenController struct {
	useCase *application_ciudadano.UpdateFCMToken
}

func NewUpdateFCMTokenController(useCase *application_ciudadano.UpdateFCMToken) *UpdateFCMTokenController {
	return &UpdateFCMTokenController{useCase: useCase}
}

// @Summary      Actualizar token FCM de ciudadano
// @Description  Actualiza u overwrita el token FCM del ciudadano autenticado en Redis.
// @Tags         Ciudadano
// @Accept       json
// @Produce      json
// @Param        body body application_ciudadano.UpdateFCMTokenInput true "FCM Token Body"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} core.ErrorResponse
// @Failure      401 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Security     BearerAuth
// @Router       /api/ciudadanos/fcm-token [patch]
func (ctrl *UpdateFCMTokenController) Run(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "no autorizado", nil)
		return
	}

	userID, ok := userIDVal.(int)
	if !ok {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "no autorizado", nil)
		return
	}

	var input application_ciudadano.UpdateFCMTokenInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", map[string]string{"detail": err.Error()})
		return
	}

	err := ctrl.useCase.Execute(ctx.Request.Context(), userID, input)
	if err != nil {
		core.RespondInternalServerError(ctx, "error interno al actualizar fcm token", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "token FCM actualizado correctamente",
	})
}
