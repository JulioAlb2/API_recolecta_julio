package postgres

import (
	"context"

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

func (r *PostgresColoniaRepository) Create(c *domain.Colonia) error {
	query := `
		INSERT INTO colonia (nombre, zona, created_at)
		VALUES ($1,$2,$3)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		c.Nombre,
		c.Zona,
		c.CreatedAt,
	)

	return err
}

func (r *PostgresColoniaRepository) GetByID(id int) (*domain.Colonia, error) {
	query := `
		SELECT id, nombre, zona, created_at
		FROM colonia
		WHERE id = $1
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
		SELECT id, nombre, zona, created_at
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

func (r *PostgresColoniaRepository) Update(c *domain.Colonia) error {
	query := `
		UPDATE colonia
		SET nombre = $1,
		    zona = $2
		WHERE id = $3
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		c.Nombre,
		c.Zona,
		c.ColoniaID,
	)

	return err
}

func (r *PostgresColoniaRepository) Delete(id int) error {
	query := `DELETE FROM colonia WHERE id = $1`

	_, err := r.db.Exec(context.Background(), query, id)
	return err
}
