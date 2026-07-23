package application

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type InboxRecord struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	Type      string `json:"type"`
	SentAt    string `json:"sent_at"`
	Delivered bool   `json:"delivered"`
	Error     string `json:"error,omitempty"`
}

type GetCitizenInboxUseCase struct {
	rdb *redis.Client
}

func NewGetCitizenInboxUseCase(rdb *redis.Client) *GetCitizenInboxUseCase {
	return &GetCitizenInboxUseCase{rdb: rdb}
}

func (uc *GetCitizenInboxUseCase) Execute(ctx context.Context, citizenID string) ([]InboxRecord, error) {
	key := fmt.Sprintf("citizen:notifications:%s", citizenID)
	vals, err := uc.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	inbox := make([]InboxRecord, 0, len(vals))
	for _, val := range vals {
		var item InboxRecord
		if err := json.Unmarshal([]byte(val), &item); err == nil {
			inbox = append(inbox, item)
		}
	}
	return inbox, nil
}
