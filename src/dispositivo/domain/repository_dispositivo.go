package domain

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/dispositivo/domain/entities"
)

type DispositivoRepository interface {
	Solicitar(ctx context.Context, dispositivo *entities.Dispositivo) error
	FindByConductorID(ctx context.Context, conductorID int) (*entities.Dispositivo, error)
	Aprobar(ctx context.Context, conductorID int) error
	Desvincular(ctx context.Context, conductorID int) error
	ListarPendientes(ctx context.Context) ([]*entities.DispositivoConductorResponse, error)
}
