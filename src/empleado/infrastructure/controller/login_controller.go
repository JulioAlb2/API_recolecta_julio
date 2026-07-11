package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/application"
	_ "github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
)

type LoginEmpleadoController struct {
	useCase *application.LoginEmpleado
}

func NewLoginEmpleadoController(useCase *application.LoginEmpleado) *LoginEmpleadoController {
	return &LoginEmpleadoController{useCase: useCase}
}

// @Summary      Login de empleado
// @Security
// @Tags         Empleado
// @Accept       json
// @Produce      json
// @Param        body body entities.LoginEmpleadoRequest true "Credenciales de login"
// @Success      200 {object} entities.LoginEmpleadoResponse
// @Failure      400 {object} core.ErrorResponse
// @Failure      401 {object} core.ErrorResponse
// @Router       /api/empleados/login [post]
func (c *LoginEmpleadoController) Run(ctx *gin.Context) {
	var input application.LoginEmpleadoInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		core.RespondBadRequest(ctx, "json inválido", err.Error())
		return
	}

	result, err := c.useCase.Execute(ctx.Request.Context(), input)
	if err != nil {
		if err.Error() == "credenciales inválidas" || err.Error() == "empleado desactivado" {
			core.RespondError(ctx, http.StatusUnauthorized, core.ErrCodeUnauthorized, err.Error(), nil)
			return
		}
		core.RespondBadRequest(ctx, err.Error(), nil)
		return
	}

	core.RespondOK(ctx, gin.H{
		"message": "login correcto",
		"token":   result.Token,
		"data":    result.Empleado,
	})
}
