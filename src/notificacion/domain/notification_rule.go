package domain

import "context"

type NotificationRule struct {
	TenantID  int    `json:"tenant_id"`
	StateCode string `json:"state_code"`
	Enabled   bool   `json:"enabled"`
	Threshold int    `json:"threshold"`
}

// INotificationRuleRepository: reglas por tenant (municipio) -- ver
// docs/07-plan-multitenancy.md Fase 8. Redis no tiene RLS, asi que tenantID
// obligatorio en cada metodo es la unica barrera de aislamiento, mismo
// criterio que modelo-reportes/clasificador-reportes.
type INotificationRuleRepository interface {
	Save(ctx context.Context, tenantID int, rule NotificationRule) error
	GetByStateCode(ctx context.Context, tenantID int, code string) (*NotificationRule, error)
	List(ctx context.Context, tenantID int) ([]NotificationRule, error)
	Delete(ctx context.Context, tenantID int, code string) error
}
