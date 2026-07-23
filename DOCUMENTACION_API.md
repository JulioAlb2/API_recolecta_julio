# Documentación técnica — API Recolecta

> Generado a partir de un análisis estático de todo el repositorio (`src/`, `docs/`, `db_*.sql`, `main.go`, `dependencies.go`) el 2026-07-21.
> Objetivo: (1) documentar todos los endpoints realmente expuestos por la API, (2) diccionario de datos completo, (3) auditar si Swagger documenta **todas** las rutas reales.

---

## 1. Resumen ejecutivo

- El proyecto es un backend **Go + Gin**, arquitectura hexagonal (`domain` / `application` / `infrastructure`) por módulo de negocio dentro de `src/`.
- Todas las rutas se registran manualmente en **[dependencies.go](dependencies.go)** (`InitDependencies()`), que arma cada caso de uso → controller → grupo de rutas y las cuelga sobre una única instancia de `gin.Engine`.
- Swagger se sirve en `/api/swagger/*any` (paquete `docs/`, generado con `swag init`).
- **Hallazgo principal:** el archivo Swagger generado (`docs/docs.go`, `docs/swagger.json`, `docs/swagger.yaml`) **está desactualizado**. No refleja el estado actual del código: le faltan rutas que sí existen y sí tienen anotación `@Router`, y sobran ~40 rutas de un módulo "Fallas y Mantenimiento" que ya no existe en el código (fue reemplazado por el módulo unificado `anomalia`). Ver sección 5.
- **Segundo hallazgo importante:** dos módulos completos (`rol`, `notificacion`) tienen dominio, casos de uso, repositorios y controladores **completos**, pero **nunca se registran** en `dependencies.go` → sus rutas **no existen en el servidor en ejecución**, aunque Swagger sí las documenta (para `rol`) o el código las anota con `@Router` sin que aparezcan en el swagger generado (para `notificacion`/push).

Total de rutas realmente montadas y accesibles: **≈93** (14 grupos/routers `.Run()`/`.RegisterRoutes()` invocados en `dependencies.go`).

---

## 2. Arquitectura y arranque

```
main.go  →  InitDependencies()  (dependencies.go)
              │
              ├─ gin.Default() + CORSMiddleware()
              ├─ GET  /api/swagger/*any   (swagger UI, paquete docs/)
              ├─ tipo-camion, camion, estado-camion, historial-asignacion,
              │  ruta, puntos-recoleccion, relleno-sanitario, ruta-camion,
              │  registro-vaciado, colonia   → registrados explícitamente
              ├─ Ciudadanos (ciudadano + domicilio), Empleado
              ├─ alerta_usuario  (bajo engine.Group("/api"))
              ├─ anomalia (Fallas)
              └─ engine.Run(":8080")
```

Módulos con código completo (dominio/aplicación/infra) que **NO se invocan** desde `dependencies.go` y por tanto no exponen ningún endpoint en el binario real:

| Módulo | Motivo |
|---|---|
| `src/rol` | `rolInfra` está importado como comentario (`//rolInfra "..."`) en `dependencies.go:59`; `NewRolRoutes(...).Register()` nunca se llama. |
| `src/notificacion` | Import en blanco `_ "…/src/notificacion/infrastructure"` (solo para compilar); `NewNotificacionRouter(engine)` y `NewPushNotificationRouter(engine)` no se invocan en ningún punto del repo. |

---

## 3. Autenticación y autorización

- JWT HS256, secreto en `JWT_SECRET` (env). Emitido en `core.GenerateToken(userID, roleID)` ([src/core/jwt.go](src/core/jwt.go)), expira a las 24h.
- Header esperado: `Authorization: Bearer <token>` (también acepta `bearer ` en minúscula).
- `core.JWTAuthMiddleware()` valida el token y setea `user_id` / `role_id` en el contexto Gin.
- `core.RequireRole(roles...int)` ([src/core/role_middleware.go](src/core/role_middleware.go)): **ADMIN (id=1) siempre pasa**, sin importar qué roles pida la ruta; el resto de roles solo si `role_id` está en la lista.
- Roles definidos ([src/core/roles.go](src/core/roles.go)):

| Constante | Valor |
|---|---|
| `ADMIN` | 1 |
| `CONDUCTOR` | 2 |
| `SUPERVISOR` | 3 |
| `COORDINADOR` | 4 |

- Login de empleados: `POST /api/empleados/login`. Login de ciudadanos: `POST /api/ciudadanos/login`.
- ⚠️ **Inconsistencia de seguridad detectada:** todos los controladores del paquete `src/Camion` (tipo-camion, ruta-camion, historial-asignación) llevan la anotación Swagger `@Security BearerAuth`, pero sus routers (`tipCamion_routes.go`, `rutaCamion_routes.go`, `historialAsignacionCamion_routes.go`) **no aplican `JWTAuthMiddleware()` ni `RequireRole()`** — a diferencia de todos los demás módulos de `src/Rutas`, que sí protegen sus grupos. Es decir, Swagger promete auth donde el servidor no la exige. Ver detalle en sección 5.4.

---

## 4. Referencia de endpoints (solo rutas realmente montadas)

Convención: **Auth** = requiere `Authorization: Bearer`; **Roles** = valores aceptados además de ADMIN (que siempre pasa); "—" = sin restricción de rol dentro de los autenticados; "público" = sin JWT.

### 4.1 Camión — tipos, historial de asignación y ruta-camión (`src/Camion`)

