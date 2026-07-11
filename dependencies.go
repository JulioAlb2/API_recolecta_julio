package main

import (
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vicpoo/API_recolecta/docs"
	alertaApplication "github.com/vicpoo/API_recolecta/src/alerta_usuario/application"
	alertaHttp "github.com/vicpoo/API_recolecta/src/alerta_usuario/infrastructure/http"
	alertaPostgres "github.com/vicpoo/API_recolecta/src/alerta_usuario/infrastructure/postgres"
	historialUseCases "github.com/vicpoo/API_recolecta/src/Camion/application"
	rutaCamionApp "github.com/vicpoo/API_recolecta/src/Camion/application"
	tipoCamionUseCases "github.com/vicpoo/API_recolecta/src/Camion/application"
	historialAdapters "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/adapters"
	rutaCamionAdapters "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/adapters"
	tipoCamionAdapters "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/adapters"
	historialControllers "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/controllers"
	rutaCamionControllers "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/controllers"
	tipoCamionControllers "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/controllers"
	historialRoutes "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/routes"
	rutaCamionRoutes "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/routes"
	tipoCamionRoutes "github.com/vicpoo/API_recolecta/src/Camion/infrastructure/routes"
	camionUseCases "github.com/vicpoo/API_recolecta/src/Rutas/application"
	estadoCamionUseCases "github.com/vicpoo/API_recolecta/src/Rutas/application"
	puntoUseCases "github.com/vicpoo/API_recolecta/src/Rutas/application"
	registroVaciadoApplication "github.com/vicpoo/API_recolecta/src/Rutas/application"
	rsApplication "github.com/vicpoo/API_recolecta/src/Rutas/application"
	rutaUseCases "github.com/vicpoo/API_recolecta/src/Rutas/application"
	camionAdapters "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/adapters"
	estadoCamionAdapters "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/adapters"
	puntoAdapters "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/adapters"
	registroVaciadoAdapters "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/adapters"
	rsAdapters "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/adapters"
	rutaAdapters "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/adapters"
	camionControllers "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	estadoCamionControllers "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	puntoControllers "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	registroVaciadoControllers "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	rsControllers "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	rutaControllers "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/controllers"
	camionRoutes "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	estadoCamionRoutes "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	puntoRoutes "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	registroVaciadoRoutesPkg "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	rsRoutes "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	rutaRoutes "github.com/vicpoo/API_recolecta/src/Rutas/infraestructure/routes"
	"github.com/vicpoo/API_recolecta/src/core"

	ciudadanosInfra "github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure"
	ciudadanosRoutes "github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure/routes"
	anomalia "github.com/vicpoo/API_recolecta/src/Fallas/infrastructure"
	incidencia "github.com/vicpoo/API_recolecta/src/Fallas/infrastructure"
	reporteConductor "github.com/vicpoo/API_recolecta/src/Fallas/infrastructure"
	recorrido "github.com/vicpoo/API_recolecta/src/recorrido/infrastructure"
	reporteFallaCritica "github.com/vicpoo/API_recolecta/src/Fallas/infrastructure"
	seguimientoFallaCritica "github.com/vicpoo/API_recolecta/src/Fallas/infrastructure"
	alertaMantenimiento "github.com/vicpoo/API_recolecta/src/Mantenimiento/infrastructure"
	registroMantenimiento "github.com/vicpoo/API_recolecta/src/Mantenimiento/infrastructure"
	reporteMantenimientoGenerado "github.com/vicpoo/API_recolecta/src/Mantenimiento/infrastructure"
	tipoMantenimiento "github.com/vicpoo/API_recolecta/src/Mantenimiento/infrastructure"
	coloniaApplication "github.com/vicpoo/API_recolecta/src/colonia/application"
	coloniaHttp "github.com/vicpoo/API_recolecta/src/colonia/infrastructure/http"
	coloniaPostgres "github.com/vicpoo/API_recolecta/src/colonia/infrastructure/postgres"
	empleadoInfra "github.com/vicpoo/API_recolecta/src/empleado/infrastructure"
	empleadoRoutes "github.com/vicpoo/API_recolecta/src/empleado/infrastructure/routes"
	_ "github.com/vicpoo/API_recolecta/src/notificacion/infrastructure"
	//rolInfra "github.com/vicpoo/API_recolecta/src/rol/infrastructure"
	//listMisAlertasUC "github.com/vicpoo/API_recolecta/src/alerta_usuario/application"
	//marcarLeidaUC "github.com/vicpoo/API_recolecta/src/alerta_usuario/application"
	//alertaRepository "github.com/vicpoo/API_recolecta/src/alerta_usuario/infrastructure/postgres"
	//ciudadanosInfra "github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure"
	//alertaHttp "github.com/vicpoo/API_recolecta/src/alerta_usuario/infrastructure/http"
	//ciudadanosRoutes "github.com/vicpoo/API_recolecta/src/Ciudadanos/infrastructure/routes"
)

