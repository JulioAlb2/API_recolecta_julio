# API Recolecta — Flujo de Fallas (Anomalías)

Este documento describe el flujo del módulo **Fallas**, que reemplaza las
cinco entidades anteriores (`Anomalia`, `Incidencia`, `ReporteConductor`,
`ReporteFallaCritica`, `SeguimientoFallaCritica`) por una **única entidad
unificada `Anomalia`**, diferenciada mediante el campo `tipo_anomalia`.

> El módulo ya compila (`go build ./...` pasa) y las entidades/controladores
> legado fueron eliminados por completo. Ver la sección
> [Estado actual / pendientes](#estado-actual--pendientes) al final para lo
> que todavía falta pulir (no son errores de compilación, son
> inconsistencias de comportamiento).

## Base URL

- Local: `http://localhost:8080`
- Prefijo: `/api/anomalias`

## Autenticación y roles

Todas las rutas de `/api/anomalias` requieren JWT y uno de estos roles:

- `ADMIN`
- `SUPERVISOR`
- `COORDINADOR`

```
Authorization: Bearer <token>
```

Sin token válido → `401 UNAUTHORIZED`. Con token pero rol no permitido →
`403 FORBIDDEN`.

---

## Por qué se unificó el modelo

Antes existían 5 tablas/entidades independientes con CRUD casi idéntico
(`Incidencia`, `ReporteConductor`, `ReporteFallaCritica`,
`SeguimientoFallaCritica`, y la propia `Anomalia`). Ahora todo vive en una
sola tabla `anomalia`, y el campo `tipo_anomalia` indica qué "concepto"
representa cada fila:

| `tipo_anomalia`             | Reemplaza a                | Notas |
|---|---|---|
| `ANOMALIA`                  | Anomalia (general)          | Caso genérico |
| `INCIDENCIA`                | Incidencia                  | Reportada en un punto de recolección |
| `REPORTE_CONDUCTOR`         | ReporteConductor             | Reportada por un conductor |
| `REPORTE_FALLA_CRITICA`     | ReporteFallaCritica           | Falla grave de camión/ruta |
| `SEGUIMIENTO_FALLA_CRITICA` | SeguimientoFallaCritica        | Comentario/seguimiento sobre un `REPORTE_FALLA_CRITICA` existente |

`SEGUIMIENTO_FALLA_CRITICA` es un caso especial: **no** crea un nuevo
incidente, sino que usa `anomalia_referencia_id` para apuntar al
`anomalia_id` del `REPORTE_FALLA_CRITICA` al que da seguimiento (auto
relación dentro de la misma tabla `anomalia`).

El módulo **Mantenimiento** completo también fue eliminado en esta
migración (`src/Mantenimiento/...` y su registro en `dependencies.go`); no
quedan referencias residuales en el código fuente (`docs/docs.go`, generado
por swagger, todavía menciona `AlertaMantenimiento` porque no se ha
regenerado, pero no afecta la compilación ni el runtime).

---

## `tipo_anomalia`: enum real, no solo un string

`TipoAnomalia` (`src/Fallas/domain/entities/tipo_anomalia.go`) dejó de ser
un simple alias de `string` y ahora es un **enum real** respaldado por
`int`:

```go
type TipoAnomalia int

const (
	TipoAnomaliaAnomalia TipoAnomalia = iota
	TipoAnomaliaIncidencia
	TipoAnomaliaReporteConductor
	TipoAnomaliaReporteFallaCritica
	TipoAnomaliaSeguimientoFallaCritica
)
```

Esto es transparente para los consumidores de la API porque el tipo
implementa:

- `MarshalJSON` / `UnmarshalJSON` → en JSON siempre se ve/envía como texto
  (`"REPORTE_FALLA_CRITICA"`, etc.), nunca como número.
- `Value` / `Scan` (`driver.Valuer` / `sql.Scanner`) → se persiste y se lee
  de Postgres como texto también.
- `String()` → representación en texto.
- `ParseTipoAnomalia(valor string) (TipoAnomalia, error)` → valida y
  convierte un string a enum; es lo que usan los controladores para
  validar `tipo_anomalia` antes de tocar la base de datos.
- `EsTipoAnomaliaValido(tipo string) bool` y `TiposAnomaliaValidos() []string`
  se mantienen para compatibilidad.

En la base de datos, la columna `anomalia.tipo_anomalia` pasó de ser
`VARCHAR(50)` + `CHECK` a un **enum nativo de Postgres**
(`tipo_anomalia_enum`), creado y migrado de forma idempotente en
`db_script.sql` / `db_constraints.sql`.

---

## Modelo de datos: `Anomalia`

```json
{
  "anomalia_id": 12,
  "tipo_anomalia": "REPORTE_FALLA_CRITICA",
  "punto_id": 5,
  "conductor_id": 3,
  "camion_id": 2,
  "ruta_id": 7,
  "anomalia_referencia_id": null,
  "descripcion": "Fuga de aceite en el motor",
  "json_ruta": "",
  "estado": "PENDIENTE",
  "eliminado": false,
  "fecha_reporte": "2026-07-14T08:30:00Z",
  "fecha_resolucion": null,
  "created_at": "2026-07-14T08:31:02Z",
  "updated_at": "2026-07-14T08:31:02Z"
}
```

| Campo | Tipo | Notas |
|---|---|---|
| `anomalia_id` | int32 | PK, autogenerado |
| `tipo_anomalia` | string (enum real en Go, `tipo_anomalia_enum` en BD) | Uno de los 5 valores de la tabla anterior |
| `punto_id` | int32 \| null | Punto de recolección relacionado |
| `conductor_id` | int32 \| null | Conductor relacionado |
| `camion_id` | int32 \| null | Camión relacionado |
| `ruta_id` | int32 \| null | Ruta relacionada |
| `anomalia_referencia_id` | int32 \| null | Solo se usa en `SEGUIMIENTO_FALLA_CRITICA`: apunta a otra `anomalia_id` |
| `descripcion` | string | Requerido |
| `json_ruta` | string | Opcional, texto libre (JSON serializado de una ruta/ubicación) |
| `estado` | string | `PENDIENTE`, `EN_PROCESO`, `RESUELTA`, u otro definido por el cliente |
| `eliminado` | bool | Borrado lógico |
| `fecha_reporte` | datetime | Requerido; si se omite/está vacío al crear, se usa la fecha/hora actual |
| `fecha_resolucion` | datetime \| null | Se llena al resolver |
| `created_at` / `updated_at` | datetime | Gestionados por el servidor |

El **borrado es lógico**: `DELETE` no borra la fila, marca
`eliminado = true`. Esto es intencional porque una fila
`SEGUIMIENTO_FALLA_CRITICA` puede referenciar a otra vía
`anomalia_referencia_id`, y borrar físicamente rompería esa relación.
`GetAll` y todos los `GetBy*` filtran siempre `eliminado = false`.

---

## Respuesta de error estandarizada

Todas las rutas del backend (no solo Fallas) usan el mismo formato de error:

```json
{
  "error": {
    "code": "NOT_FOUND",
    "message": "Anomalía no encontrado",
    "details": { "identifier": "12" }
  }
}
```

| Código | HTTP | Cuándo ocurre en este módulo |
|---|---:|---|
| `VALIDATION_ERROR` | 400 | Falta un campo requerido en el body, `tipo_anomalia` no es uno de los 5 valores válidos, o falta un query param requerido (`estado`, `tipo_anomalia`, `fecha_inicio`/`fecha_fin`) |
| `BAD_REQUEST` | 400 | Formato de fecha inválido en el body (`fecha_reporte` / `fecha_resolucion`) |
| `INVALID_INPUT` | 400 | El `id`/`puntoId`/`choferId`/`camionId`/`rutaId`/`referenciaId` de la URL no es un entero válido |
| `NOT_FOUND` | 404 | El `anomalia_id` no existe — **solo confirmado en `GET /{id}`**, ver aviso abajo |
| `INTERNAL_SERVER_ERROR` | 500 | Error no controlado, incluye hoy `PUT`/`DELETE` sobre un `id` inexistente (ver pendientes) |

---

## Endpoints

### 1. Crear anomalía

`POST /api/anomalias/`

**Body:**

```json
{
  "tipo_anomalia": "REPORTE_FALLA_CRITICA",
  "punto_id": null,
  "conductor_id": 3,
  "camion_id": 2,
  "ruta_id": 7,
  "anomalia_referencia_id": null,
  "descripcion": "Fuga de aceite en el motor",
  "json_ruta": "",
  "estado": "PENDIENTE",
  "fecha_reporte": "2026-07-14 08:30:00"
}
```

- Requeridos: `tipo_anomalia`, `descripcion`, `estado`, `fecha_reporte`.
- `tipo_anomalia` debe ser uno de los 5 valores del enum; si no lo es, se
  responde `400 VALIDATION_ERROR` con el detalle en `tipo_anomalia`.
- `fecha_reporte` acepta los formatos `YYYY-MM-DD HH:MM:SS`,
  `YYYY-MM-DDTHH:MM:SSZ` o `YYYY-MM-DD` (ver `parseFecha` en
  `src/Fallas/infrastructure/helpers.go`); cualquier otro formato (p. ej.
  con offset de zona horaria `+00:00`) responde `400 BAD_REQUEST`.

**Respuesta 201** — el objeto `Anomalia` creado (ver modelo de datos arriba).

**Errores:**
- `400 VALIDATION_ERROR` — falta `tipo_anomalia`, `descripcion`, `estado`
  o `fecha_reporte`; o `tipo_anomalia` inválido.
- `400 BAD_REQUEST` — `fecha_reporte` con formato no reconocido.

---

### 2. Listar todas

`GET /api/anomalias/`

**Respuesta 200:**

```json
{ "data": [ { "anomalia_id": 1, "...": "..." } ] }
```

Solo devuelve filas con `eliminado = false`, ordenadas por
`fecha_reporte DESC`.

---

### 3. Obtener por ID

`GET /api/anomalias/{id}`

**Respuesta 200:** objeto `Anomalia`.
**Errores:** `400 INVALID_INPUT` (id no numérico) · `404 NOT_FOUND` (no existe).

---

### 4. Actualizar

`PUT /api/anomalias/{id}`

**Body:**

```json
{
  "tipo_anomalia": "REPORTE_FALLA_CRITICA",
  "punto_id": null,
  "conductor_id": 3,
  "camion_id": 2,
  "ruta_id": 7,
  "anomalia_referencia_id": null,
  "descripcion": "Fuga de aceite reparada",
  "json_ruta": "",
  "estado": "RESUELTA",
  "eliminado": false,
  "fecha_reporte": "2026-07-14 08:30:00",
  "fecha_resolucion": "2026-07-15 10:00:00"
}
```

- Requeridos: `tipo_anomalia`, `descripcion`, `estado`, `fecha_reporte`.
- `fecha_resolucion` es opcional (string en alguno de los formatos de
  `parseFecha`, o `null`/omitido).

**Respuesta 200:** objeto `Anomalia` actualizado.
**Errores:** `400 INVALID_INPUT` (id no numérico) · `400 VALIDATION_ERROR`
(body inválido o `tipo_anomalia` inválido) · `400 BAD_REQUEST` (fechas mal
formateadas) · `404 NOT_FOUND` documentado, pero **hoy no se dispara** — un
`id` inexistente responde `500 INTERNAL_SERVER_ERROR` (ver pendientes).

---

### 5. Eliminar (borrado lógico)

`DELETE /api/anomalias/{id}`

**Respuesta 200:**

```json
{ "message": "Anomalía eliminada exitosamente" }
```

**Errores:** `400 INVALID_INPUT` (id no numérico) · `404 NOT_FOUND`
documentado, pero **hoy no se dispara** — un `id` inexistente responde
`500 INTERNAL_SERVER_ERROR` (ver pendientes).

---

### 6. Por punto de recolección

`GET /api/anomalias/punto/{puntoId}`

**Respuesta 200:** arreglo de `Anomalia` (sin envolver en `data`).
**Errores:** `400 INVALID_INPUT` (id no numérico).

---

### 7. Por conductor/chofer

`GET /api/anomalias/chofer/{choferId}`

Filtra por `conductor_id`. **Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 INVALID_INPUT` (id no numérico).

---

### 8. Por camión

`GET /api/anomalias/camion/{camionId}`

Filtra por `camion_id`. **Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 INVALID_INPUT` (id no numérico).

---

### 9. Por ruta

`GET /api/anomalias/ruta/{rutaId}`

Filtra por `ruta_id`. **Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 INVALID_INPUT` (id no numérico).

---

### 10. Por anomalía de referencia (seguimientos)

`GET /api/anomalias/referencia/{referenciaId}`

Devuelve los registros cuyo `anomalia_referencia_id` apunta a
`{referenciaId}` — típicamente los `SEGUIMIENTO_FALLA_CRITICA` de un
`REPORTE_FALLA_CRITICA`. **Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 INVALID_INPUT` (id no numérico).

---

### 11. Por estado

`GET /api/anomalias/estado?estado=PENDIENTE`

**Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 VALIDATION_ERROR` (falta `estado`).

---

### 12. Por tipo de anomalía

`GET /api/anomalias/tipo?tipo_anomalia=REPORTE_FALLA_CRITICA`

**Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 VALIDATION_ERROR` (falta `tipo_anomalia` o no es uno de
los 5 valores válidos).

---

### 13. Por rango de fechas

`GET /api/anomalias/por-fecha?fecha_inicio=2026-07-01&fecha_fin=2026-07-15`

**Respuesta 200:** arreglo de `Anomalia`.
**Errores:** `400 VALIDATION_ERROR` (falta `fecha_inicio` o `fecha_fin`).
Si el formato de fecha es inválido, hoy responde `500 INTERNAL_SERVER_ERROR`
en vez de `400 BAD_REQUEST` (el controlador no distingue ese caso).

---

## Nota sobre formas de respuesta inconsistentes

Hoy la forma de la respuesta 200 **no es uniforme** entre endpoints:

| Endpoint | Forma real de la respuesta |
|---|---|
| Crear / Obtener por ID / Actualizar | objeto `Anomalia` "pelado" |
| Listar todas | `{ "data": [...] }` |
| Eliminar | `{ "message": "..." }` |
| Por punto / chofer / camión / ruta / referencia / estado / tipo / fecha | arreglo `[...]` "pelado" |

Los comentarios `@Success` de Swagger dicen que todas devuelven
`AnomaliaResponse` / `AnomaliaListResponse` (con `success`, `message`,
`data`, `code`), pero el código real (`core.RespondOK`, `core.RespondCreated`,
`c.JSON(...)`) no arma ese envoltorio. Si el frontend/mobile va a
consumir esto, conviene decidir una única forma y alinear código + swagger.

De forma similar, `CreateAnomaliaController` y `UpdateAnomaliaController`
no hacen `ShouldBindJSON` directo sobre `entities.CreateAnomaliaRequest` /
`entities.UpdateAnomaliaRequest` (los structs documentados en Swagger,
definidos en `anomalia_swagger.go`); usan un `struct` anónimo local con
campos ligeramente distintos (`fecha_reporte` como `string` requerido en
vez de `time.Time` opcional, `estado` requerido). El contrato real es el
de los ejemplos de este documento, no el de los structs de Swagger.

---

## Estado actual / pendientes

El proyecto **compila** (`go build ./...` sin errores) y los archivos
legado (`Incidencia`, `ReporteConductor`, `ReporteFallaCritica`,
`SeguimientoFallaCritica`, `Mantenimiento`) ya fueron eliminados. Lo que
queda pendiente es afinar comportamiento, no arreglar el build:

1. **`PUT` y `DELETE` no devuelven `404` con un `id` inexistente.**
   `PostgresAnomaliaRepository.Update` / `.Delete` devuelven el error
   `fmt.Errorf("anomalía con ID %d no encontrada", id)`, pero
   `UpdateAnomaliaController` y `DeleteAnomaliaController` siguen
   comparando contra el string literal `"anomalia not found"` (que nunca
   se produce). Resultado: hoy devuelven `500 INTERNAL_SERVER_ERROR` en
   vez de `404 NOT_FOUND`. `GetAnomaliaByIdController` sí funciona
   correctamente porque además compara `anomalia == nil` como respaldo.
   Lo más simple es exportar un error tipado
   (p. ej. `var ErrAnomaliaNotFound = errors.New(...)`) y compararlo con
   `errors.Is` en los tres controladores.

2. **Unificar la forma de la respuesta 200** entre endpoints (objeto
   pelado vs. `{ "data": [...] }` vs. `{ "message": "..." }`) y alinearla
   con los comentarios `@Success` de Swagger, o actualizar Swagger para
   reflejar la realidad.

3. **Alinear `CreateAnomaliaController` / `UpdateAnomaliaController`** con
   `entities.CreateAnomaliaRequest` / `entities.UpdateAnomaliaRequest` (o
   viceversa, actualizar esos structs de Swagger) — hoy usan un `struct`
   anónimo local con reglas de requerido distintas a las documentadas en
   `anomalia_swagger.go`.

4. **`GetAnomaliasByFechaRangeController`** trata cualquier error del caso
   de uso (incluido un formato de fecha inválido) como
   `500 INTERNAL_SERVER_ERROR`; sería más correcto distinguir el caso de
   parseo inválido y responder `400 BAD_REQUEST`.

5. `docs/docs.go` (generado por `swag`) todavía referencia
   `AlertaMantenimiento` y otros modelos del módulo Mantenimiento ya
   eliminado; conviene regenerar la documentación Swagger (`swag init`)
   para que quede en sync con el código actual.
