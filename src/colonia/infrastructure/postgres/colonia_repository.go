package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/colonia/domain"
	"github.com/vicpoo/API_recolecta/src/core"
)

type PostgresColoniaRepository struct {
	db *pgxpool.Pool
}

func NewColoniaRepository() domain.ColoniaRepository {
	db := core.GetBD()

	return &PostgresColoniaRepository{
		db: db,
	}
}

func (r *PostgresColoniaRepository) Create(ctx context.Context, tenantID int, c *domain.Colonia) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			INSERT INTO colonia (nombre, zona, created_at, tenant_id)
			VALUES ($1,$2,$3,$4)
		`

		_, err := tx.Exec(ctx, query, c.Nombre, c.Zona, c.CreatedAt, tenantID)
		return err
	})
}

func (r *PostgresColoniaRepository) GetByID(id int) (*domain.Colonia, error) {
	query := `
		SELECT colonia_id, nombre, zona, created_at
		FROM colonia
		WHERE colonia_id = $1
	`

	row := r.db.QueryRow(context.Background(), query, id)

	var c domain.Colonia
	err := row.Scan(
		&c.ColoniaID,
		&c.Nombre,
		&c.Zona,
		&c.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *PostgresColoniaRepository) GetAll() ([]domain.Colonia, error) {
	query := `
		SELECT colonia_id, nombre, zona, created_at
		FROM colonia
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var colonias []domain.Colonia

	for rows.Next() {
		var c domain.Colonia
		if err := rows.Scan(
			&c.ColoniaID,
			&c.Nombre,
			&c.Zona,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}

		colonias = append(colonias, c)
	}

	return colonias, nil
}

func (r *PostgresColoniaRepository) Update(ctx context.Context, tenantID int, c *domain.Colonia) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE colonia
			SET nombre = $1,
			    zona = $2
			WHERE colonia_id = $3
		`

		_, err := tx.Exec(ctx, query, c.Nombre, c.Zona, c.ColoniaID)
		return err
	})
}

func (r *PostgresColoniaRepository) Delete(ctx context.Context, tenantID int, id int) error {
	return core.RunInTenantTx(ctx, r.db, tenantID, func(tx pgx.Tx) error {
		query := `
			UPDATE colonia
			SET eliminado = true
			WHERE colonia_id = $1
		`

		_, err := tx.Exec(ctx, query, id)
		return err
	})
}
