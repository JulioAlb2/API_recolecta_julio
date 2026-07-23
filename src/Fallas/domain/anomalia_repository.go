// anomalia_repository.go
package domain

import (
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type IAnomalia interface {
	// Operaciones CRUD básicas
	Save(anomalia *entities.Anomalia) error
	Update(anomalia *entities.Anomalia) error
	Delete(id int32) error
	GetAll() ([]entities.Anomalia, error)
	GetByID(id int32) (*entities.Anomalia, error)

	// Métodos específicos para Anomalia (cubren los antiguos dominios
	// Incidencia, ReporteConductor, ReporteFallaCritica y
	// SeguimientoFallaCritica, ahora unificados por TipoAnomalia)
	GetByTipoAnomalia(tipoAnomalia entities.TipoAnomalia) ([]entities.Anomalia, error)
	GetByEstado(estado string) ([]entities.Anomalia, error)
	GetByPuntoID(puntoID int32) ([]entities.Anomalia, error)
	GetByConductorID(conductorID int32) ([]entities.Anomalia, error)
	GetByCamionID(camionID int32) ([]entities.Anomalia, error)
	GetByRutaID(rutaID int32) ([]entities.Anomalia, error)
	GetByReferenciaID(referenciaID int32) ([]entities.Anomalia, error)
	GetByFechaRange(fechaInicio, fechaFin string) ([]entities.Anomalia, error)

	// ActualizarPipeline persiste el resultado del pipeline
	// modelo_reportes -> clasificador_reportes sobre una anomalia ya
	// existente. Es una actualizacion dirigida (no pasa por Update/la
	// entidad completa) para no interferir con el flujo CRUD normal.
	// Los punteros nil dejan la columna correspondiente sin tocar.
	ActualizarPipeline(anomaliaID int32, estadoPipeline string, nivelRiesgo *string, inferenciaID *int32, categoria *string, subtipo *string, accion *string, pipelineError *string) error
}
