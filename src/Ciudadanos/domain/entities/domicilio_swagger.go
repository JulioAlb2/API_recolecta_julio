package entities

// CreateDomicilioRequest cuerpo para POST /api/domicilios.
// latitud y longitud son opcionales (WGS84); deben enviarse juntas si se usan.
type CreateDomicilioRequest struct {
	CiudadanoID int      `json:"ciudadano_id" example:"42"`
	ColoniaID   int      `json:"colonia_id" binding:"required" example:"1"`
	Alias       string   `json:"alias" binding:"required" example:"Casa"`
	Calle       string   `json:"calle" binding:"required" example:"Calle Olmo"`
	Numero      string   `json:"numero" binding:"required" example:"123"`
	Referencia  *string  `json:"referencia,omitempty" example:"Frente al parque"`
	Latitud     *float64 `json:"latitud,omitempty" example:"16.6278"`
	Longitud    *float64 `json:"longitud,omitempty" example:"-93.1045"`
}

// UpdateDomicilioRequest cuerpo para PUT /api/domicilios/{id}. Todos los campos son opcionales.
type UpdateDomicilioRequest struct {
	ColoniaID  *int     `json:"colonia_id,omitempty" example:"1"`
	Alias      *string  `json:"alias,omitempty" example:"Casa principal"`
	Calle      *string  `json:"calle,omitempty" example:"Calle Olmo"`
	Numero     *string  `json:"numero,omitempty" example:"125"`
	Referencia *string  `json:"referencia,omitempty" example:"Portón azul"`
	Latitud    *float64 `json:"latitud,omitempty" example:"16.6280"`
	Longitud   *float64 `json:"longitud,omitempty" example:"-93.1048"`
}

// DomicilioResponse respuesta de GET /api/domicilios/{id}.
type DomicilioResponse struct {
	Success bool      `json:"success" example:"true"`
	Message string    `json:"message" example:"domicilio obtenido correctamente"`
	Data    Domicilio `json:"data"`
	Code    int       `json:"code" example:"200"`
}

// DomicilioIDResponse respuesta de POST /api/domicilios.
type DomicilioIDResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"domicilio creado correctamente"`
	ID      int    `json:"id" example:"1"`
	Code    int    `json:"code" example:"201"`
}

// DomicilioDetailResponse respuesta detallada de domicilio.
type DomicilioDetailResponse struct {
	Success bool      `json:"success" example:"true"`
	Data    Domicilio `json:"data"`
}

// DomicilioListResponse respuesta de GET /api/domicilios.
type DomicilioListResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"domicilios listados correctamente"`
	Data    []Domicilio `json:"data"`
	Code    int         `json:"code" example:"200"`
}

// DomicilioMessageResponse respuesta sin cuerpo de domicilio (PUT/DELETE).
type DomicilioMessageResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"domicilio actualizado correctamente"`
	Code    int    `json:"code,omitempty" example:"200"`
}
