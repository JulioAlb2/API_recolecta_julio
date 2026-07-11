package controller_ciudadano

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
)

type LoginCiudadanoController struct {
	useCase *application_ciudadano.LoginCiudadano
}

func NewLoginCiudadanoController(useCase *application_ciudadano.LoginCiudadano) *LoginCiudadanoController {
	return &LoginCiudadanoController{useCase: useCase}
}

// @Summary      Login de ciudadano
// @Security
// @Tags         Ciudadano
// @Accept       json
// @Produce      json
// @Param        body body entities.LoginCiudadanoRequest true "Body"
// @Success      200 {object} entities.LoginCiudadanoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      401 {object} core.ErrorResponse
// @Failure      500 {object} core.ErrorResponse
// @Router       /api/ciudadanos/login [post]
func (c *LoginCiudadanoController) Run(ctx *gin.Context) {
	var input entities.LoginCiudadanoRequest

	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", map[string]string{"detail": err.Error()})
		return
	}

	appInput := application_ciudadano.LoginCiudadanoInput{
		Email:    input.Email,
		Password: input.Password,
	}

	result, err := c.useCase.Execute(ctx.Request.Context(), appInput)
	if err != nil {
		core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, "credenciales inválidas", nil)
		return
	}

	core.RespondOK(ctx, entities.LoginCiudadanoResponse{
		Success: true,
		Message: "login correcto",
		Token:   result.Token,
		Data:    *result.Ciudadano,
		Code:    http.StatusOK,
	})
}
