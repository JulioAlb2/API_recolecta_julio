package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type RedisNotificationRuleRepository struct {
	rdb *redis.Client
}

func NewRedisNotificationRuleRepository(rdb *redis.Client) *RedisNotificationRuleRepository {
	return &RedisNotificationRuleRepository{rdb: rdb}
}

func ruleKey(tenantID int, code string) string {
	return fmt.Sprintf("rules:state:%d:%s", tenantID, code)
}

func (r *RedisNotificationRuleRepository) Save(ctx context.Context, tenantID int, rule domain.NotificationRule) error {
	rule.TenantID = tenantID
	data, err := json.Marshal(rule)
	if err != nil {
		return err
	}
	return r.rdb.Set(ctx, ruleKey(tenantID, rule.StateCode), data, 0).Err()
}

func (r *RedisNotificationRuleRepository) GetByStateCode(ctx context.Context, tenantID int, code string) (*domain.NotificationRule, error) {
	val, err := r.rdb.Get(ctx, ruleKey(tenantID, code)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("rule not found: %s", code)
	}
	if err != nil {
		return nil, err
	}
	var rule domain.NotificationRule
	if err := json.Unmarshal([]byte(val), &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

func (r *RedisNotificationRuleRepository) List(ctx context.Context, tenantID int) ([]domain.NotificationRule, error) {
	keys, err := r.rdb.Keys(ctx, fmt.Sprintf("rules:state:%d:*", tenantID)).Result()
	if err != nil {
		return nil, err
	}
	rules := make([]domain.NotificationRule, 0, len(keys))
	for _, key := range keys {
		val, err := r.rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}
		var rule domain.NotificationRule
		if err := json.Unmarshal([]byte(val), &rule); err != nil {
			continue
		}
		rules = append(rules, rule)
	}
	return rules, nil
}

func (r *RedisNotificationRuleRepository) Delete(ctx context.Context, tenantID int, code string) error {
	return r.rdb.Del(ctx, ruleKey(tenantID, code)).Err()
}
