package application

import (
	"context"

	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type ManageNotificationRulesUseCase struct {
	repo domain.INotificationRuleRepository
}

func NewManageNotificationRulesUseCase(repo domain.INotificationRuleRepository) *ManageNotificationRulesUseCase {
	return &ManageNotificationRulesUseCase{repo: repo}
}

func (uc *ManageNotificationRulesUseCase) GetByStateCode(ctx context.Context, tenantID int, code string) (*domain.NotificationRule, error) {
	return uc.repo.GetByStateCode(ctx, tenantID, code)
}

func (uc *ManageNotificationRulesUseCase) List(ctx context.Context, tenantID int) ([]domain.NotificationRule, error) {
	return uc.repo.List(ctx, tenantID)
}

func (uc *ManageNotificationRulesUseCase) Upsert(ctx context.Context, tenantID int, rule domain.NotificationRule) error {
	return uc.repo.Save(ctx, tenantID, rule)
}

func (uc *ManageNotificationRulesUseCase) Delete(ctx context.Context, tenantID int, code string) error {
	return uc.repo.Delete(ctx, tenantID, code)
}