// archivo para hacer las instancias de los controllers, casos de uso y repositories, etc.
// @title           API Recolecta
// @version         1.0
// @description     API para gestión de recolección de residuos
// @host            localhost:8080
// @BasePath        /
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description     Pega el token del login. Puedes usar solo el JWT o con prefijo Bearer (ej: Bearer eyJhbG...)
// @security BearerAuth
func InitDependencies() {
	// En producción las variables vienen inyectadas por Docker; .env es opcional.
	_ = godotenv.Load()

	// Host absoluto para Swagger "Try it out" (sin abrir cualquier origen).
	// Prioridad: SWAGGER_HOST → NGROK_DOMAIN. Nunca dejar localhost en despliegue.
	docs.SwaggerInfo.Host = ""
	docs.SwaggerInfo.Schemes = []string{"https"}

	rawHost := strings.TrimSpace(os.Getenv("SWAGGER_HOST"))
	if rawHost == "" {
		rawHost = strings.TrimSpace(os.Getenv("NGROK_DOMAIN"))
	}
	rawHost = strings.TrimPrefix(strings.TrimPrefix(rawHost, "https://"), "http://")
	rawHost = strings.TrimSpace(rawHost)

	if rawHost == "" || strings.EqualFold(rawHost, "auto") {
		log.Println("Swagger: SWAGGER_HOST/NGROK_DOMAIN vacío — UI usará el origen del navegador")
		docs.SwaggerInfo.Host = ""
		docs.SwaggerInfo.Schemes = []string{}
	} else {
		docs.SwaggerInfo.Host = rawHost
		// ngrok y producción van por HTTPS; http solo como respaldo local.
		if strings.Contains(rawHost, "ngrok") || strings.ToLower(os.Getenv("ENVIRONMENT")) == "production" {
			docs.SwaggerInfo.Schemes = []string{"https"}
		} else {
			docs.SwaggerInfo.Schemes = []string{"https", "http"}
		}
	}
	log.Printf("Swagger host=%q schemes=%v", docs.SwaggerInfo.Host, docs.SwaggerInfo.Schemes)
	configureSwaggerDocs()

	engine := gin.Default()
	engine.Use(core.CORSMiddleware())
	engine.Use(core.InjectNgrokSkipHeaderInSwagger())
	engine.GET("/api/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.PersistAuthorization(true),
	))

	db := core.GetBD()

	//tipo camion
	tipoCamionRepository := tipoCamionAdapters.NewPostgresTipoCamion()
	saveTipoCamionUc := tipoCamionUseCases.NewSaveTipoCamionUseCase(tipoCamionRepository)
	listAllTipoCamionUc := tipoCamionUseCases.NewListAllTipoCamion(tipoCamionRepository)
	getTipoCamionUc := tipoCamionUseCases.NewGetTipoCamionByNameUseCase(tipoCamionRepository)
	deleteTipoCamionByIdUc := tipoCamionUseCases.NewDeleteTipoCamionUseCase(tipoCamionRepository)

	createTipoCamionCtr := tipoCamionControllers.NewCreateTipoCamionController(saveTipoCamionUc)
	getAllTipoCamionCtr := tipoCamionControllers.NewGetAllTipoCamionController(listAllTipoCamionUc)
	getTipoCamionByNameCtr := tipoCamionControllers.NewGetTipoCamionByNameController(getTipoCamionUc)
	deleteTipoCamionByIdCtr := tipoCamionControllers.NewDeleteTipoCamionController(deleteTipoCamionByIdUc)

	tipoCamionRoutes := tipoCamionRoutes.NewTipoCamionRoutes(
		engine,
		createTipoCamionCtr,
		getAllTipoCamionCtr,
		getTipoCamionByNameCtr,
		deleteTipoCamionByIdCtr,
	)
	tipoCamionRoutes.Run()

	//camion
	camionRepository := camionAdapters.NewPostgresCamion()
	saveCamionUc := camionUseCases.NewSaveCamionUseCase(camionRepository)
	listAllCamionUc := camionUseCases.NewListCamionUseCase(camionRepository)
	updateCamionUc := camionUseCases.NewUpdateCamionUseCase(camionRepository)
	deleteCamionByIdUc := camionUseCases.NewDeleteCamionUseCase(camionRepository)
	getCamionByIdUc := camionUseCases.NewGetCamionByIDUseCase(camionRepository)
	getCamionByPlacaUc := camionUseCases.NewGetCamionByPlacaUseCase(camionRepository)
	getCamionByModeloUc := camionUseCases.NewGetCamionByModeloUseCase(camionRepository)

	createCamionCtr := camionControllers.NewCreateCamionController(saveCamionUc)
	getAllCamionCtr := camionControllers.NewGetAllCamionController(listAllCamionUc)
	updateCamionCtr := camionControllers.NewUpdateCamionController(updateCamionUc)
	deleteCamionByIdCtr := camionControllers.NewDeleteCamionController(deleteCamionByIdUc)
	getCamionByIdCtr := camionControllers.NewGetCamionByIDController(getCamionByIdUc)
	getCamionByPlacaCtr := camionControllers.NewGetCamionByPlacaController(getCamionByPlacaUc)
	getCamionByModeloCtr := camionControllers.NewGetCamionByModeloController(getCamionByModeloUc)

	camionRoutes := camionRoutes.NewCamionRoutes(
		engine, createCamionCtr,
		getAllCamionCtr,
		getCamionByIdCtr,
		updateCamionCtr,
		deleteCamionByIdCtr,
		getCamionByPlacaCtr,
		getCamionByModeloCtr,
	)
	camionRoutes.Run()

	//estado camion
	estadoCamionRepository := estadoCamionAdapters.NewPostgresEstadoCamion()

	saveEstadoCamionUc := estadoCamionUseCases.NewSaveEstadoCamionUseCase(estadoCamionRepository)
	listEstadoCamionUc := estadoCamionUseCases.NewListAllEstadoCamionUseCase(estadoCamionRepository)
	getEstadoCamionByIdUc := estadoCamionUseCases.NewGetByIdEstadoCamionUseCase(estadoCamionRepository)
	updateEstadoCamionUc := estadoCamionUseCases.NewUpdateEstadoCamionUseCase(estadoCamionRepository)
	deleteEstadoCamionUc := estadoCamionUseCases.NewDeleteEstadoCamionUseCase(estadoCamionRepository)

	createEstadoCamionCtr := estadoCamionControllers.NewCreateEstadoCamionController(saveEstadoCamionUc)
	getAllEstadoCamionCtr := estadoCamionControllers.NewGetAllEstadoCamionController(listEstadoCamionUc)
	getEstadoCamionByIdCtr := estadoCamionControllers.NewGetEstadoCamionByIdController(getEstadoCamionByIdUc)
	updateEstadoCamionCtr := estadoCamionControllers.NewUpdateEstadoCamionController(updateEstadoCamionUc)
	deleteEstadoCamionCtr := estadoCamionControllers.NewDeleteEstadoCamionController(deleteEstadoCamionUc)

	estadoCamionRoutes := estadoCamionRoutes.NewEstadoCamionRoutes(
		engine,
		createEstadoCamionCtr,
		getAllEstadoCamionCtr,
		getEstadoCamionByIdCtr,
		deleteEstadoCamionCtr,
		updateEstadoCamionCtr,
	)

	estadoCamionRoutes.Run()

	// ================================
	// HISTORIAL ASIGNACION CAMION
	// ================================
	historialRepository := historialAdapters.NewPostgresHistorialAsignacionCamion()

	createHistorialUC := historialUseCases.NewSaveHistorialAsignacionCamionUseCase(historialRepository)
	getAllHistorialUC := historialUseCases.NewListAllHistorialAsignacionCamionUseCase(historialRepository)
	getHistorialByIdUC := historialUseCases.NewGetHistorialAsignacionCamionByIdUseCase(historialRepository)
	updateHistorialUC := historialUseCases.NewUpdateHistorialAsignacionCamionUseCase(historialRepository)
	deleteHistorialUC := historialUseCases.NewDeleteHistorialAsignacionCamionUseCase(historialRepository)

	getByCamionUC := historialUseCases.NewGetHistorialByCamionUseCase(historialRepository)
	getByChoferUC := historialUseCases.NewGetHistorialByChoferUseCase(historialRepository)
	getActivoByCamionUC := historialUseCases.NewGetActivoByCamionUseCase(historialRepository)
	getActivoByChoferUC := historialUseCases.NewGetActivoByChoferUseCase(historialRepository)

	darDeBajaUC := historialUseCases.NewDarDeBajaHistorialAsignacionUseCase(historialRepository)
	cerrarCamionUC := historialUseCases.NewCerrarAsignacionActivaCamionUseCase(historialRepository)
	cerrarChoferUC := historialUseCases.NewCerrarAsignacionActivaChoferUseCase(historialRepository)

	createHistorialCtr := historialControllers.NewCreateHistorialAsignacionCamionController(createHistorialUC)
	getAllHistorialCtr := historialControllers.NewGetAllHistorialAsignacionCamionController(getAllHistorialUC)
	getHistorialByIdCtr := historialControllers.NewGetHistorialAsignacionByIdController(getHistorialByIdUC)
	updateHistorialCtr := historialControllers.NewUpdateHistorialAsignacionCamionController(updateHistorialUC)
	deleteHistorialCtr := historialControllers.NewDeleteHistorialAsignacionCamionController(deleteHistorialUC)

	getByCamionCtr := historialControllers.NewGetHistorialByCamionController(getByCamionUC)
	getByChoferCtr := historialControllers.NewGetHistorialByChoferController(getByChoferUC)
	getActivoByCamionCtr := historialControllers.NewGetActivoByCamionController(getActivoByCamionUC)
	getActivoByChoferCtr := historialControllers.NewGetActivoByChoferController(getActivoByChoferUC)

	darDeBajaCtr := historialControllers.NewDarDeBajaHistorialAsignacionController(darDeBajaUC)
	cerrarCamionCtr := historialControllers.NewCerrarAsignacionActivaCamionController(cerrarCamionUC)
	cerrarChoferCtr := historialControllers.NewCerrarAsignacionActivaChoferController(cerrarChoferUC)

	historialRoutes := historialRoutes.NewHistorialAsignacionCamionRoutes(
		engine,
		createHistorialCtr,
		getAllHistorialCtr,
		getHistorialByIdCtr,
		updateHistorialCtr,
		deleteHistorialCtr,
		getByCamionCtr,
		getByChoferCtr,
		getActivoByCamionCtr,
		getActivoByChoferCtr,
		darDeBajaCtr,
		cerrarCamionCtr,
		cerrarChoferCtr,
	)

	historialRoutes.Run()

	// ================================
	// RUTA
	// ================================

	rutaRepository := rutaAdapters.NewPostgresRuta()

	createRutaUc := rutaUseCases.NewCreateRutaUseCase(rutaRepository)
	getAllRutaUc := rutaUseCases.NewListAllRutaUseCase(rutaRepository)
	getRutaByIdUc := rutaUseCases.NewGetRutaByIdUseCase(rutaRepository)
	updateRutaUc := rutaUseCases.NewUpdateRutaUseCase(rutaRepository)
	deleteRutaUc := rutaUseCases.NewDeleteRutaUseCase(rutaRepository)
	getRutasActivasUc := rutaUseCases.NewGetRutaActivasUseCase(rutaRepository)

	createRutaCtr := rutaControllers.NewCreateRutaController(createRutaUc)
	getAllRutaCtr := rutaControllers.NewGetAllRutaController(getAllRutaUc)
	getRutaByIdCtr := rutaControllers.NewGetRutaByIdController(getRutaByIdUc)
	updateRutaCtr := rutaControllers.NewUpdateRutaController(updateRutaUc)
	deleteRutaCtr := rutaControllers.NewDeleteRutaController(deleteRutaUc)
	getRutasActivasCtr := rutaControllers.NewGetRutaActivasController(getRutasActivasUc)

	rutaRoutes := rutaRoutes.NewRutaRoutes(
		engine,
		createRutaCtr,
		getAllRutaCtr,
		getRutaByIdCtr,
		updateRutaCtr,
		deleteRutaCtr,
		getRutasActivasCtr,
	)

	rutaRoutes.Run()

	puntoRepository := puntoAdapters.NewPostgresPuntoRecoleccion()

	createPuntoUC := puntoUseCases.NewSavePuntoRecoleccionUseCase(puntoRepository)
	updatePuntoUC := puntoUseCases.NewUpdatePuntoRecoleccionUseCase(puntoRepository)
	getAllPuntoUC := puntoUseCases.NewListAllPuntoRecoleccionUseCase(puntoRepository)
	getPuntoByIdUC := puntoUseCases.NewGetPuntoRecoleccionByIdUseCase(puntoRepository)
	getPuntoByRutaUC := puntoUseCases.NewGetPuntoRecoleccionByRutaUseCase(puntoRepository)
	deletePuntoUC := puntoUseCases.NewDeletePuntoRecoleccionUseCase(puntoRepository)

	createPuntoCTR := puntoControllers.NewCreatePuntoRecoleccionController(createPuntoUC)
	updatePuntoCTR := puntoControllers.NewUpdatePuntoRecoleccionController(updatePuntoUC)
	getAllPuntoCTR := puntoControllers.NewGetAllPuntoRecoleccionController(getAllPuntoUC)
	getPuntoByIdCTR := puntoControllers.NewGetPuntoRecoleccionByIdController(getPuntoByIdUC)
	getPuntoByRutaCTR := puntoControllers.NewGetPuntoRecoleccionByRutaController(getPuntoByRutaUC)
	deletePuntoCTR := puntoControllers.NewDeletePuntoRecoleccionController(deletePuntoUC)

	puntoRoutes := puntoRoutes.NewPuntoRecoleccionRoutes(
		engine,
		createPuntoCTR,
		getAllPuntoCTR,
		getPuntoByIdCTR,
		getPuntoByRutaCTR,
		updatePuntoCTR,
		deletePuntoCTR,
	)

	puntoRoutes.Run()

	rellenoRepo := rsAdapters.NewPostgresRellenoSanitario()

	createRellenoUC := rsApplication.NewSaveRellenoSanitarioUseCase(rellenoRepo)
	updateRellenoUC := rsApplication.NewUpdateRellenoSanitarioUseCase(rellenoRepo)
	getAllRellenoUC := rsApplication.NewListRellenoSanitarioUseCase(rellenoRepo)
	getRellenoByIDUC := rsApplication.NewGetRellenoSanitarioByIdUseCase(rellenoRepo)
	deleteRellenoUC := rsApplication.NewDeleteRellenoSanitarioUseCase(rellenoRepo)
	getRellenoByNombreUC := rsApplication.NewGetRellenoSanitarioByNombreUseCase(rellenoRepo)
	existsRellenoUC := rsApplication.NewExistsRellenoSanitarioByIdUseCase(rellenoRepo)

	createRellenoController := rsControllers.NewCreateRellenoSanitarioController(createRellenoUC)
	updateRellenoController := rsControllers.NewUpdateRellenoSanitarioController(updateRellenoUC)
	getAllRellenoController := rsControllers.NewGetAllRellenoSanitarioController(getAllRellenoUC)
	getRellenoByIDController := rsControllers.NewGetRellenoSanitarioByIDController(getRellenoByIDUC)
	deleteRellenoController := rsControllers.NewDeleteRellenoSanitarioController(deleteRellenoUC)
	getRellenoByNombreController := rsControllers.NewGetRellenoSanitarioByNombreController(getRellenoByNombreUC)
	existsRellenoController := rsControllers.NewExistsRellenoSanitarioByIdController(existsRellenoUC)

	rellenoRoutes := rsRoutes.NewRellenoSanitarioRoutes(
		engine,
		createRellenoController,
		getAllRellenoController,
		getRellenoByIDController,
		updateRellenoController,
		deleteRellenoController,
		getRellenoByNombreController,
		existsRellenoController,
	)

	rellenoRoutes.Run()

	repository := rutaCamionAdapters.NewPostgresRutaCamion()

	// ===============================
	// USE CASES
	// ===============================
	createRutaCamionUC := rutaCamionApp.NewSaveRutaCamionUseCase(repository)
	updateRutaCamionUC := rutaCamionApp.NewUpdateRutaCamionUseCase(repository)
	getAllRutaCamionUC := rutaCamionApp.NewListAllRutaCamionUseCase(repository)
	getRutaCamionByIDUC := rutaCamionApp.NewGetRutaCamionByIDUseCase(repository)
	getRutaCamionByCamionIDUC := rutaCamionApp.NewGetRutaCamionByCamionIDUseCase(repository)
	getRutaCamionByRutaIDUC := rutaCamionApp.NewGetRutaCamionByRutaIDUseCase(repository)
	existsRutaCamionUC := rutaCamionApp.NewExistsRutaCamionByIDUseCase(repository)
	deleteRutaCamionUC := rutaCamionApp.NewDeleteRutaCamionUseCase(repository)

	// ===============================
	// CONTROLLERS
	// ===============================
	createRutaCamionController :=
		rutaCamionControllers.NewCreateRutaCamionController(createRutaCamionUC)

	updateRutaCamionController :=
		rutaCamionControllers.NewUpdateRutaCamionController(updateRutaCamionUC)

	getAllRutaCamionController :=
		rutaCamionControllers.NewGetAllRutaCamionController(getAllRutaCamionUC)

	getRutaCamionByIDController :=
		rutaCamionControllers.NewGetRutaCamionByIDController(getRutaCamionByIDUC)

	getRutaCamionByCamionIDController :=
		rutaCamionControllers.NewGetRutaCamionByCamionIDController(getRutaCamionByCamionIDUC)

	getRutaCamionByRutaIDController :=
		rutaCamionControllers.NewGetRutaCamionByRutaIDController(getRutaCamionByRutaIDUC)

	existsRutaCamionController :=
		rutaCamionControllers.NewExistsRutaCamionByIDController(existsRutaCamionUC)

	deleteRutaCamionController :=
		rutaCamionControllers.NewDeleteRutaCamionController(deleteRutaCamionUC)

	rutaCamionRoutes := rutaCamionRoutes.NewRutaCamionRoutes(
		engine,
		createRutaCamionController,
		getAllRutaCamionController,
		getRutaCamionByIDController,
		getRutaCamionByCamionIDController,
		getRutaCamionByRutaIDController,
		existsRutaCamionController,
		updateRutaCamionController,
		deleteRutaCamionController,
	)

	rutaCamionRoutes.Run()

	// ===============================
	// REGISTRO VACIADO
	// ===============================

	// Repository
	registroVaciadoRepository := registroVaciadoAdapters.NewPostgresRegistroVaciado()

	// ===============================
	// USE CASES
	// ===============================
	createRegistroVaciadoUC := registroVaciadoApplication.NewCreateRegistroVaciadoUseCase(registroVaciadoRepository)
	getAllRegistroVaciadoUC := registroVaciadoApplication.NewListAllRegistroVaciadoUseCase(registroVaciadoRepository)
	getRegistroVaciadoByIDUC := registroVaciadoApplication.NewGetRegistroVaciadoByIDUseCase(registroVaciadoRepository)
	getRegistroVaciadoByRellenoIDUC := registroVaciadoApplication.NewGetRegistroVaciadoByRellenoIDUseCase(registroVaciadoRepository)
	getRegistroVaciadoByRutaCamionIDUC := registroVaciadoApplication.NewGetRegistroVaciadoByRutaCamionIDUseCase(registroVaciadoRepository)
	existsRegistroVaciadoUC := registroVaciadoApplication.NewExistsRegistroVaciadoUseCase(registroVaciadoRepository)
	deleteRegistroVaciadoUC := registroVaciadoApplication.NewDeleteRegistroVaciadoUseCase(registroVaciadoRepository)

	// ===============================
	// CONTROLLERS
	// ===============================
	createRegistroVaciadoController := registroVaciadoControllers.NewCreateRegistroVaciadoController(createRegistroVaciadoUC)
	getAllRegistroVaciadoController := registroVaciadoControllers.NewGetAllRegistroVaciadoController(getAllRegistroVaciadoUC)
	getRegistroVaciadoByIDController := registroVaciadoControllers.NewGetRegistroVaciadoByIDController(getRegistroVaciadoByIDUC)
	getRegistroVaciadoByRellenoIDController := registroVaciadoControllers.NewGetRegistroVaciadoByRellenoIDController(getRegistroVaciadoByRellenoIDUC)
	getRegistroVaciadoByRutaCamionIDController := registroVaciadoControllers.NewGetRegistroVaciadoByRutaCamionIDController(getRegistroVaciadoByRutaCamionIDUC)
	existsRegistroVaciadoController := registroVaciadoControllers.NewExistsRegistroVaciadoController(existsRegistroVaciadoUC)
	deleteRegistroVaciadoController := registroVaciadoControllers.NewDeleteRegistroVaciadoController(deleteRegistroVaciadoUC)

	// ===============================
	// ROUTES
	// ===============================
	registroVaciadoRoutes := registroVaciadoRoutesPkg.NewRegistroVaciadoRoutes(
		engine,
		createRegistroVaciadoController,
		getAllRegistroVaciadoController,
		getRegistroVaciadoByIDController,
		getRegistroVaciadoByRellenoIDController,
		getRegistroVaciadoByRutaCamionIDController,
		existsRegistroVaciadoController,
		deleteRegistroVaciadoController,
	)

	registroVaciadoRoutes.Run()

	// ===============================
	// COLONIA
	// ===============================

	coloniaRepository := coloniaPostgres.NewColoniaRepository()

	createColoniaUC := coloniaApplication.NewCreateColonia(coloniaRepository)
	getColoniaUC := coloniaApplication.NewGetColonia(coloniaRepository)
	listColoniasUC := coloniaApplication.NewListColonias(coloniaRepository)
	updateColoniaUC := coloniaApplication.NewUpdateColonia(coloniaRepository)
	deleteColoniaUC := coloniaApplication.NewDeleteColonia(coloniaRepository)

	coloniaController := coloniaHttp.NewColoniaController(
		createColoniaUC,
		getColoniaUC,
		listColoniasUC,
		updateColoniaUC,
		deleteColoniaUC,
	)

	coloniaController.RegisterRoutes(engine)

	// ===============================
	// Ciudadanos
	//===============================

	/*
		ciudadanosRoutes.CiudadanoRoutes(
			engine,
			ciudadanoDeps.CreateCiudadanoController,
			ciudadanoDeps.GetCiudadanoController,
			ciudadanoDeps.ListCiudadanoController,
			ciudadanoDeps.UpdateCiudadanoController,
			ciudadanoDeps.DeleteCiudadanoController,
			ciudadanoDeps.LoginCiudadanoController,
		)

		ciudadanosRoutes.DomicilioRoutes(
			engine,
			domicilioDeps.DomicilioController,
		)

	*/
	//==============================
	//Ciudadano

	ciudadanoDeps := ciudadanosInfra.InitCiudadanoDependencies(db)
	domicilioDeps := ciudadanosInfra.InitDomicilioDependencies(db)
	empleadoDeps := empleadoInfra.InitEmpleadoDependencies(db)

	ciudadanosRoutes.CiudadanoRoutes(
		engine,
		ciudadanoDeps.CreateCiudadanoController,
		ciudadanoDeps.GetCiudadanoController,
		ciudadanoDeps.ListCiudadanoController,
		ciudadanoDeps.UpdateCiudadanoController,
		ciudadanoDeps.DeleteCiudadanoController,
		ciudadanoDeps.LoginCiudadanoController,
	)

	ciudadanosRoutes.DomicilioRoutes(
		engine,
		domicilioDeps.DomicilioController,
	)

	empleadoRoutes.EmpleadoRoutes(
		engine,
		empleadoDeps.CreateEmpleadoController,
		empleadoDeps.ListEmpleadoController,
		empleadoDeps.GetEmpleadoController,
		empleadoDeps.UpdateEmpleadoController,
		empleadoDeps.DeleteEmpleadoController,
		empleadoDeps.LoginEmpleadoController,
	)

	// ===============================
	// ALERTA USUARIO
	// ===============================
	alertaRepository := alertaPostgres.NewPostgresAlertaRepository(db)
	createAlertaUC := alertaApplication.NewCreateAlerta(alertaRepository)
	listMisAlertasUC := alertaApplication.NewListMisAlertas(alertaRepository)
	marcarLeidaUC := alertaApplication.NewMarcarLeida(alertaRepository)
	alertaController := alertaHttp.NewAlertaController(createAlertaUC, listMisAlertasUC, marcarLeidaUC)

	apiGroup := engine.Group("/api")
	alertaController.RegisterRoutes(apiGroup)

	// ===============================
	// FALLAS Y MANTENIMIENTO
	// ===============================

	anomaliaRoutes := anomalia.NewAnomaliaRouter(engine)

	anomaliaRoutes.Run()

	incidenciaRoutes := incidencia.NewIncidenciaRouter(engine)

	incidenciaRoutes.Run()

	reporteConductorRoutes := reporteConductor.NewReporteConductorRouter(engine)

	reporteConductorRoutes.Run()

	registroMantenimientoRoutes := registroMantenimiento.NewRegistroMantenimientoRouter(engine)

	registroMantenimientoRoutes.Run()

	reporteFallaCriticaRoutes := reporteFallaCritica.NewReporteFallaCriticaRouter(engine)

	reporteFallaCriticaRoutes.Run()

	alertaMantenimientoRoutes := alertaMantenimiento.NewAlertaMantenimientoRouter(engine)

	alertaMantenimientoRoutes.Run()

	reporteMantenimientoGeneradoRoutes := reporteMantenimientoGenerado.NewReporteMantenimientoGeneradoRouter(engine)

	reporteMantenimientoGeneradoRoutes.Run()

	seguimientoFallaCriticaRoutes := seguimientoFallaCritica.NewSeguimientoFallaCriticaRouter(engine)

	seguimientoFallaCriticaRoutes.Run()

	tipoMantenimientoRoutes := tipoMantenimiento.NewTipoMantenimientoRouter(engine)

	tipoMantenimientoRoutes.Run()

	recorridoRoutes := recorrido.NewRecorridoRouter(engine)
	recorridoRoutes.Run()

	engine.Run(":8080")
}
