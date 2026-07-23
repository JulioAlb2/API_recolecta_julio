-- Agrega columnas para el pipeline modelo_reportes -> clasificador_reportes
-- sobre la tabla anomalia (reportes de ciudadanos/conductores).
--
-- Correr manualmente contra la BD que ya esta corriendo (el contenedor de
-- postgres solo ejecuta db_script.sql en el primer arranque del volumen, no
-- en cada `docker compose up`):
--
--   docker exec -i postgres_db psql -U <DB_USER> -d <DB_NAME> < migrations/2026-07-22_pipeline_reportes_anomalia.sql
--
-- (los valores de DB_USER/DB_NAME estan en tu .env). Es idempotente: se
-- puede correr varias veces sin romper nada gracias a IF NOT EXISTS.

ALTER TABLE anomalia
  ADD COLUMN IF NOT EXISTS estado_pipeline VARCHAR(30) NOT NULL DEFAULT 'pendiente',
  ADD COLUMN IF NOT EXISTS nivel_riesgo VARCHAR(20),
  ADD COLUMN IF NOT EXISTS inferencia_id INTEGER,
  ADD COLUMN IF NOT EXISTS categoria_clasificada VARCHAR(50),
  ADD COLUMN IF NOT EXISTS subtipo_clasificado VARCHAR(50),
  ADD COLUMN IF NOT EXISTS accion_sugerida VARCHAR(50),
  ADD COLUMN IF NOT EXISTS pipeline_error TEXT;

COMMENT ON COLUMN anomalia.estado_pipeline IS
  'pendiente | rechazado (fraude segun modelo_reportes) | clasificado (paso por clasificador_reportes) | error (algun microservicio fallo/timeout)';
COMMENT ON COLUMN anomalia.accion_sugerida IS
  'Accion devuelta por clasificador_reportes: block_edge | inflate_weight | marcar_mantenimiento | none. Input futuro del algoritmo genetico de rutas.';
