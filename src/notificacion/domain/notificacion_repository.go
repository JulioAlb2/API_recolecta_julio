// notificacion_repository.go
package domain

import (
	"github.com/vicpoo/API_recolecta/src/notificacion/domain/entities"
)

type INotificacion interface {
	// Operaciones de consulta (lectura)
	GetAll() ([]entities.Notificacion, error)
	GetByID(id int32) (*entities.Notificacion, error)
	GetByUsuarioID(usuarioID int32) ([]entities.Notificacion, error)
	GetActivasByUsuarioID(usuarioID int32) ([]entities.Notificacion, error)
	GetByTipo(tipo string) ([]entities.Notificacion, error)
	GetActivas() ([]entities.Notificacion, error)
	GetInactivas() ([]entities.Notificacion, error)
	GetByFechaRange(fechaInicio, fechaFin string) ([]entities.Notificacion, error)
	GetGlobales() ([]entities.Notificacion, error)
	GetByCamionID(camionID int32) ([]entities.Notificacion, error)
	GetByCamionYTipo(camionID int32, tipo string) ([]entities.Notificacion, error)
	GetByUsuarioYTipo(usuarioID int32, tipo string) ([]entities.Notificacion, error)
	
	// Métodos de actualización de estado
	MarcarComoLeida(id int32) error
	MarcarComoActiva(id int32) error
	MarcarTodasComoLeidas(usuarioID int32) error
	
	// Métodos de conteo
	CountActivasByUsuarioID(usuarioID int32) (int64, error)
	CountByUsuarioID(usuarioID int32) (int64, error)
	CountByTipo(tipo string) (int64, error)
	CountByCamionID(camionID int32) (int64, error)
}