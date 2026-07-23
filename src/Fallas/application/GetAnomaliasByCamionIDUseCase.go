// GetAnomaliasByCamionIDUseCase.go
package application

import (
	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
)

type GetAnomaliasByCamionIDUseCase struct {
	repo repositories.IAnomalia
}

func NewGetAnomaliasByCamionIDUseCase(repo repositories.IAnomalia) *GetAnomaliasByCamionIDUseCase {
	return &GetAnomaliasByCamionIDUseCase{repo: repo}
}

func (uc *GetAnomaliasByCamionIDUseCase) Run(camionID int32) ([]entities.Anomalia, error) {
	return uc.repo.GetByCamionID(camionID)
}
