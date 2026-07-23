package core

import (
	"context"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// RunInTenantTx abre una transaccion, fija app.current_tenant con is_local=true
// (solo para esta transaccion, no la sesion/conexion completa) y corre fn.
// pgxpool reutiliza conexiones fisicas entre requests de tenants distintos: si
// la variable quedara fijada a nivel de sesion, una conexion podria arrastrar
// el tenant de un request anterior al siguiente. Acotarla a la transaccion
// hace que Postgres la resetee sola en el COMMIT/ROLLBACK.
func RunInTenantTx(ctx context.Context, pool *pgxpool.Pool, tenantID int, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, "SELECT set_config('app.current_tenant', $1, true)", strconv.Itoa(tenantID)); err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
