package infrastructure

import (
	"github.com/jackc/pgx/v5/pgxpool"

	application_ciudadano "github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_ciudadano"
	application_domicilio "github.com/vicpoo/API_recolecta/src/Ciudadanos/application/application_domicilio"

	controller_ciudadano "github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure/controller/controller_ciudadano"
	httpDomicilio "github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure/controller/controller_domicilio"

	"github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure/repository"
)

type CiudadanoDependencies struct {
	CreateCiudadanoController *controller_ciudadano.CreateCiudadanoController
	GetCiudadanoController    *controller_ciudadano.GetCiudadanoController
	ListCiudadanoController   *controller_ciudadano.ListCiudadanoController
	UpdateCiudadanoController *controller_ciudadano.UpdateCiudadanoController
	DeleteCiudadanoController *controller_ciudadano.DeleteCiudadanoController
	LoginCiudadanoController  *controller_ciudadano.LoginCiudadanoController
	UpdateFCMTokenController  *controller_ciudadano.UpdateFCMTokenController
}

func InitCiudadanoDependencies(db *pgxpool.Pool) *CiudadanoDependencies {
	ciudadanoRepo := repository.NewCiudadanoPostgresRepository(db)

	createUseCase := application_ciudadano.NewCreateCiudadano(ciudadanoRepo)
	getUseCase := application_ciudadano.NewViewOneCiudadano(ciudadanoRepo)
	listUseCase := application_ciudadano.NewViewAllCiudadano(ciudadanoRepo)
	updateUseCase := application_ciudadano.NewUpdateCiudadano(ciudadanoRepo)
	deleteUseCase := application_ciudadano.NewDeleteCiudadano(ciudadanoRepo)
	loginUseCase := application_ciudadano.NewLoginCiudadano(ciudadanoRepo)
	updateFCMUseCase := application_ciudadano.NewUpdateFCMToken()

	createController := controller_ciudadano.NewCreateCiudadanoController(createUseCase)
	getController := controller_ciudadano.NewGetCiudadanoController(getUseCase)
	listController := controller_ciudadano.NewListCiudadanoController(listUseCase)
	updateController := controller_ciudadano.NewUpdateCiudadanoController(updateUseCase)
	deleteController := controller_ciudadano.NewDeleteCiudadanoController(deleteUseCase)
	loginController := controller_ciudadano.NewLoginCiudadanoController(loginUseCase)
	updateFCMController := controller_ciudadano.NewUpdateFCMTokenController(updateFCMUseCase)

	return &CiudadanoDependencies{
		CreateCiudadanoController: createController,
		GetCiudadanoController:    getController,
		ListCiudadanoController:   listController,
		UpdateCiudadanoController: updateController,
		DeleteCiudadanoController: deleteController,
		LoginCiudadanoController:  loginController,
		UpdateFCMTokenController:  updateFCMController,
	}
}

type DomicilioDependencies struct {
	DomicilioController *httpDomicilio.DomicilioController
}

func InitDomicilioDependencies(db *pgxpool.Pool) *DomicilioDependencies {
	domicilioRepo := repository.NewDomicilioPostgresRepository(db)

	createUseCase := application_domicilio.NewCreateDomicilio(domicilioRepo)
	getUseCase := application_domicilio.NewGetDomicilio(domicilioRepo)
	updateUseCase := application_domicilio.NewUpdateDomicilio(domicilioRepo)
	listUseCase := application_domicilio.NewListDomicilios(domicilioRepo)
	deleteUseCase := application_domicilio.NewDeleteDomicilio(domicilioRepo)

	domicilioController := httpDomicilio.NewDomicilioController(
		createUseCase,
		getUseCase,
		updateUseCase,
		listUseCase,
		deleteUseCase,
	)

	return &DomicilioDependencies{
		DomicilioController: domicilioController,
	}
}