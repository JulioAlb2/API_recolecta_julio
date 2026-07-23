package application

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type FailedNotificationRecord struct {
	UserID    string    `json:"user_id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	Type      string    `json:"type"`
	Error     string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

type QueryFailedNotificationsUseCase struct {
	rdb *redis.Client
}

func NewQueryFailedNotificationsUseCase(rdb *redis.Client) *QueryFailedNotificationsUseCase {
	return &QueryFailedNotificationsUseCase{rdb: rdb}
}

func (uc *QueryFailedNotificationsUseCase) Execute(ctx context.Context, start, end time.Time) ([]FailedNotificationRecord, error) {
	minScore := fmt.Sprintf("%d", start.Unix())
	maxScore := fmt.Sprintf("%d", end.Unix())

	vals, err := uc.rdb.ZRangeByScore(ctx, "notifications:failed", &redis.ZRangeBy{
		Min: minScore,
		Max: maxScore,
	}).Result()
	if err != nil {
		return nil, err
	}

	records := make([]FailedNotificationRecord, 0, len(vals))
	for _, val := range vals {
		var r FailedNotificationRecord
		if err := json.Unmarshal([]byte(val), &r); err == nil {
			records = append(records, r)
		}
	}
	return records, nil
}
