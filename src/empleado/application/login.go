package application

import (
	"context"
	"errors"
	"strings"

	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/empleado/domain"
	"github.com/vicpoo/API_recolecta/src/empleado/domain/entities"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

type LoginEmpleadoInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginEmpleadoOutput struct {
	Empleado *entities.Empleado `json:"empleado"`
	Token    string             `json:"token"`
}

type LoginEmpleado struct {
	repo domain.EmpleadoRepository
}

func NewLoginEmpleado(repo domain.EmpleadoRepository) *LoginEmpleado {
	return &LoginEmpleado{repo: repo}
}

func (uc *LoginEmpleado) Execute(ctx context.Context, in LoginEmpleadoInput) (*LoginEmpleadoOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	password := strings.TrimSpace(in.Password)

	if email == "" {
		return nil, errors.New("email es requerido")
	}

	if password == "" {
		return nil, errors.New("password es requerido")
	}

	empleado, err := uc.repo.FindByMail(ctx, email)
	if err != nil {
		return nil, err
	}

	if empleado == nil {
		return nil, errors.New("credenciales inválidas")
	}

	if empleado.Desactivado {
		return nil, errors.New("empleado desactivado")
	}

	if err := passwordSecurity.Verify(empleado.Password, password); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := core.GenerateToken(empleado.ID, empleado.RolID, empleado.TenantID)
	if err != nil {
		return nil, err
	}

	return &LoginEmpleadoOutput{
		Empleado: empleado,
		Token:    token,
	}, nil
}
