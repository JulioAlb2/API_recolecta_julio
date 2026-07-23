// GetAnomaliasByTipoAnomaliaUseCase.go
package application

import (
	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByTipoAnomaliaUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByTipoAnomaliaUseCase(repo repositories.IAnomalia) *GetAnomaliasByTipoAnomaliaUseCase {
	return &GetAnomaliasByTipoAnomaliaUseCase{repo: repo}
}

func (uc *GetAnomaliasByTipoAnomaliaUseCase) Run(tipoAnomalia entities.TipoAnomalia) ([]entities.Anomalia, error) {
	return uc.repo.GetByTipoAnomalia(tipoAnomalia)
}