Prefijos `/api/tipo-camion`, `/api/ruta-camion`, `/api/historial-asignacion`. **Sin middleware de autenticación** (ver §3).

| Método | Ruta | Descripción | Body |
|---|---|---|---|
| POST | `/api/tipo-camion/` | Crear tipo de camión | `{nombre, descripcion}` |
| GET | `/api/tipo-camion/` | Listar tipos de camión | — |
| GET | `/api/tipo-camion/nombre/{nombre}` | Buscar tipo de camión por nombre | — |
| DELETE | `/api/tipo-camion/{id}` | Eliminar tipo de camión | — |
| POST | `/api/ruta-camion/` | Crear asignación ruta-camión | `{ruta_id, camion_id, fecha}` |
| GET | `/api/ruta-camion/` | Listar todas | — |
| GET | `/api/ruta-camion/{id}` | Obtener por id | — |
| GET | `/api/ruta-camion/camion/{camion_id}` | Buscar por camión | — |
| GET | `/api/ruta-camion/ruta/{ruta_id}` | Buscar por ruta | — |
| GET | `/api/ruta-camion/exists/{id}` | Verificar existencia | — |
| PUT | `/api/ruta-camion/{id}` | Actualizar | `{ruta_id, camion_id, fecha}` |
| DELETE | `/api/ruta-camion/{id}` | Eliminar (lógico) | — |
| POST | `/api/historial-asignacion/` | Crear asignación chofer↔camión | `{id_chofer, id_camion, fecha_asignacion}` |
| GET | `/api/historial-asignacion/` | Listar todo el historial | — |
| GET | `/api/historial-asignacion/{id}` | Obtener por id | — |
| PUT | `/api/historial-asignacion/{id}` | Actualizar registro | `{...}` |
| DELETE | `/api/historial-asignacion/{id}` | Eliminar registro | — |
| GET | `/api/historial-asignacion/camion/{camionId}` | Historial de un camión | — |
| GET | `/api/historial-asignacion/chofer/{choferId}` | Historial de un chofer | — |
| GET | `/api/historial-asignacion/activo/camion/{camionId}` | Asignación activa del camión | — |
| GET | `/api/historial-asignacion/activo/chofer/{choferId}` | Asignación activa del chofer | — |
| PUT | `/api/historial-asignacion/baja/{id}` | Dar de baja una asignación | — |
| PUT | `/api/historial-asignacion/cerrar/camion/{camionId}` | Cerrar asignación activa por camión | — |
| PUT | `/api/historial-asignacion/cerrar/chofer/{choferId}` | Cerrar asignación activa por chofer | — |

### 4.2 Camión — CRUD y estado (`src/Rutas`, prefijo `/api/camion`, `/api/estado-camion`)

Auth: JWT + roles `ADMIN, CONDUCTOR, SUPERVISOR, COORDINADOR` (en la práctica: cualquier empleado autenticado).

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/camion/` | Crear camión |
| GET | `/api/camion/` | Listar camiones |
| GET | `/api/camion/{id}` | Obtener por id |
| PUT | `/api/camion/{id}` | Actualizar |
| DELETE | `/api/camion/{id}` | Eliminar (baja lógica) |
| GET | `/api/camion/placa/{placa}` | Buscar por placa |
| GET | `/api/camion/modelo?modelo=` | Buscar por modelo (query param) |
| POST | `/api/estado-camion/` | Registrar estado de camión |
| GET | `/api/estado-camion/` | Listar estados |
| GET | `/api/estado-camion/camion/{id}` | Estado(s) de un camión |
| PUT | `/api/estado-camion/{id}` | Actualizar estado |
| DELETE | `/api/estado-camion/{id}` | Eliminar estado |

### 4.3 Rutas, puntos de recolección, relleno sanitario, registro de vaciado (`src/Rutas`)

Auth: JWT + roles `ADMIN, CONDUCTOR, SUPERVISOR, COORDINADOR`.

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/rutas/` | Crear ruta |
| GET | `/api/rutas/` | Listar rutas |
| GET | `/api/rutas/{id}` | Obtener por id |
| PUT | `/api/rutas/{id}` | Actualizar |
| DELETE | `/api/rutas/{id}` | Eliminar |
| GET | `/api/rutas/activas` | Listar rutas activas |
| POST | `/api/puntos-recoleccion/` | Crear punto de recolección |
| GET | `/api/puntos-recoleccion/` | Listar puntos |
| GET | `/api/puntos-recoleccion/{id}` | Obtener por id |
| GET | `/api/puntos-recoleccion/ruta/{rutaId}` | Puntos de una ruta |
| PUT | `/api/puntos-recoleccion/{id}` | Actualizar |
| DELETE | `/api/puntos-recoleccion/{id}` | Eliminar |
| POST | `/api/relleno-sanitario/` | Crear relleno sanitario |
| GET | `/api/relleno-sanitario/` | Listar |
| GET | `/api/relleno-sanitario/{id}` | Obtener por id |
| PUT | `/api/relleno-sanitario/{id}` | Actualizar |
| DELETE | `/api/relleno-sanitario/{id}` | Eliminar |
| GET | `/api/relleno-sanitario/buscar?nombre=` | Buscar por nombre |
| GET | `/api/relleno-sanitario/exists/{id}` | Verificar existencia |
| POST | `/api/registro-vaciado/` | Crear registro de vaciado |
| GET | `/api/registro-vaciado/` | Listar |
| GET | `/api/registro-vaciado/{id}` | Obtener por id |
| GET | `/api/registro-vaciado/relleno/{relleno_id}` | Vaciados de un relleno |
| GET | `/api/registro-vaciado/ruta-camion/{ruta_camion_id}` | Vaciados de una asignación ruta-camión |
| GET | `/api/registro-vaciado/exists/{id}` | Verificar existencia |
| DELETE | `/api/registro-vaciado/{id}` | Eliminar (no tiene `update`) |

