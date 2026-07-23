package infrastructure

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisNotificationRepository struct {
	rdb *redis.Client
}

func NewRedisNotificationRepository(rdb *redis.Client) *RedisNotificationRepository {
	return &RedisNotificationRepository{rdb: rdb}
}

func (r *RedisNotificationRepository) GetTokensByUserIDs(ctx context.Context, userIDs []string) (map[string]string, error) {
	result := make(map[string]string)
	for _, uid := range userIDs {
		// Intenta obtener de user:<uid> (campo fcm_token)
		userKey := "user:" + uid
		val, err := r.rdb.HGet(ctx, userKey, "fcm_token").Result()
		if err == redis.Nil || val == "" {
			// Alternativa: intenta obtener de fcm:ciudadano:<uid> (valor string)
			legacyKey := "fcm:ciudadano:" + uid
			val, err = r.rdb.Get(ctx, legacyKey).Result()
			if err == redis.Nil {
				continue
			}
		}
		if err != nil {
			return nil, err
		}
		if val != "" {
			result[uid] = val
		}
	}
	return result, nil
}
