package entities


// swagger:model CreatePuntoRecoleccionRequest
type CreatePuntoRecoleccionRequest struct {
	PuntoID int32   `json:"punto_id"`
	RutaID  int32   `json:"ruta_id"`
	CP      string  `json:"cp"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// swagger:model UpdatePuntoRecoleccionRequest
type UpdatePuntoRecoleccionRequest struct {
	PuntoID int32   `json:"punto_id"`
	RutaID  int32   `json:"ruta_id"`
	CP      string  `json:"cp"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
}

// swagger:model PuntoRecoleccionResponse
type PuntoRecoleccionResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    PuntoRecoleccion `json:"data"`
	Code    int          `json:"code"`
}

// swagger:model PuntoRecoleccionDetailResponse
type PuntoRecoleccionDetailResponse struct {
	Success bool         `json:"success"`
	Data    PuntoRecoleccion `json:"data"`
}

// swagger:model PuntoRecoleccionListResponse
type PuntoRecoleccionListResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    []PuntoRecoleccion `json:"data"`
	Code    int            `json:"code"`
}

// swagger:model PuntoRecoleccionListSimpleResponse
type PuntoRecoleccionListSimpleResponse struct {
	Success bool           `json:"success"`
	Data    []PuntoRecoleccion `json:"data"`
}

// swagger:model PuntoRecoleccionMessageResponse
type PuntoRecoleccionMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
