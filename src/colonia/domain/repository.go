package domain

import "context"

type ColoniaRepository interface {
	Create(ctx context.Context, tenantID int, c *Colonia) error
	GetByID(id int) (*Colonia, error)
	GetAll() ([]Colonia, error)
	Update(ctx context.Context, tenantID int, c *Colonia) error
	Delete(ctx context.Context, tenantID int, id int) error
}
