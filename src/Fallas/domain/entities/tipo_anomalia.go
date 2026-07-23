// tipo_anomalia.go
package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// TipoAnomalia es un enum real (no un simple string) que identifica cuál de
// los antiguos conceptos (Anomalia, Incidencia, ReporteConductor,
// ReporteFallaCritica, SeguimientoFallaCritica) representa cada registro de
// la tabla unificada anomalia. Los únicos valores válidos son las
// constantes predefinidas abajo.
type TipoAnomalia int

const (
	TipoAnomaliaAnomalia TipoAnomalia = iota
	TipoAnomaliaIncidencia
	TipoAnomaliaReporteConductor
	TipoAnomaliaReporteFallaCritica
	TipoAnomaliaSeguimientoFallaCritica
)

// nombrePorTipo mapea cada valor del enum a su representación en texto,
// que es la misma que se almacena en la columna tipo_anomalia de la BD.
var nombrePorTipo = map[TipoAnomalia]string{
	TipoAnomaliaAnomalia:                "ANOMALIA",
	TipoAnomaliaIncidencia:              "INCIDENCIA",
	TipoAnomaliaReporteConductor:        "REPORTE_CONDUCTOR",
	TipoAnomaliaReporteFallaCritica:     "REPORTE_FALLA_CRITICA",
	TipoAnomaliaSeguimientoFallaCritica: "SEGUIMIENTO_FALLA_CRITICA",
}

var tipoPorNombre = map[string]TipoAnomalia{
	"ANOMALIA":                  TipoAnomaliaAnomalia,
	"INCIDENCIA":                TipoAnomaliaIncidencia,
	"REPORTE_CONDUCTOR":         TipoAnomaliaReporteConductor,
	"REPORTE_FALLA_CRITICA":     TipoAnomaliaReporteFallaCritica,
	"SEGUIMIENTO_FALLA_CRITICA": TipoAnomaliaSeguimientoFallaCritica,
}

// TiposAnomaliaValidos contiene todos los valores permitidos (en texto)
// para tipo_anomalia.
func TiposAnomaliaValidos() []string {
	return []string{
		TipoAnomaliaAnomalia.String(),
		TipoAnomaliaIncidencia.String(),
		TipoAnomaliaReporteConductor.String(),
		TipoAnomaliaReporteFallaCritica.String(),
		TipoAnomaliaSeguimientoFallaCritica.String(),
	}
}

// String devuelve la representación en texto del enum (la misma que se
// guarda en BD y se expone en la API).
func (t TipoAnomalia) String() string {
	if nombre, ok := nombrePorTipo[t]; ok {
		return nombre
	}
	return "DESCONOCIDO"
}

// EsValido indica si el valor pertenece a alguno de los 5 tipos soportados.
func (t TipoAnomalia) EsValido() bool {
	_, ok := nombrePorTipo[t]
	return ok
}

// ParseTipoAnomalia convierte un texto (recibido en un request o leído de
// la BD) al enum TipoAnomalia. Devuelve error si el valor no coincide con
// ninguno de los 5 tipos soportados.
func ParseTipoAnomalia(valor string) (TipoAnomalia, error) {
	if tipo, ok := tipoPorNombre[valor]; ok {
		return tipo, nil
	}
	return 0, fmt.Errorf("tipo_anomalia inválido: %q. Valores permitidos: %v", valor, TiposAnomaliaValidos())
}

// EsTipoAnomaliaValido verifica si el texto pertenece al enum soportado.
func EsTipoAnomaliaValido(tipo string) bool {
	_, ok := tipoPorNombre[tipo]
	return ok
}

// MarshalJSON serializa el enum como su representación en texto.
func (t TipoAnomalia) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// UnmarshalJSON deserializa el enum desde su representación en texto,
// rechazando cualquier valor que no sea uno de los predefinidos.
func (t *TipoAnomalia) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	tipo, err := ParseTipoAnomalia(s)
	if err != nil {
		return err
	}
	*t = tipo
	return nil
}

// Value implementa driver.Valuer para que pgx/database-sql persistan el
// enum como su representación en texto en la columna tipo_anomalia.
func (t TipoAnomalia) Value() (driver.Value, error) {
	if !t.EsValido() {
		return nil, fmt.Errorf("tipo_anomalia inválido: %d", t)
	}
	return t.String(), nil
}

// Scan implementa sql.Scanner para leer el enum desde la columna
// tipo_anomalia de la BD.
func (t *TipoAnomalia) Scan(value interface{}) error {
	if value == nil {
		return fmt.Errorf("tipo_anomalia no puede ser nulo")
	}

	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("tipo de dato no soportado para tipo_anomalia: %T", value)
	}

	tipo, err := ParseTipoAnomalia(s)
	if err != nil {
		return err
	}
	*t = tipo
	return nil
}