### 4.4 Colonia (`src/colonia`, prefijo `/api/colonia`)

| Método | Ruta | Auth |
|---|---|---|
| GET | `/api/colonia` | Público |
| GET | `/api/colonia/{id}` | Público |
| POST | `/api/colonia` | JWT + `ADMIN` |
| PUT | `/api/colonia/{id}` | JWT + `ADMIN` |
| DELETE | `/api/colonia/{id}` | JWT + `ADMIN` |

### 4.5 Ciudadanos y domicilios (`src/Ciudadanos`)

| Método | Ruta | Auth |
|---|---|---|
| POST | `/api/ciudadanos` | Público (auto-registro) |
| POST | `/api/ciudadanos/login` | Público |
| GET | `/api/ciudadanos` | JWT + `ADMIN` |
| GET | `/api/ciudadanos/{id}` | JWT + `ADMIN` |
| PATCH | `/api/ciudadanos/{id}` | JWT + `ADMIN` |
| DELETE | `/api/ciudadanos/{id}` | JWT + `ADMIN` |
| POST | `/api/domicilios` | JWT (cualquier rol) |
| GET | `/api/domicilios` | JWT (cualquier rol) |
| GET | `/api/domicilios/{id}` | JWT (cualquier rol) |
| PUT | `/api/domicilios/{id}` | JWT (cualquier rol) |
| DELETE | `/api/domicilios/{id}` | JWT (cualquier rol) |

### 4.6 Empleados (`src/empleado`, prefijo `/api/empleados`)

| Método | Ruta | Auth |
|---|---|---|
| POST | `/api/empleados/login` | Público |
| POST | `/api/empleados/` | JWT + `ADMIN` |
| GET | `/api/empleados/` | JWT + `ADMIN` |
| GET | `/api/empleados/{id}` | JWT + `ADMIN` |
| PATCH | `/api/empleados/{id}` | JWT + `ADMIN` |
| DELETE | `/api/empleados/{id}` | JWT + `ADMIN` |

### 4.7 Alertas de usuario (`src/alerta_usuario`, prefijo `/api/alertas`)

Todo bajo JWT.

| Método | Ruta | Roles |
|---|---|---|
| POST | `/api/alertas` | `ADMIN, SUPERVISOR` |
| GET | `/api/alertas` | Cualquiera autenticado (lista las propias) |
| PUT | `/api/alertas/{id}/leida` | Cualquiera autenticado |

### 4.8 Fallas / Anomalías (`src/Fallas`, prefijo `/api/anomalias`)

Tabla unificada que reemplaza los antiguos conceptos `Anomalia`, `Incidencia`, `ReporteConductor`, `ReporteFallaCritica`, `SeguimientoFallaCritica` (ver §6.9). Auth: JWT + roles `ADMIN, SUPERVISOR, COORDINADOR`.

| Método | Ruta | Descripción |
|---|---|---|
| POST | `/api/anomalias/` | Crear anomalía |
| GET | `/api/anomalias/` | Listar todas |
| GET | `/api/anomalias/{id}` | Obtener por id |
| PUT | `/api/anomalias/{id}` | Actualizar |
| DELETE | `/api/anomalias/{id}` | Eliminar |
| GET | `/api/anomalias/punto/{puntoId}` | Por punto de recolección |
| GET | `/api/anomalias/chofer/{choferId}` | Por chofer/conductor |
| GET | `/api/anomalias/camion/{camionId}` | Por camión |
| GET | `/api/anomalias/ruta/{rutaId}` | Por ruta |
| GET | `/api/anomalias/referencia/{referenciaId}` | Seguimientos de una anomalía (auto-relación) |
| GET | `/api/anomalias/estado?estado=` | Filtrar por estado (`PENDIENTE`/`EN_PROCESO`/`RESUELTA`) |
| GET | `/api/anomalias/tipo?tipo_anomalia=` | Filtrar por tipo (ver enum §6.9) |
| GET | `/api/anomalias/por-fecha?fecha_inicio=&fecha_fin=` | Filtrar por rango de fechas |

### 4.9 Código muerto — no accesible en runtime (documentado por completitud)

Estas rutas **existen en el código y compilan**, pero como `dependencies.go` nunca invoca su registro, **no responden en el servidor real** (404, ya que Gin nunca las registra):

**`src/rol`** (prefijo `/api/roles` si se activara):
`POST /api/roles`, `GET /api/roles`, `PUT /api/roles/{id}`, `DELETE /api/roles/{id}` — todas `RequireRole(ADMIN)`.

**`src/notificacion`** (prefijo `/api/notificaciones` si se activara): ~28 endpoints CRUD + búsquedas + contadores + acciones (`/emergencia`, `/falla`, `/mantenimiento`, `/notificar`, `/enviar-multiples`, `/enviar-todos`, `/marcar-leida`, `/reactivar`, etc.)

**`src/notificacion` (push/realtime)** (prefijo `/api/notifications`, `/api/notifications/rules`, `/api/realtime/ws` si se activara): envío push vía FCM, reglas de notificación por estado, sesiones realtime de admin (WS upgrade token, heartbeat).

