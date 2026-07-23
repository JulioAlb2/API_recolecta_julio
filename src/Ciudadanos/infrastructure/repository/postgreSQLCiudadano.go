package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/domain/entities"
)

type CiudadanoPostgresRepository struct {
	db *pgxpool.Pool
}

func NewCiudadanoPostgresRepository(db *pgxpool.Pool) *CiudadanoPostgresRepository {
	return &CiudadanoPostgresRepository{db: db}
}

func (r *CiudadanoPostgresRepository) Create(ctx context.Context, c *entities.Ciudadano) (int, error) {
	const q = `
		INSERT INTO ciudadano (email, alias, password, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var id int
	err := r.db.QueryRow(ctx, q, c.Email, c.Alias, c.Password, c.CreatedAt).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *CiudadanoPostgresRepository) GetByID(ctx context.Context, id int) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE id = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, id).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CiudadanoPostgresRepository) List(ctx context.Context) ([]entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		ORDER BY id DESC
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []entities.Ciudadano

	for rows.Next() {
		var c entities.Ciudadano
		if err := rows.Scan(
			&c.ID,
			&c.TenantID,
			&c.Email,
			&c.Alias,
			&c.Password,
			&c.CreatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	return out, rows.Err()
}

func (r *CiudadanoPostgresRepository) Update(ctx context.Context, c *entities.Ciudadano) error {
	const q = `
		UPDATE ciudadano
		SET email = $1,
		    alias = $2,
		    password = $3
		WHERE id = $4
	`

	cmdTag, err := r.db.Exec(ctx, q, c.Email, c.Alias, c.Password, c.ID)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("ciudadano no encontrado")
	}

	return nil
}

func (r *CiudadanoPostgresRepository) Delete(ctx context.Context, id int) error {
	const q = `DELETE FROM ciudadano WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, q, id)
	if err != nil {
		return err
	}

	if cmdTag.RowsAffected() == 0 {
		return errors.New("ciudadano no encontrado")
	}

	return nil
}

func (r *CiudadanoPostgresRepository) FindByEmail(ctx context.Context, email string) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE email = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, email).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *CiudadanoPostgresRepository) FindByAlias(ctx context.Context, alias string) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE alias = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, alias).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}

func (r *CiudadanoPostgresRepository) FindByEmailOrAlias(ctx context.Context, value string) (*entities.Ciudadano, error) {
	const q = `
		SELECT id, tenant_id, email, alias, password, created_at
		FROM ciudadano
		WHERE email = $1 OR alias = $1
	`

	var c entities.Ciudadano
	err := r.db.QueryRow(ctx, q, value).Scan(
		&c.ID,
		&c.TenantID,
		&c.Email,
		&c.Alias,
		&c.Password,
		&c.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &c, nil
}
