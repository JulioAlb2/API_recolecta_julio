package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vicpoo/API_recolecta/src/dispositivo/application"
	"github.com/vicpoo/API_recolecta/src/dispositivo/infrastructure/controller"
	"github.com/vicpoo/API_recolecta/src/dispositivo/infrastructure/repository"
)

type DispositivoDependencies struct {
	DispositivoController *controller.DispositivoController
}

func InitDispositivoDependencies(db *pgxpool.Pool) *DispositivoDependencies {
	repo := repository.NewPostgresDispositivoRepository(db)
	useCases := application.NewDispositivoUseCases(repo)
	ctrl := controller.NewDispositivoController(useCases)

	return &DispositivoDependencies{
		DispositivoController: ctrl,
	}
}