---

## 5. Auditoría de Swagger — ¿documenta TODAS las rutas?

**No.** Swagger está desincronizado del código en varias direcciones. Resumen cuantitativo:

- `docs/swagger.json` declara **130 rutas**.
- Rutas realmente montadas en el servidor: **≈93**.
- De esas 93, **3 no aparecen en Swagger** aunque sí tienen anotación `@Router` en el código (falta regenerar).
- **~40 rutas** en `docs/swagger.json` no corresponden a ningún código existente (paquetes/controladores que ya no existen).
- **4 rutas** (`/api/roles*`) están en Swagger y en el código, pero el módulo no está montado → Swagger promete un endpoint que da 404.
- **~28 rutas** de `notificacion` (CRUD) sí tienen `@Router` en el código, están en Swagger, pero el módulo tampoco está montado → mismo problema.
- **~15 rutas** de `notificacion` push/realtime ni siquiera tienen anotaciones Swagger (sus controladores no llevan `@Router`), y tampoco están montadas.

### 5.1 Rutas reales sin documentar en Swagger (regenerar `swag init`)

El código ya tiene la anotación correcta, pero el swagger generado (`docs/docs.go`, `swagger.json/yaml`) quedó desactualizado:

| Ruta real | Archivo con la anotación |
|---|---|
| `GET /api/anomalias/camion/{camionId}` | `src/Fallas/infrastructure/GetAnomaliasByCamionIDController.go` |
| `GET /api/anomalias/ruta/{rutaId}` | `src/Fallas/infrastructure/GetAnomaliasByRutaIDController.go` |
| `GET /api/anomalias/referencia/{referenciaId}` | `src/Fallas/infrastructure/GetAnomaliasByReferenciaIDController.go` |

**Acción:** correr `swag init` (o el comando que use el proyecto) para regenerar `docs/`.

### 5.2 Rutas "fantasma" en Swagger (ya no existen en el código)

`docs/swagger.json` sigue conteniendo rutas de un módulo previo de "Fallas y Mantenimiento" más granular, que fue reemplazado por la tabla unificada `anomalia` (ver commit "Fallas arregladas con todo tipo de anomalias"). No hay ni un solo archivo `.go` en `src/` que registre estos paths:

```
/api/alertas-mantenimiento/...            (9 rutas)
/api/incidencias/...                      (5 rutas)
/api/registros-mantenimiento/...          (6 rutas)
/api/reportes-conductor/...               (6 rutas)
/api/reportes-falla-critica/...           (5 rutas)
/api/reportes-mantenimiento-generado/...  (5 rutas)
/api/seguimientos-falla-critica/...       (4 rutas)
/api/tipos-mantenimiento/...              (2 rutas)
```

**Acción:** regenerar Swagger desde cero (borrar `docs/` y volver a correr `swag init`) para eliminar este ruido; de lo contrario cualquier consumidor de la doc probará endpoints que no existen.

### 5.3 Módulos documentados pero no montados (404 real)

| Módulo | Rutas en swagger.json | Estado real |
|---|---|---|
| `rol` | `/api/roles`, `/api/roles/{id}` (4 combinaciones método+ruta) | Router nunca registrado → 404 |
| `notificacion` (CRUD) | `/api/notificaciones/...` (~28 combinaciones) | Router nunca registrado → 404 |

**Acción:** decidir si estos módulos deben activarse (agregar las llamadas correspondientes en `dependencies.go`) o eliminarse del repo si ya no se usan. Documentarlos como si funcionaran es engañoso para cualquier equipo frontend/móvil que consuma el swagger.

### 5.4 Anotaciones Swagger con datos incorrectos (bug de copy-paste)

Los `@Success` / `@Param body` de varios controladores apuntan a un DTO equivocado, aparentemente copiado de otro controlador como plantilla y no actualizado. Esto no afecta la ruta en sí (`@Router` es correcto) pero sí el **schema** que Swagger publica para esa respuesta/petición:

- **Todo el paquete `src/Camion` (24 controladores: tipo-camion, ruta-camion, historial-asignación)** documenta sus respuestas con `entities.HistorialAsignacionCamionResponse` / `...ListResponse` / `...MessageResponse` y su `@Param body` con `entities.CreateHistorialAsignacionCamionRequest`, **incluso en los endpoints de `tipo-camion` y `ruta-camion`**, que son recursos distintos y deberían tener su propio DTO (`TipoCamionResponse`, `RutaCamionResponse`, etc., que no existen en `src/Camion/domain/entities/historial_asignacion_swagger.go`).
- **Varios controladores de `src/Rutas`** (p. ej. `getByModelo_controller.go`, `createRuta_controller.go`, `getAllCamion_controller.go`, `getAllRuta_controller.go`, `getRutaById_controller.go`, `updateRuta_controller.go`, `createCamion_controller.go`) documentan sus respuestas/bodies como `entities.EstadoCamionResponse` / `entities.EstadoCamionListResponse` / `entities.CreateEstadoCamionRequest`, cuando en realidad manejan `Camion` o `Ruta`.

**Acción:** revisar cada `*_swagger.go` de `src/Camion/domain/entities` y `src/Rutas/domain/entities` y corregir los `@Param`/`@Success` de los controladores para que apunten al DTO real del recurso.

### 5.5 `@Security BearerAuth` sin protección real

