package domain

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type DomicilioRepository interface {
	Create(ctx context.Context, d *entities.Domicilio) (int, error)
	GetByID(ctx context.Context, id int) (*entities.Domicilio, error)
	List(ctx context.Context) ([]entities.Domicilio, error)
	ListByCiudadanoID(ctx context.Context, ciudadanoID int) ([]entities.Domicilio, error)
	Update(ctx context.Context, d *entities.Domicilio) error
	DeleteByCiudadano(ctx context.Context, id int, ciudadanoID int) error
	FindByAlias(ctx context.Context, alias string, ciudadanoID int) (*entities.Domicilio, error)
}