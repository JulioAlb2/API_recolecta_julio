package application_ciudadano

import (
	"context"
	"errors"
	"strings"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
	"github.com/vicpoo/API_recolecta/src/core"
	passwordSecurity "github.com/vicpoo/API_recolecta/src/security/password"
)

type LoginCiudadanoInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginCiudadanoOutput struct {
	Ciudadano *entities.Ciudadano `json:"ciudadano"`
	Token     string              `json:"token"`
}

type LoginCiudadano struct {
	repo domain.CiudadanoRepository
}

func NewLoginCiudadano(repo domain.CiudadanoRepository) *LoginCiudadano {
	return &LoginCiudadano{repo: repo}
}

func (uc *LoginCiudadano) Execute(ctx context.Context, in LoginCiudadanoInput) (*LoginCiudadanoOutput, error) {
	email := strings.TrimSpace(strings.ToLower(in.Email))
	password := strings.TrimSpace(in.Password)

	if email == "" {
		return nil, errors.New("email es requerido")
	}

	if password == "" {
		return nil, errors.New("password es requerido")
	}

	ciudadano, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if ciudadano == nil {
		return nil, errors.New("credenciales inválidas")
	}

	if err := passwordSecurity.Verify(ciudadano.Password, password); err != nil {
		return nil, errors.New("credenciales inválidas")
	}

	token, err := core.GenerateToken(ciudadano.ID, 0, ciudadano.TenantID)
	if err != nil {
		return nil, err
	}

	return &LoginCiudadanoOutput{
		Ciudadano: ciudadano,
		Token:     token,
	}, nil
}
