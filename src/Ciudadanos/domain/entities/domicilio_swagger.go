package entities

// swagger:model CreateDomicilioRequest
type CreateDomicilioRequest struct {
	CiudadanoID int     `json:"ciudadano_id"`
	ColoniaID   int     `json:"colonia_id" binding:"required"`
	Alias       string  `json:"alias" binding:"required"`
	Calle       string  `json:"calle" binding:"required"`
	Numero      string  `json:"numero" binding:"required"`
	Referencia  *string `json:"referencia,omitempty"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
}

// swagger:model UpdateDomicilioRequest
type UpdateDomicilioRequest struct {
	ColoniaID  *int    `json:"colonia_id,omitempty"`
	Alias      *string `json:"alias,omitempty"`
	Calle      *string `json:"calle,omitempty"`
	Numero     *string `json:"numero,omitempty"`
	Referencia *string `json:"referencia,omitempty"`
}

// swagger:model DomicilioResponse
type DomicilioResponse struct {
	Success bool      `json:"success"`
	Message string    `json:"message"`
	Data    Domicilio `json:"data"`
	Code    int       `json:"code"`
}

// swagger:model DomicilioIDResponse
type DomicilioIDResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	ID      int    `json:"id"`
	Code    int    `json:"code"`
}

// swagger:model DomicilioDetailResponse
type DomicilioDetailResponse struct {
	Success bool      `json:"success"`
	Data    Domicilio `json:"data"`
}

// swagger:model DomicilioListResponse
type DomicilioListResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    []Domicilio `json:"data"`
	Code    int         `json:"code"`
}

// swagger:model DomicilioMessageResponse
type DomicilioMessageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
}
