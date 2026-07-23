package application_ciudadano

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/vicpoo/API_recolecta/src/core"
)

type UpdateFCMTokenInput struct {
	FCMToken string `json:"fcm_token"`
}

type UpdateFCMToken struct{}

func NewUpdateFCMToken() *UpdateFCMToken {
	return &UpdateFCMToken{}
}

func (uc *UpdateFCMToken) Execute(ctx context.Context, citizenID int, in UpdateFCMTokenInput) error {
	token := strings.TrimSpace(in.FCMToken)
	if token == "" {
		return errors.New("fcm_token es requerido")
	}

	rdb, err := core.ConnectRedis()
	if err != nil {
		return err
	}

	legacyKey := fmt.Sprintf("fcm:ciudadano:%d", citizenID)
	userKey := fmt.Sprintf("user:%d", citizenID)

	pipe := rdb.TxPipeline()
	pipe.Set(ctx, legacyKey, token, 0)
	pipe.HSet(ctx, userKey,
		"fcm_token", token,
		"fcm_status", "valid",
		"updated_at", time.Now().UTC().Format(time.RFC3339),
	)

	_, err = pipe.Exec(ctx)
	return err
}
