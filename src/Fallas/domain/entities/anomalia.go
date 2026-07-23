// anomalia.go
package entities

import (
	"time"
)

// Anomalia es la entidad unificada que reemplaza a Anomalia, Incidencia,
// ReporteConductor, ReporteFallaCritica y SeguimientoFallaCritica.
// El campo TipoAnomalia (ver tipo_anomalia.go) indica cuál de esos
// conceptos representa cada registro.
//
// SeguimientoFallaCritica es un caso especial: es un seguimiento/comentario
// sobre un registro REPORTE_FALLA_CRITICA ya existente, por lo que usa
// AnomaliaReferenciaID para apuntar al AnomaliaID del reporte al que da
// seguimiento (auto-relación dentro de la misma tabla).
type Anomalia struct {
	AnomaliaID           int32        `json:"anomalia_id" gorm:"column:anomalia_id;primaryKey;autoIncrement"`
	TipoAnomalia         TipoAnomalia `json:"tipo_anomalia" gorm:"column:tipo_anomalia;not null"`
	PuntoID              *int32       `json:"punto_id" gorm:"column:punto_id"`
	ConductorID          *int32       `json:"conductor_id" gorm:"column:conductor_id"`
	CamionID             *int32       `json:"camion_id" gorm:"column:camion_id"`
	RutaID               *int32       `json:"ruta_id" gorm:"column:ruta_id"`
	AnomaliaReferenciaID *int32       `json:"anomalia_referencia_id" gorm:"column:anomalia_referencia_id"`
	Descripcion          string       `json:"descripcion" gorm:"column:descripcion;type:text;not null"`
	JsonRuta             string       `json:"json_ruta" gorm:"column:json_ruta;type:text"`
	Estado               string       `json:"estado" gorm:"column:estado;type:varchar(30)"`
	Eliminado            bool         `json:"eliminado" gorm:"column:eliminado;default:false"`
	FechaReporte         time.Time    `json:"fecha_reporte" gorm:"column:fecha_reporte;not null"`
	FechaResolucion      *time.Time   `json:"fecha_resolucion" gorm:"column:fecha_resolucion"`
	CreatedAt            time.Time    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt            time.Time    `json:"updated_at" gorm:"column:updated_at"`

	// Pipeline modelo_reportes -> clasificador_reportes. Se llenan por
	// ActualizarPipeline en background, no al crear la anomalia (por eso
	// EstadoPipeline arranca en "pendiente" via DEFAULT de la columna).
	EstadoPipeline        string  `json:"estado_pipeline" gorm:"column:estado_pipeline"`
	NivelRiesgo           *string `json:"nivel_riesgo" gorm:"column:nivel_riesgo"`
	InferenciaID          *int32  `json:"inferencia_id" gorm:"column:inferencia_id"`
	CategoriaClasificada  *string `json:"categoria_clasificada" gorm:"column:categoria_clasificada"`
	SubtipoClasificado    *string `json:"subtipo_clasificado" gorm:"column:subtipo_clasificado"`
	AccionSugerida        *string `json:"accion_sugerida" gorm:"column:accion_sugerida"`
	PipelineError         *string `json:"pipeline_error" gorm:"column:pipeline_error"`
}

// Setters
func (a *Anomalia) SetAnomaliaID(anomaliaID int32) {
	a.AnomaliaID = anomaliaID
}

func (a *Anomalia) SetTipoAnomalia(tipoAnomalia TipoAnomalia) {
	a.TipoAnomalia = tipoAnomalia
}

func (a *Anomalia) SetPuntoID(puntoID *int32) {
	a.PuntoID = puntoID
}

func (a *Anomalia) SetConductorID(conductorID *int32) {
	a.ConductorID = conductorID
}

func (a *Anomalia) SetCamionID(camionID *int32) {
	a.CamionID = camionID
}

func (a *Anomalia) SetRutaID(rutaID *int32) {
	a.RutaID = rutaID
}

func (a *Anomalia) SetAnomaliaReferenciaID(referenciaID *int32) {
	a.AnomaliaReferenciaID = referenciaID
}

func (a *Anomalia) SetDescripcion(descripcion string) {
	a.Descripcion = descripcion
}

func (a *Anomalia) SetJsonRuta(jsonRuta string) {
	a.JsonRuta = jsonRuta
}

func (a *Anomalia) SetEstado(estado string) {
	a.Estado = estado
}

func (a *Anomalia) SetEliminado(eliminado bool) {
	a.Eliminado = eliminado
}

func (a *Anomalia) SetFechaReporte(fechaReporte time.Time) {
	a.FechaReporte = fechaReporte
}

func (a *Anomalia) SetFechaResolucion(fechaResolucion *time.Time) {
	a.FechaResolucion = fechaResolucion
}

func (a *Anomalia) SetUpdatedAt(updatedAt time.Time) {
	a.UpdatedAt = updatedAt
}

// Getters
func (a *Anomalia) GetAnomaliaID() int32 {
	return a.AnomaliaID
}

func (a *Anomalia) GetTipoAnomalia() TipoAnomalia {
	return a.TipoAnomalia
}

func (a *Anomalia) GetPuntoID() *int32 {
	return a.PuntoID
}

func (a *Anomalia) GetConductorID() *int32 {
	return a.ConductorID
}

func (a *Anomalia) GetCamionID() *int32 {
	return a.CamionID
}

func (a *Anomalia) GetRutaID() *int32 {
	return a.RutaID
}

func (a *Anomalia) GetAnomaliaReferenciaID() *int32 {
	return a.AnomaliaReferenciaID
}

func (a *Anomalia) GetDescripcion() string {
	return a.Descripcion
}

