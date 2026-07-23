package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/colonia/domain"
)

type DeleteColonia struct {
	repo domain.ColoniaRepository
}

func NewDeleteColonia(repo domain.ColoniaRepository) *DeleteColonia {
	return &DeleteColonia{repo}
}

func (uc *DeleteColonia) Execute(ctx context.Context, tenantID int, id int) error {
	return uc.repo.Delete(ctx, tenantID, id)
}
