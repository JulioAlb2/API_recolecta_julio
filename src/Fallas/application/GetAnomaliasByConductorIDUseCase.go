// GetAnomaliasByConductorIDUseCase.go
package application

import (
	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByConductorIDUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByConductorIDUseCase(repo repositories.IAnomalia) *GetAnomaliasByConductorIDUseCase {
	return &GetAnomaliasByConductorIDUseCase{repo: repo}
}

func (uc *GetAnomaliasByConductorIDUseCase) Run(conductorID int32) ([]entities.Anomalia, error) {
	return uc.repo.GetByConductorID(conductorID)
}
