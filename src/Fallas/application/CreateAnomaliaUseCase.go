// CreateAnomaliaUseCase.go
package application

import (
	"time"

	repositories "github.com/vicpoo/API_recolecta/src/Fallas/domain"
	"github.com/vicpoo/API_recolecta/src/Fallas/domain/entities"
	alertaDomain "github.com/vicpoo/API_recolecta/src/alerta_usuario/domain"
)

type CreateAnomaliaUseCase struct {
	repo       repositories.IAnomalia
	alertaRepo alertaDomain.AlertaUsuarioRepository
	// pipelineUseCase es opcional: si es nil (ej. tests que no lo inyectan),
	// simplemente no se dispara el pipeline de validacion/clasificacion.
	pipelineUseCase *ProcesarPipelineAnomaliaUseCase
}

func NewCreateAnomaliaUseCase(repo repositories.IAnomalia, alertaRepo alertaDomain.AlertaUsuarioRepository, pipelineUseCase *ProcesarPipelineAnomaliaUseCase) *CreateAnomaliaUseCase {
	return &CreateAnomaliaUseCase{repo: repo, alertaRepo: alertaRepo, pipelineUseCase: pipelineUseCase}
}

// tiposConPipeline son los tipos de Anomalia que representan texto libre
// reportado por un ciudadano o conductor, y que por lo tanto tiene sentido
// pasar por modelo_reportes/clasificador_reportes. REPORTE_FALLA_CRITICA y
// SEGUIMIENTO_FALLA_CRITICA quedan fuera: son seguimientos/fallas de camion
// ya capturados por otro flujo, no reportes de calle que el clasificador
// sepa categorizar.
var tiposConPipeline = map[string]bool{
	"ANOMALIA":          true,
	"INCIDENCIA":        true,
	"REPORTE_CONDUCTOR": true,
}

func (uc *CreateAnomaliaUseCase) Run(anomalia *entities.Anomalia) (*entities.Anomalia, error) {
	err := uc.repo.Save(anomalia)
	if err != nil {
		return nil, err
	}

	// OJO: string(anomalia.TipoAnomalia) NO llama a TipoAnomalia.String() --
	// TipoAnomalia es un enum basado en int, y convertir un int a string en
	// Go hace una conversion de rune (numero -> caracter Unicode), no
	// devuelve el texto del enum. Hay que llamar a .String() explicitamente.
	tipo := anomalia.TipoAnomalia.String()

	// Si es una anomalía crítica o incidencia, alertar a supervisores
	if tipo == "REPORTE_FALLA_CRITICA" || tipo == "INCIDENCIA" {
		creadoPor := 0
		if anomalia.ConductorID != nil {
			creadoPor = int(*anomalia.ConductorID)
		}
		alerta := &alertaDomain.AlertaUsuario{
			Titulo:    "Alerta Operativa: " + tipo,
			Mensaje:   anomalia.Descripcion,
			UsuarioID: 1, // Administrador general
			Leida:     false,
			CreadoPor: creadoPor,
			CreatedAt: time.Now(),
		}
		_ = uc.alertaRepo.Create(alerta)
	}

	// Dispara en background la validacion/clasificacion del reporte. No
	// bloquea la respuesta al ciudadano/conductor: modelo_reportes y
	// clasificador_reportes pueden tardar y no deberian retrasar el 201.
	if uc.pipelineUseCase != nil && tiposConPipeline[tipo] {
		origen := "ciudadano"
		if anomalia.ConductorID != nil {
			origen = "conductor"
		}
		// TenantID: por ahora la entidad Anomalia no lo expone (no hay
		// multi-tenant real todavia en este dominio); se usa el default 1
		// que ya usa la columna tenant_id en BD.
		go uc.pipelineUseCase.Run(anomalia.AnomaliaID, anomalia.Descripcion, 1, &origen)
	}

	return anomalia, nil
}