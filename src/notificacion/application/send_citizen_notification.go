package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type tokenRepository interface {
	GetTokensByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error)
}

type SendCitizenNotificationUseCase struct {
	sender domain.PushNotificationSender
	repo   tokenRepository
}

func NewSendCitizenNotificationUseCase(sender domain.PushNotificationSender, repo tokenRepository) *SendCitizenNotificationUseCase {
	return &SendCitizenNotificationUseCase{sender: sender, repo: repo}
}

func (uc *SendCitizenNotificationUseCase) Execute(ctx context.Context, userIDs []string, notification *domain.PushNotification) (map[string]domain.SendResult, error) {
	tokens, err := uc.repo.GetTokensByUserIDs(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	finalResults := make(map[string]domain.SendResult)
	validTokens := make(map[string]string)

	for _, uid := range userIDs {
		if token, ok := tokens[uid]; ok && token != "" {
			validTokens[uid] = token
		} else {
			finalResults[uid] = domain.SendResult{
				Success: false,
				Error:   "no FCM token found for user",
			}
		}
	}

	if len(validTokens) > 0 {
		results, err := uc.sender.Send(ctx, validTokens, notification)
		if err != nil {
			return nil, err
		}
		for uid, res := range results {
			finalResults[uid] = res
		}
	}

	// Registrar resultados en Redis para auditoría y bandeja de entrada
	rdb := core.GetRedis()
	now := time.Now()
	nowUnix := float64(now.Unix())
	nowStr := now.Format(time.RFC3339)

	for citizenID, res := range finalResults {
		// 1. Añadir a bandeja individual del ciudadano
		inboxKey := fmt.Sprintf("citizen:notifications:%s", citizenID)
		inboxRecord := map[string]interface{}{
			"title":     notification.Title,
			"body":      notification.Body,
			"type":      notification.Type,
			"sent_at":   nowStr,
			"delivered": res.Success,
		}
		if !res.Success {
			inboxRecord["error"] = res.Error
		}
		inboxData, _ := json.Marshal(inboxRecord)

		pipe := rdb.Pipeline()
		pipe.LPush(ctx, inboxKey, string(inboxData))
		pipe.LTrim(ctx, inboxKey, 0, 49) // Mantener últimas 50
		pipe.Expire(ctx, inboxKey, 30*24*time.Hour)

		// 2. Si falló, registrar globalmente en ZSET para consultas administrativas por fecha
		if !res.Success {
			failedRecord := map[string]interface{}{
				"user_id":   citizenID,
				"title":     notification.Title,
				"body":      notification.Body,
				"type":      notification.Type,
				"error":     res.Error,
				"timestamp": now,
			}
			failedData, _ := json.Marshal(failedRecord)
			pipe.ZAdd(ctx, "notifications:failed", redis.Z{
				Score:  nowUnix,
				Member: string(failedData),
			})
			thirtyDaysAgo := now.Add(-30 * 24 * time.Hour).Unix()
			pipe.ZRemRangeByScore(ctx, "notifications:failed", "-inf", fmt.Sprintf("%d", thirtyDaysAgo))
		}
		_, _ = pipe.Exec(ctx)
	}

	return finalResults, nil
}