func (a *Anomalia) GetJsonRuta() string {
	return a.JsonRuta
}

func (a *Anomalia) GetEstado() string {
	return a.Estado
}

func (a *Anomalia) GetEliminado() bool {
	return a.Eliminado
}

func (a *Anomalia) GetFechaReporte() time.Time {
	return a.FechaReporte
}

func (a *Anomalia) GetFechaResolucion() *time.Time {
	return a.FechaResolucion
}

func (a *Anomalia) GetCreatedAt() time.Time {
	return a.CreatedAt
}

func (a *Anomalia) GetUpdatedAt() time.Time {
	return a.UpdatedAt
}

// NewAnomalia crea una anomalía lista para persistir, fijando timestamps.
func NewAnomalia(tipoAnomalia TipoAnomalia, puntoID, conductorID, camionID, rutaID, anomaliaReferenciaID *int32, descripcion, jsonRuta, estado string, fechaReporte time.Time) *Anomalia {
	if fechaReporte.IsZero() {
		fechaReporte = time.Now()
	}
	now := time.Now()
	return &Anomalia{
		TipoAnomalia:         tipoAnomalia,
		PuntoID:              puntoID,
		ConductorID:          conductorID,
		CamionID:             camionID,
		RutaID:               rutaID,
		AnomaliaReferenciaID: anomaliaReferenciaID,
		Descripcion:          descripcion,
		JsonRuta:             jsonRuta,
		Estado:               estado,
		Eliminado:            false,
		FechaReporte:         fechaReporte,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// NewAnomaliaCompleta reconstruye una anomalía completa (actualización o consulta).
func NewAnomaliaCompleta(anomaliaID int32, tipoAnomalia TipoAnomalia, puntoID, conductorID, camionID, rutaID, anomaliaReferenciaID *int32, descripcion, jsonRuta, estado string, eliminado bool, fechaReporte time.Time, fechaResolucion *time.Time, createdAt, updatedAt time.Time) *Anomalia {
	return &Anomalia{
		AnomaliaID:           anomaliaID,
		TipoAnomalia:         tipoAnomalia,
		PuntoID:              puntoID,
		ConductorID:          conductorID,
		CamionID:             camionID,
		RutaID:               rutaID,
		AnomaliaReferenciaID: anomaliaReferenciaID,
		Descripcion:          descripcion,
		JsonRuta:             jsonRuta,
		Estado:               estado,
		Eliminado:            eliminado,
		FechaReporte:         fechaReporte,
		FechaResolucion:      fechaResolucion,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}

// Método para marcar como resuelta
func (a *Anomalia) MarcarResuelta(fechaResolucion time.Time) {
	a.Estado = "RESUELTA"
	a.FechaResolucion = &fechaResolucion
	a.UpdatedAt = time.Now()
}

// Método para marcar como pendiente
func (a *Anomalia) MarcarPendiente() {
	a.Estado = "PENDIENTE"
	a.FechaResolucion = nil
	a.UpdatedAt = time.Now()
}

// Método para marcar como en proceso
func (a *Anomalia) MarcarEnProceso() {
	a.Estado = "EN_PROCESO"
	a.FechaResolucion = nil
	a.UpdatedAt = time.Now()
}

// Método para marcar como eliminada (borrado lógico)
func (a *Anomalia) MarcarEliminada() {
	a.Eliminado = true
	a.UpdatedAt = time.Now()
}

// Método para restaurar (deshacer borrado lógico)
func (a *Anomalia) Restaurar() {
	a.Eliminado = false
	a.UpdatedAt = time.Now()
}

// Método para verificar si está resuelta
func (a *Anomalia) EstaResuelta() bool {
	return a.Estado == "RESUELTA"
}

// Método para verificar si está pendiente
func (a *Anomalia) EstaPendiente() bool {
	return a.Estado == "PENDIENTE"
}

// Método para verificar si está en proceso
func (a *Anomalia) EstaEnProceso() bool {
	return a.Estado == "EN_PROCESO"
}

// Método para agregar más información a la descripción
func (a *Anomalia) AgregarDescripcion(descripcionAdicional string) {
	if a.Descripcion == "" {
		a.Descripcion = descripcionAdicional
	} else {
		a.Descripcion = a.Descripcion + "\n" + descripcionAdicional
	}
	a.UpdatedAt = time.Now()
}

// Método para verificar si tiene punto asociado
func (a *Anomalia) TienePunto() bool {
	return a.PuntoID != nil
}

// Método para verificar si tiene fecha de resolución
func (a *Anomalia) TieneFechaResolucion() bool {
	return a.FechaResolucion != nil
}

// Método para verificar si es un seguimiento sobre otra anomalía
func (a *Anomalia) EsSeguimiento() bool {
	return a.TipoAnomalia == TipoAnomaliaSeguimientoFallaCritica
}

// Método para obtener una descripción resumida (primeros 100 caracteres)
func (a *Anomalia) GetDescripcionResumida() string {
	if len(a.Descripcion) <= 100 {
		return a.Descripcion
	}
	return a.Descripcion[:100] + "..."
}

// Método para obtener fecha de reporte formateada
func (a *Anomalia) GetFechaReporteFormateada() string {
	return a.FechaReporte.Format("02/01/2006 15:04:05")
}

// Método para obtener fecha de resolución formateada (si existe)
func (a *Anomalia) GetFechaResolucionFormateada() string {
	if a.FechaResolucion != nil {
		return a.FechaResolucion.Format("02/01/2006 15:04:05")
	}
	return ""
}

// TableName especifica el nombre de la tabla para GORM
func (Anomalia) TableName() string {
	return "anomalia"
}