Los 24 controladores de `src/Camion` (tipo-camion, ruta-camion, historial-asignación) llevan `// @Security BearerAuth`, pero sus tres archivos de rutas (`tipCamion_routes.go`, `rutaCamion_routes.go`, `historialAsignacionCamion_routes.go`) **no aplican ningún middleware**. Es decir, Swagger le dice al consumidor que necesita un token, pero el servidor acepta la petición sin él. Contrastar con `estadoCamion_routes.go`, `ruta_routes.go`, etc. (mismo paquete `src/Rutas`), que sí llaman `.Use(core.JWTAuthMiddleware(), core.RequireRole(...))`.

**Acción:** decidir si estas rutas deben protegerse (agregar el middleware, lo más probable dado que el resto del dominio Camión/Ruta sí lo exige) o si la anotación Swagger debe quitarse porque son intencionalmente públicas.

### 5.6 Push notifications / realtime sin ninguna anotación Swagger

`src/notificacion/infrastructure/push_notification_routes.go` registra 11 rutas (`/api/notifications/citizens/send`, `/api/notifications/events/...`, `/api/notifications/rules...`, `/api/realtime/ws/...`) cuyos controladores (`send_citizen_notification_controller.go`, `truck_state_event_controller.go`, `notification_rules_controller.go`, `admin_realtime_session_controller.go`) no tienen ninguna anotación `@Router`/`@Summary`. Por eso ni siquiera aparecen como "fantasma" en swagger.json — simplemente son invisibles para la herramienta. (Recuerda además que, actualmente, tampoco están montadas en el servidor — sección 4.9.)

---

## 6. Diccionario de datos

Fuente: [db_script.sql](db_script.sql) + [db_constraints.sql](db_constraints.sql) + [db_indexes.sql](db_indexes.sql), cruzado con las entidades Go de cada módulo. Motor: PostgreSQL. Todas las tablas con `updated_at` reciben un trigger genérico `update_timestamp` que refresca esa columna en cada `UPDATE` (ver `db_script.sql`, función `update_updated_at_column`).

### 6.1 `rol`

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SMALLINT | NO (PK) | Asignado manualmente, no `SERIAL` (ver `core/roles.go`: 1=ADMIN, 2=CONDUCTOR, 3=SUPERVISOR, 4=COORDINADOR) |
| `nombre` | VARCHAR(50) | NO | `UNIQUE` (`uq_rol_nombre`) |
| `active` | BOOLEAN | — | Default `TRUE` |

Entidad Go (`src/rol/domain/entities/rol.go`): `ID int`, `Nombre string`, `Eliminado bool` — nótese que el campo Go se llama `Eliminado` pero la columna real es `active` (semántica invertida: `Eliminado` no es literalmente `NOT active`, revisar mapeo en el repositorio si se reactiva el módulo).

### 6.2 `empleado`

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SERIAL | NO (PK) | |
| `nombre` | VARCHAR(100) | NO | |
| `apellidos` | VARCHAR(100) | NO | |
| `mail` | VARCHAR(100) | NO | `UNIQUE` (`uq_mail_empleado`), índice `idx_empleado_mail` |
| `password` | VARCHAR(100) | NO | Hash (bcrypt, ver `src/security/password`) |
| `username` | VARCHAR(100) | NO | `UNIQUE` (`uq_username_empleado`), índice `idx_empleado_username` |
| `desactivado` | BOOLEAN | — | Default `FALSE` |
| `rol_id` | SMALLINT | NO | FK → `rol(id)` (`fk_rol`) |
| `created_at` | TIMESTAMP | — | Default `CURRENT_TIMESTAMP` |
| `updated_at` | TIMESTAMP | — | Trigger automático |
| `deleted_at` | TIMESTAMP | sí | Baja lógica |

Entidad Go: `Empleado{ID, Nombre, Apellidos, Mail, Username, Password (json:"-"), Desactivado, RolID, CreatedAt, UpdatedAt, DeletedAt *time.Time}`.

### 6.3 `licencia` — ⚠️ sin módulo Go asociado

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SERIAL | NO (PK) | |
| `licencia` | VARCHAR(100) | NO | `UNIQUE` (`uq_licencia`), índice `idx_licencia` |
| `tipo_licencia` | SMALLINT | NO | |
| `fecha_vencimiento` | DATE | NO | `CHECK (fecha_vencimiento > CURRENT_DATE)` |
| `id_empleado` | INTEGER | NO | FK → `empleado(id)` |
| `created_at` | TIMESTAMP | NO | |
| `updated_at` | TIMESTAMP | — | |

No existe ningún paquete `src/licencia` ni referencias a esta tabla fuera del propio SQL: la tabla existe en el esquema pero **no tiene API**.

### 6.4 `camion`

| Columna DB | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SERIAL | NO (PK) | |
| `placa` | VARCHAR(20) | NO | `UNIQUE` (`uq_placa_camion`), índice `idx_placa_camion` |
| `modelo` | VARCHAR(50) | NO | |
| `rentado` | BOOLEAN | — | Default `FALSE` |
| `estado` | VARCHAR(20) | NO | Valores usados por la app: `OPERATIVO`, `MANTENIMIENTO`, `FUERA_SERVICIO`, `BAJA` (no hay `CHECK` en BD, solo convención en Go) |
| `tipo_id` | INTEGER | NO | FK → `tipo_camion(id)` (`fk_tipo_camion`) |
| `created_at` | TIMESTAMP | NO | |
| `updated_at` | TIMESTAMP | — | |
| `deleted_at` | TIMESTAMP | sí | Baja lógica |

