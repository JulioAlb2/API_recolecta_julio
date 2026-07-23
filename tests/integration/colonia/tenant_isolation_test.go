package colonia_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vicpoo/API_recolecta/src/core"
)

// TestColoniaTenantIsolation prueba de extremo a extremo que la Row-Level
// Security de Fase 5 aisla tenants entre si sobre la tabla colonia.
//
// El rol de conexion configurado hoy (DB_USER) es superusuario en este
// entorno de desarrollo -- ver docs/08-multitenancy-implementado.md -- y los
// superusuarios ignoran RLS sin importar FORCE ROW LEVEL SECURITY. Por eso
// esta prueba crea su propio rol restringido (no superusuario) para verificar
// que las politicas en si funcionan como se diseñaron, independientemente de
// que el rol configurado actualmente en .env no las este aprovechando.
func TestColoniaTenantIsolation(t *testing.T) {
	ctx := context.Background()
	adminPool := core.GetBD()

	const testRole = "tenant_isolation_test_role"
	const testRolePassword = "tenant_isolation_test_pw"
	const tenantA = 9001
	const tenantB = 9002

	dropTestRole := func() {
		// DROP ROLE falla si el rol todavia tiene privilegios GRANTed
		// pendientes (sobre colonia y su secuencia) -- DROP OWNED BY los
		// revoca antes de poder soltar el rol.
		_, _ = adminPool.Exec(ctx, "DROP OWNED BY "+testRole)
		_, _ = adminPool.Exec(ctx, "DROP ROLE IF EXISTS "+testRole)
	}

	dropTestRole()
	if _, err := adminPool.Exec(ctx, fmt.Sprintf("CREATE ROLE %s LOGIN PASSWORD '%s'", testRole, testRolePassword)); err != nil {
		t.Fatalf("no se pudo crear rol de prueba: %v", err)
	}
	if _, err := adminPool.Exec(ctx, "GRANT SELECT, INSERT, UPDATE ON colonia TO "+testRole); err != nil {
		t.Fatalf("no se pudo otorgar permisos sobre colonia: %v", err)
	}
	if _, err := adminPool.Exec(ctx, "GRANT USAGE, SELECT ON SEQUENCE colonia_colonia_id_seq TO "+testRole); err != nil {
		t.Fatalf("no se pudo otorgar permisos de secuencia: %v", err)
	}
	t.Cleanup(dropTestRole)

	if _, err := adminPool.Exec(ctx,
		"INSERT INTO tenant (tenant_id, nombre) VALUES ($1, 'Tenant Test A'), ($2, 'Tenant Test B') ON CONFLICT (tenant_id) DO NOTHING",
		tenantA, tenantB,
	); err != nil {
		t.Fatalf("no se pudo sembrar tenants de prueba: %v", err)
	}
	t.Cleanup(func() {
		_, _ = adminPool.Exec(ctx, "DELETE FROM colonia WHERE tenant_id IN ($1,$2)", tenantA, tenantB)
		_, _ = adminPool.Exec(ctx, "DELETE FROM tenant WHERE tenant_id IN ($1,$2)", tenantA, tenantB)
	})

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		testRole, testRolePassword, os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	restrictedPool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("no se pudo conectar con el rol restringido: %v", err)
	}
	defer restrictedPool.Close()

	var coloniaID int
	err = core.RunInTenantTx(ctx, restrictedPool, tenantA, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			"INSERT INTO colonia (nombre, zona, created_at, tenant_id) VALUES ($1,$2,NOW(),$3) RETURNING colonia_id",
			"Colonia Aislamiento Test", "test", tenantA,
		).Scan(&coloniaID)
	})
	if err != nil {
		t.Fatalf("no se pudo crear colonia como tenant A: %v", err)
	}
	t.Cleanup(func() {
		_, _ = adminPool.Exec(ctx, "DELETE FROM colonia WHERE colonia_id = $1", coloniaID)
	})

	// Tenant B no debe poder ver la colonia de tenant A.
	var countAsB int
	err = core.RunInTenantTx(ctx, restrictedPool, tenantB, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "SELECT COUNT(*) FROM colonia WHERE colonia_id = $1", coloniaID).Scan(&countAsB)
	})
	if err != nil {
		t.Fatalf("query como tenant B fallo: %v", err)
	}
	if countAsB != 0 {
		t.Fatalf("tenant B no deberia poder ver la colonia de tenant A, pero countAsB=%d", countAsB)
	}

	// Tenant A si debe poder ver su propia colonia.
	var countAsA int
	err = core.RunInTenantTx(ctx, restrictedPool, tenantA, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx, "SELECT COUNT(*) FROM colonia WHERE colonia_id = $1", coloniaID).Scan(&countAsA)
	})
	if err != nil {
		t.Fatalf("query como tenant A fallo: %v", err)
	}
	if countAsA != 1 {
		t.Fatalf("tenant A deberia poder ver su propia colonia, countAsA=%d", countAsA)
	}

	// Tenant B no debe poder modificarla (RLS filtra el UPDATE a 0 filas).
	err = core.RunInTenantTx(ctx, restrictedPool, tenantB, func(tx pgx.Tx) error {
		cmdTag, err := tx.Exec(ctx, "UPDATE colonia SET zona = 'hackeado' WHERE colonia_id = $1", coloniaID)
		if err != nil {
			return err
		}
		if cmdTag.RowsAffected() != 0 {
			return fmt.Errorf("tenant B pudo modificar %d fila(s) de la colonia de tenant A", cmdTag.RowsAffected())
		}
		return nil
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
}
