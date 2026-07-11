package application_domicilio

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain"
	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type CreateDomicilioInput struct {
	CiudadanoID int      `json:"ciudadano_id"`
	ColoniaID   int      `json:"colonia_id"`
	Alias       string   `json:"alias"`
	Calle       string   `json:"calle"`
	Numero      string   `json:"numero"`
	Referencia  *string  `json:"referencia,omitempty"`
	Latitud     *float64 `json:"latitud,omitempty"`
	Longitud    *float64 `json:"longitud,omitempty"`
}

type CreateDomicilio struct {
	repo domain.DomicilioRepository
}

func NewCreateDomicilio(repo domain.DomicilioRepository) *CreateDomicilio {
	return &CreateDomicilio{repo: repo}
}

func validateCoordenadas(latitud, longitud *float64) error {
	if latitud == nil && longitud == nil {
		return nil
	}
	if latitud == nil || longitud == nil {
		return errors.New("latitud y longitud deben enviarse juntas")
	}
	if *latitud < -90 || *latitud > 90 {
		return fmt.Errorf("latitud inválida: %v", *latitud)
	}
	if *longitud < -180 || *longitud > 180 {
		return fmt.Errorf("longitud inválida: %v", *longitud)
	}
	return nil
}

func (uc *CreateDomicilio) Execute(ctx context.Context, in CreateDomicilioInput) (int, error) {
	in.Alias = strings.TrimSpace(in.Alias)
	in.Calle = strings.TrimSpace(in.Calle)
	in.Numero = strings.TrimSpace(in.Numero)

	if in.CiudadanoID <= 0 {
		return 0, errors.New("ciudadano_id es requerido")
	}
	if in.ColoniaID <= 0 {
		return 0, errors.New("colonia_id es requerido")
	}
	if in.Alias == "" {
		return 0, errors.New("alias es requerido")
	}
	if in.Calle == "" {
		return 0, errors.New("calle es requerida")
	}
	if in.Numero == "" {
		return 0, errors.New("numero es requerido")
	}
	if err := validateCoordenadas(in.Latitud, in.Longitud); err != nil {
		return 0, err
	}

	existingByAlias, err := uc.repo.FindByAlias(ctx, in.Alias)
	if err != nil {
		return 0, err
	}
	if existingByAlias != nil {
		return 0, errors.New("el alias del domicilio ya está registrado")
	}

	d := &entities.Domicilio{
		CiudadanoID: in.CiudadanoID,
		ColoniaID:   in.ColoniaID,
		Alias:       in.Alias,
		Calle:       in.Calle,
		Numero:      in.Numero,
		Referencia:  in.Referencia,
		Latitud:     in.Latitud,
		Longitud:    in.Longitud,
		CreatedAt:   time.Now(),
	}

	return uc.repo.Create(ctx, d)
}