Entidad Go (`Camion`, `src/Rutas/domain/entities/Camion.go`) expone campos **calculados, no persistidos**: `DisponibilidadID`, `NombreDisponibilidad`, `ColorDisponibilidad`. Se derivan en `PostgresCamion.mapEstadoToDisponibilidad()` a partir de la columna `estado` (mapeo fijo en código: 1/OPERATIVO/green, 2/MANTENIMIENTO/orange, 3/FUERA_SERVICIO/red, 4/BAJA/grey). No confundir con una tabla `disponibilidad` — no existe tal tabla.

### 6.5 `tipo_camion`

| Columna | Tipo | Nulo |
|---|---|---|
| `id` | SMALLINT | NO (PK, asignado manualmente) |
| `nombre` | VARCHAR(50) | NO |
| `descripcion` | VARCHAR(255) | sí |

### 6.6 `historial_asignacion_camion`

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id_historial` | SERIAL | NO (PK) | |
| `id_chofer` | INTEGER | NO | FK → `empleado(id)` (`fk_chofer_historial`) |
| `id_camion` | INTEGER | NO | FK → `camion(id)` (`fk_camion_historial`) |
| `fecha_asignacion` | DATE | NO | `CHECK (fecha_asignacion <= CURRENT_DATE)` |
| `fecha_baja` | DATE | sí | `CHECK (fecha_baja IS NULL OR fecha_baja >= fecha_asignacion)` |
| `eliminado` | BOOLEAN | — | Default `FALSE` |
| `created_at` | TIMESTAMP | NO | |
| `updated_at` | TIMESTAMP | — | |
| `deleted_at` | TIMESTAMP | sí | |

`UNIQUE (id_chofer, id_camion, fecha_asignacion)` (`uq_chofer_camion_historial`). Entidad Go usa punteros (`*int`, `*time.Time`) para `IDChofer`/`IDCamion`/`FechaAsignacion`/`FechaBaja` pese a que en BD `id_chofer`/`id_camion`/`fecha_asignacion` son `NOT NULL`.

### 6.7 `ruta_camion`

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `ruta_camion_id` | SERIAL | NO (PK) | |
| `ruta_id` | INTEGER | NO | FK → `ruta(id)` (`fk_ruta_asignada`) |
| `camion_id` | INTEGER | NO | FK → `camion(id)` (`fk_camion_asignado_ruta`) |
| `fecha` | DATE | NO | `CHECK (fecha <= CURRENT_DATE)` |
| `eliminado` | BOOLEAN | — | Default `FALSE` |
| `created_at` | TIMESTAMP | NO | |
| `updated_at` | TIMESTAMP | — | |
| `deleted_at` | TIMESTAMP | sí | |

### 6.8 `ruta`, `punto_recoleccion`, `relleno_sanitario`, `estado_camion`, `registro_vaciado`

**`ruta`**

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SERIAL | NO (PK) | |
| `nombre` | VARCHAR(100) | NO | `UNIQUE` (`uq_nombre_ruta`), `CHECK (nombre <> '')` |
| `descripcion` | VARCHAR(255) | NO | |
| `colonia_id` | INTEGER | NO | FK → `colonia(colonia_id)` (`fk_colonia_ruta`) — ⚠️ la entidad Go `Ruta` (`src/Rutas/domain/entities/Ruta.go`) **no expone `colonia_id`**, aunque la columna es `NOT NULL` en BD |
| `json_ruta` | JSON | NO | Traza geográfica de la ruta |
| `created_at` | TIMESTAMP | NO | |
| `updated_at` | TIMESTAMP | — | |
| `deleted_at` | TIMESTAMP | sí | Índice parcial `idx_registro_activo_asignacion_ruta` |

**`punto_recoleccion`**

| Columna DB | Tipo | Nulo |
|---|---|---|
| `id` | SERIAL | NO (PK) |
| `ruta_id` | INTEGER | NO, FK → `ruta(id)` |
| `direccion` | VARCHAR(255) | NO |
| `created_at` / `updated_at` / `deleted_at` | TIMESTAMP | — |

⚠️ La entidad Go `PuntoRecoleccion` expone `CP`, `Lat`, `Lon`, `Eliminado` — campos que **no existen** como columnas en `db_script.sql` (que solo define `direccion`). Indica que el repositorio construye/mapea estos campos desde otro origen o que el script SQL versionado en el repo no está sincronizado con el esquema real en producción. **Revisar antes de confiar en el diccionario para este recurso.**

**`relleno_sanitario`**

| Columna | Tipo | Nulo |
|---|---|---|
| `relleno_id` | SERIAL | NO (PK) |
| `nombre` | VARCHAR(100) | NO |
| `direccion` | VARCHAR(255) | NO |
| `es_rentado` | BOOLEAN | Default `FALSE` |
| `capacidad_toneladas` | DOUBLE PRECISION | NO |
| `eliminado` | BOOLEAN | Default `FALSE` |
| `created_at` / `updated_at` / `deleted_at` | TIMESTAMP | — |

**`estado_camion`**

| Columna | Tipo | Nulo |
|---|---|---|
| `estado_id` | SERIAL | NO (PK) |
| `camion_id` | INTEGER | NO, FK → `camion(id)` |
| `estado` | VARCHAR(50) | NO |
| `observaciones` | TEXT | sí |
| `timestamp` | TIMESTAMP | NO |
| `created_at` | TIMESTAMP | NO |

**`registro_vaciado`**

| Columna | Tipo | Nulo |
|---|---|---|
| `vaciado_id` | SERIAL | NO (PK) |
| `relleno_id` | INTEGER | NO, FK → `relleno_sanitario(relleno_id)` |
| `ruta_camion_id` | INTEGER | NO, FK → `ruta_camion(ruta_camion_id)` |
| `hora` | TIMESTAMP | NO |
| `created_at` | TIMESTAMP | NO |

### 6.9 `anomalia` (dominio Fallas)

Tabla **unificada**: sustituye lo que antes eran 5 tablas/conceptos independientes (`Anomalia`, `Incidencia`, `ReporteConductor`, `ReporteFallaCritica`, `SeguimientoFallaCritica`). El campo `tipo_anomalia` indica cuál de esos 5 conceptos representa cada fila.

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `anomalia_id` | SERIAL | NO (PK) | |
| `tipo_anomalia` | ENUM `tipo_anomalia_enum` | NO | Valores: `ANOMALIA`, `INCIDENCIA`, `REPORTE_CONDUCTOR`, `REPORTE_FALLA_CRITICA`, `SEGUIMIENTO_FALLA_CRITICA` |
| `punto_id` | INTEGER | sí | FK → `punto_recoleccion(id)` |
| `conductor_id` | INTEGER | sí | FK → `empleado(id)` |
| `camion_id` | INTEGER | sí | FK → `camion(id)` |
| `ruta_id` | INTEGER | sí | FK → `ruta(id)` |
| `anomalia_referencia_id` | INTEGER | sí | FK → `anomalia(anomalia_id)` — **auto-relación**; usada por `SEGUIMIENTO_FALLA_CRITICA` para apuntar al `REPORTE_FALLA_CRITICA` original |
| `descripcion` | TEXT | NO | `CHECK (descripcion <> '')` |
| `json_ruta` | TEXT | sí | |
| `estado` | VARCHAR(30) | sí | `CHECK (estado IN ('PENDIENTE','EN_PROCESO','RESUELTA') OR estado IS NULL)` |
| `eliminado` | BOOLEAN | — | Default `FALSE`, índice parcial |
| `fecha_reporte` | TIMESTAMP | NO | Índice |
| `fecha_resolucion` | TIMESTAMP | sí | `CHECK (fecha_resolucion IS NULL OR fecha_resolucion >= fecha_reporte)` |
| `created_at` / `updated_at` | TIMESTAMP | — | |

Índices: `tipo_anomalia`, `estado`, `punto_id`, `conductor_id`, `camion_id`, `ruta_id`, `anomalia_referencia_id`, `fecha_reporte`, parcial sobre `eliminado = false`.

### 6.10 `colonia`

| Columna | Tipo | Nulo |
|---|---|---|
| `colonia_id` | SERIAL | NO (PK) |
| `nombre` | VARCHAR(100) | NO, índice `idx_nombre_colonia` |
| `zona` | VARCHAR(50) | NO |
| `created_at` | TIMESTAMP | NO |

### 6.11 `ciudadano` / `domicilio`

**`ciudadano`**

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SERIAL | NO (PK) | |
| `email` | VARCHAR(100) | NO | `UNIQUE`, `CHECK` formato email (`chk_email_ciudadano`) |
| `alias` | VARCHAR(100) | NO | `UNIQUE` |
| `password` | VARCHAR(100) | NO | Hash; en Go `json:"-"` (nunca se serializa) |
| `created_at` / `updated_at` | TIMESTAMP | — | |

**`domicilio`**

| Columna | Tipo | Nulo | Notas |
|---|---|---|---|
| `id` | SERIAL | NO (PK) | |
| `alias` | VARCHAR(100) | NO | `CHECK (alias <> '')` |
| `calle` | VARCHAR(100) | NO | |
| `numero` | VARCHAR(20) | NO | |
| `referencia` | VARCHAR(255) | sí | |
| `ciudadano_id` | INTEGER | NO | FK → `ciudadano(id)` |
| `colonia_id` | INTEGER | NO | FK → `colonia(colonia_id)` |
| `deleted_at` | TIMESTAMP | sí | |
| `created_at` | TIMESTAMP | NO | |

### 6.12 `aviso` — ⚠️ sin módulo Go asociado

| Columna | Tipo | Nulo |
|---|---|---|
| `id` | SERIAL | NO (PK) |
| `enviado_por` | INTEGER | NO |
| `tipo_aviso` | VARCHAR(50) | NO |
| `descripcion` | VARCHAR(255) | NO |
| `entidad_involucrada` | VARCHAR(100) | NO |
| `estado` | SMALLINT | NO |
| `created_at` / `updated_at` / `deleted_at` | TIMESTAMP | — |

Al igual que `licencia`, esta tabla se crea en el esquema pero **ningún paquete Go la referencia**.

### 6.13 `tipo_mantenimiento`, `alerta_mantenimiento`, `registro_mantenimiento` — legado

Estas tres tablas se crean en `db_script.sql`, pero **no tienen módulo Go activo**. Su único rastro en el código es `src/core/db.go` (`AutoMigrateDatabase`), una rutina de desarrollo que — solo si `ENVIRONMENT=development` o `DEBUG=true` — verifica si el esquema está desactualizado y, de ser así, **las hace `DROP TABLE ... CASCADE`** junto con otras antes de re-ejecutar `db_script.sql`/`db_constraints.sql`/`db_indexes.sql`. Son remanentes del antiguo módulo "Fallas y Mantenimiento" fragmentado, consolidado ahora en `anomalia` (§6.9).

| Tabla | Columnas clave |
|---|---|
| `tipo_mantenimiento` | `tipo_mantenimiento_id` PK, `nombre`, `categoria`, `eliminado` |
| `alerta_mantenimiento` | `alerta_id` PK, `camion_id` FK→camion, `tipo_mantenimiento_id` FK, `descripcion`, `observaciones`, `atendido` |
| `registro_mantenimiento` | `registro_id` PK, `alerta_id` FK→alerta_mantenimiento, `camion_id` FK, `coordinador_id`, `mecanico_responsable`, `fecha_realizada`, `kilometraje_mantenimiento` (`CHECK >= 0`), `observaciones` |

### 6.14 `notificacion` — ⚠️ tabla no está en el esquema versionado

El repositorio `src/notificacion/infrastructure/postgres_notificacion_repository.go` hace `INSERT/SELECT/DELETE FROM notificacion` y también consulta una tabla `usuario` (`SELECT user_id FROM usuario WHERE eliminado = false OR eliminado IS NULL`). **Ninguna de las dos tablas (`notificacion`, `usuario`) se crea en `db_script.sql`, `db_constraints.sql` ni `db_indexes.sql`.** Si el módulo llegara a activarse tal cual, fallaría contra una base de datos provisionada solo con estos scripts. Columnas inferidas de la entidad Go (`src/notificacion/domain/entities/notificacion.go`):

| Columna (inferida) | Tipo Go | Notas |
|---|---|---|
| `notificacion_id` | int32 (PK) | |
| `usuario_id` | *int32 | `NULL` = notificación global |
| `tipo` | string | `"falla"`, `"mantenimiento"`, `"emergencia"`, `"ruta"`, etc. |
| `titulo` | string | |
| `mensaje` | string | |
| `activa` | bool | Se usa como proxy de "no leída" (`MarcarComoLeida()` pone `activa=false`) |
| `id_camion_relacionado` | *int32 | |
| `id_falla_relacionado` | *int32 | |
| `id_mantenimiento_relacionado` | *int32 | |
| `creado_por` | *int32 | |
| `created_at` | time.Time | |

### 6.15 `alerta_usuario` — ⚠️ tabla no está en el esquema versionado

`src/alerta_usuario/infrastructure/postgres/alerta_repository.go` opera sobre una tabla `alerta_usuario` que tampoco existe en ninguno de los tres scripts SQL del repo (no confundir con `alerta_mantenimiento`, que sí existe pero es de otro dominio). Como este módulo **sí está montado** (`/api/alertas`), cualquier entorno provisionado únicamente con `db_script.sql` + `db_constraints.sql` + `db_indexes.sql` fallará en tiempo de ejecución al usarlo. Columnas según la entidad Go (`src/alerta_usuario/domain/entity.go`):

| Columna (inferida) | Tipo Go | Notas |
|---|---|---|
| `alerta_id` | int (PK) | |
| `usuario_id` | int | Receptor de la alerta |
| `titulo` | string | |
| `mensaje` | string | |
| `leida` | bool | |
| `creado_por` | int | Supervisor/admin que la generó |
| `created_at` | time.Time | |

**Acción recomendada:** añadir el `CREATE TABLE` de `alerta_usuario` (y, si se reactiva el módulo, `notificacion`/`usuario`) a `db_script.sql`, o documentar en qué migración externa se crean.

---

## 7. Hallazgos y recomendaciones (priorizado)

1. **Regenerar Swagger** (`swag init`) — corrige automáticamente §5.1 y reduce el ruido de §5.2. Es la acción de mayor impacto/menor esfuerzo.
2. **Decidir el destino de `rol` y `notificacion`**: si son necesarios, cablearlos en `dependencies.go`; si no, quitarlos del repo o marcarlos claramente como WIP para que Swagger no prometa endpoints inexistentes (§4.9, §5.3).
3. **Proteger o des-anotar** las rutas de `src/Camion` (tipo-camion, ruta-camion, historial-asignación): hoy son públicas pese a decir `@Security BearerAuth` (§3, §5.5).
4. **Corregir los DTOs de Swagger** mal referenciados en `src/Camion` y en varios controladores de `src/Rutas` (§5.4) — son correctos en `Fallas`, `colonia`, `alerta_usuario`, `Ciudadanos`, `empleado`, así que hay plantilla de referencia dentro del propio repo.
5. **Sincronizar el esquema SQL versionado**: agregar `alerta_usuario` (usado por un módulo activo) y, si se reactiva `notificacion`, también `notificacion`/`usuario`. Revisar por qué `punto_recoleccion` en Go trae `cp/lat/lon` que no están en `db_script.sql` (§6.8).
6. **Tablas huérfanas** (`licencia`, `aviso`, `tipo_mantenimiento`, `alerta_mantenimiento`, `registro_mantenimiento`): decidir si se eliminan del esquema o se retoma su API.

---

## 8. Otros documentos del repo relacionados

- [README_FALLAS_API.md](README_FALLAS_API.md) — guía de uso del módulo Fallas/Anomalías.
- [README_FRONTEND_API.md](README_FRONTEND_API.md) / [README_MOBILE_API.md](README_MOBILE_API.md) — guías de consumo por plataforma.
- [CHANGELOG.md](CHANGELOG.md) — historial de cambios.

Este documento (`DOCUMENTACION_API.md`) es el único que cruza **código real montado** vs **Swagger** vs **esquema de base de datos**; los demás README describen el uso desde la perspectiva de un consumidor y no deben tomarse como fuente de verdad sobre qué está realmente expuesto.
