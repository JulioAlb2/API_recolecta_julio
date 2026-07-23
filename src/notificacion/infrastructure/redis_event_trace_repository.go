package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type RedisEventTraceRepository struct {
	rdb *redis.Client
}

func NewRedisEventTraceRepository(rdb *redis.Client) *RedisEventTraceRepository {
	return &RedisEventTraceRepository{rdb: rdb}
}

func (r *RedisEventTraceRepository) TryAcquireDeduplication(ctx context.Context, hash string, _ *domain.TruckStateEvent) (bool, error) {
	key := fmt.Sprintf("event:dedup:%s", hash)
	return r.rdb.SetNX(ctx, key, "1", 24*time.Hour).Result()
}

func (r *RedisEventTraceRepository) SaveTrace(ctx context.Context, trace *domain.EventTraceRecord) error {
	data, err := json.Marshal(trace)
	if err != nil {
		return err
	}
	key := fmt.Sprintf("event:trace:%s", trace.EventID)
	if err := r.rdb.Set(ctx, key, data, 30*24*time.Hour).Err(); err != nil {
		return err
	}
	truckKey := fmt.Sprintf("truck:events:%d", trace.TruckID)
	if err := r.rdb.LPush(ctx, truckKey, trace.EventID).Err(); err != nil {
		return err
	}
	return r.rdb.Expire(ctx, truckKey, 30*24*time.Hour).Err()
}

func (r *RedisEventTraceRepository) GetByEventID(ctx context.Context, eventID string) (*domain.EventTraceRecord, error) {
	key := fmt.Sprintf("event:trace:%s", eventID)
	val, err := r.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("trace not found: %s", eventID)
	}
	if err != nil {
		return nil, err
	}
	var trace domain.EventTraceRecord
	if err := json.Unmarshal([]byte(val), &trace); err != nil {
		return nil, err
	}
	return &trace, nil
}

func (r *RedisEventTraceRepository) ListByTruckID(ctx context.Context, truckID int32, limit int) ([]domain.EventTraceRecord, error) {
	truckKey := fmt.Sprintf("truck:events:%d", truckID)
	eventIDs, err := r.rdb.LRange(ctx, truckKey, 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	traces := make([]domain.EventTraceRecord, 0, len(eventIDs))
	for _, eventID := range eventIDs {
		trace, err := r.GetByEventID(ctx, eventID)
		if err != nil {
			continue
		}
		traces = append(traces, *trace)
	}
	return traces, nil
}
