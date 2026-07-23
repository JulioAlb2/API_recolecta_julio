package application

import (
	"context"
	"time"

	"github.com/vicpoo/API_recolecta/src/colonia/domain"
)

type CreateColonia struct {
	repo domain.ColoniaRepository
}

func NewCreateColonia(repo domain.ColoniaRepository) *CreateColonia {
	return &CreateColonia{repo}
}

func (uc *CreateColonia) Execute(ctx context.Context, tenantID int, c *domain.Colonia) error {
	c.CreatedAt = time.Now()
	return uc.repo.Create(ctx, tenantID, c)
}
