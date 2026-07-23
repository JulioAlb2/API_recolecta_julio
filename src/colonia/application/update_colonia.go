package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/colonia/domain"
)

type UpdateColonia struct {
	repo domain.ColoniaRepository
}

func NewUpdateColonia(repo domain.ColoniaRepository) *UpdateColonia {
	return &UpdateColonia{repo}
}

func (uc *UpdateColonia) Execute(ctx context.Context, tenantID int, c *domain.Colonia) error {
	return uc.repo.Update(ctx, tenantID, c)
}
