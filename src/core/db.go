package core

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
	err  error
)

func ConnectPostgres() (*pgxpool.Pool, error) {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_NAME"),
		)

		fmt.Printf("[db-connect] Usando DSN: postgres://%s:***@%s:%s/%s?sslmode=disable\n", os.Getenv("DB_USER"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

		config, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			err = fmt.Errorf("error parsing DSN: %w", parseErr)
			return
		}

		config.MaxConns = 25
		config.MinConns = 5
		config.MaxConnLifetime = time.Hour
		config.MaxConnIdleTime = 30 * time.Minute
		config.HealthCheckPeriod = time.Minute

		var tempPool *pgxpool.Pool
		var lastErr error
		maxRetries := 5

		for i := 1; i <= maxRetries; i++ {
			fmt.Printf("[db-connect] Intentando conectar a la base de datos (intento %d de %d)...\n", i, maxRetries)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			
			tempPool, lastErr = pgxpool.NewWithConfig(ctx, config)
			if lastErr == nil {
				lastErr = tempPool.Ping(ctx)
				if lastErr == nil {
					cancel()
					pool = tempPool
					break
				}
				tempPool.Close()
			}
			cancel()

			fmt.Printf("[db-connect] Intento %d falló: %v. Reintentando en 3 segundos...\n", i, lastErr)
			time.Sleep(3 * time.Second)
		}

		if pool == nil {
			err = fmt.Errorf("no se pudo conectar a la base de datos tras %d intentos: %w", maxRetries, lastErr)
			return
		}

		// Run auto migration to auto-heal/sync outdated development schemas
		AutoMigrateDatabase(pool)
	})

	return pool, err
}

func ClosePool() {
	if pool != nil {
		pool.Close()
	}
}

func GetBD() *pgxpool.Pool {
	db, err := ConnectPostgres()
	if err != nil {
		panic(fmt.Sprintf("Error al conectar a la base de datos: %v", err))
	}
	return db
}

func AutoMigrateDatabase(db *pgxpool.Pool) {
	ctx := context.Background()

	// Only run this in development/debug mode
	if os.Getenv("ENVIRONMENT") != "development" && os.Getenv("DEBUG") != "true" {
		return
	}

	fmt.Println("[db-migrate] Checking if database schema is up-to-date...")

	// Check 1: Does table 'domicilio' have column 'calle'?
	var hasCalle bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='domicilio' AND column_name='calle')").Scan(&hasCalle)

	// Check 2: Does 'alerta_mantenimiento' table exist?
	var hasAlertaTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='alerta_mantenimiento')").Scan(&hasAlertaTable)

	// Check 3: Does 'ruta_camion' table exist?
	var hasRutaCamionTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='ruta_camion')").Scan(&hasRutaCamionTable)

	// Check 4: Does 'historial_asignacion_camion' table exist?
	var hasHistorialTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='historial_asignacion_camion')").Scan(&hasHistorialTable)

	// Check 5: Does 'registro_mantenimiento' table have 'registro_id' column?
	var hasRegistroID bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='registro_mantenimiento' AND column_name='registro_id')").Scan(&hasRegistroID)

	// Check 6: Does 'colonia' table have 'colonia_id' column?
	var hasColoniaID bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='colonia' AND column_name='colonia_id')").Scan(&hasColoniaID)

	// Check 7: Does 'relleno_sanitario' table exist?
	var hasRellenoTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='relleno_sanitario')").Scan(&hasRellenoTable)

	// Check 8: Does 'estado_camion' table exist?
	var hasEstadoCamionTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='estado_camion')").Scan(&hasEstadoCamionTable)

	// Check 9: Does 'registro_vaciado' table exist?
	var hasRegistroVaciadoTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='registro_vaciado')").Scan(&hasRegistroVaciadoTable)

	// Check 10: Does 'alerta_usuario' table exist?
	var hasAlertaUsuarioTable bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='alerta_usuario')").Scan(&hasAlertaUsuarioTable)

	// Check 11: Does 'tipo_camion' table have 'created_at' column?
	var hasTipoCamionCreatedAt bool
	_ = db.QueryRow(ctx, "SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='tipo_camion' AND column_name='created_at')").Scan(&hasTipoCamionCreatedAt)

	needsUpdate := !hasCalle || !hasAlertaTable || !hasRutaCamionTable || !hasHistorialTable || !hasRegistroID || !hasColoniaID || !hasRellenoTable || !hasEstadoCamionTable || !hasRegistroVaciadoTable || !hasAlertaUsuarioTable || !hasTipoCamionCreatedAt

	if !needsUpdate {
		fmt.Println("[db-migrate] Database schema is up-to-date. No action needed.")
		return
	}

	fmt.Println("[db-migrate] Database schema is outdated! Recreating tables...")

	// Drop old/conflicting tables and views/triggers
	dropQuery := `
		DROP TABLE IF EXISTS 
			domicilio, 
			historial_asignacion, 
			historial_asignacion_camion, 
			registro_mantenimiento, 
			tipo_mantenimiento, 
			alerta_mantenimiento, 
			registro_asignacion_ruta, 
			ruta_camion,
			colonia,
			relleno_sanitario,
			estado_camion,
			registro_vaciado,
			alerta_usuario,
			camion,
			tipo_camion
		CASCADE;
	`
	_, err := db.Exec(ctx, dropQuery)
	if err != nil {
		fmt.Printf("[db-migrate] Warning dropping tables: %v\n", err)
	}

	// Load and execute SQL files
	executeSQLFile(ctx, db, "db_script.sql")
	executeSQLFile(ctx, db, "db_constraints.sql")
	executeSQLFile(ctx, db, "db_indexes.sql")

	fmt.Println("[db-migrate] Database schema successfully updated/recreated!")
}

func executeSQLFile(ctx context.Context, db *pgxpool.Pool, filepath string) {
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		// Try parent directory fallback
		bytes, err = os.ReadFile("../" + filepath)
		if err != nil {
			fmt.Printf("[db-migrate] Error reading %s: %v\n", filepath, err)
			return
		}
	}

	content := string(bytes)
	lines := strings.Split(content, "\n")
	var cleanLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "CREATE DATABASE") || strings.HasPrefix(trimmed, "\\c") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}
	cleanContent := strings.Join(cleanLines, "\n")

	_, err = db.Exec(ctx, cleanContent)
	if err != nil {
		fmt.Printf("[db-migrate] Error executing %s: %v\n", filepath, err)
	} else {
		fmt.Printf("[db-migrate] Successfully executed %s\n", filepath)
	}
}
