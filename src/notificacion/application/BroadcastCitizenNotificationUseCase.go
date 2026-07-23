package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/core"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type BroadcastCitizenNotificationUseCase struct {
	sender domain.PushNotificationSender
	rdb    *redis.Client
}

func NewBroadcastCitizenNotificationUseCase(sender domain.PushNotificationSender, rdb *redis.Client) *BroadcastCitizenNotificationUseCase {
	return &BroadcastCitizenNotificationUseCase{sender: sender, rdb: rdb}
}

// ExecuteBroadcast sends a push notification to all citizens
func (uc *BroadcastCitizenNotificationUseCase) ExecuteBroadcast(ctx context.Context, title, body string) (map[string]domain.SendResult, error) {
	db := core.GetBD()
	rows, err := db.Query(ctx, "SELECT id FROM ciudadano")
	if err != nil {
		return nil, fmt.Errorf("error querying citizens: %w", err)
	}
	defer rows.Close()

	var citizenIDs []string
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err == nil {
			citizenIDs = append(citizenIDs, fmt.Sprintf("%d", id))
		}
	}

	notification := &domain.PushNotification{
		Title: title,
		Body:  body,
		Type:  "BROADCAST",
	}

	return uc.sendAndLog(ctx, citizenIDs, notification)
}

// ExecuteRouteBroadcast sends a push notification to all citizens residing in the colony of a route
func (uc *BroadcastCitizenNotificationUseCase) ExecuteRouteBroadcast(ctx context.Context, routeID int32, title, body string) (map[string]domain.SendResult, error) {
	db := core.GetBD()

	query := `
		SELECT DISTINCT d.ciudadano_id 
		FROM domicilio d
		JOIN ruta r ON r.colonia_id = d.colonia_id
		WHERE r.id = $1 AND d.deleted_at IS NULL AND r.deleted_at IS NULL
	`
	rows, err := db.Query(ctx, query, routeID)
	if err != nil {
		return nil, fmt.Errorf("error querying citizens by route: %w", err)
	}
	defer rows.Close()

	var citizenIDs []string
	for rows.Next() {
		var cid int32
		if err := rows.Scan(&cid); err == nil {
			citizenIDs = append(citizenIDs, fmt.Sprintf("%d", cid))
		}
	}

	if len(citizenIDs) == 0 {
		return nil, fmt.Errorf("no citizens found for route %d", routeID)
	}

	notification := &domain.PushNotification{
		Title: title,
		Body:  body,
		Type:  "ROUTE_BROADCAST",
		Data: map[string]string{
			"route_id": fmt.Sprintf("%d", routeID),
		},
	}

	return uc.sendAndLog(ctx, citizenIDs, notification)
}

// ExecutePointBroadcast sends a push notification to all citizens near a recollection point
func (uc *BroadcastCitizenNotificationUseCase) ExecutePointBroadcast(ctx context.Context, pointID string, title, body string) (map[string]domain.SendResult, error) {
	ptKey := fmt.Sprintf("point:%s", pointID)
	latStr, err := uc.rdb.HGet(ctx, ptKey, "lat").Result()
	lonStr, _ := uc.rdb.HGet(ctx, ptKey, "lon").Result()

	if err != nil || latStr == "" || lonStr == "" {
		return nil, fmt.Errorf("coordinates not found in Redis for point %s", pointID)
	}

	ptLat, _ := strconv.ParseFloat(latStr, 64)
	ptLon, _ := strconv.ParseFloat(lonStr, 64)

	// Fetch search radius from rules (default 200m for point broadcast)
	radius := 200.0

	// Search geofenced addresses
	doms, err := uc.rdb.GeoRadius(ctx, "domicilios:geo", ptLon, ptLat, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   "m",
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("error querying geolocation from Redis: %w", err)
	}

	var citizenIDs []string
	for _, d := range doms {
		domKey := fmt.Sprintf("domicilio:%s", d.Name)
		cid, err := uc.rdb.HGet(ctx, domKey, "ciudadano_id").Result()
		if err == nil && cid != "" {
			citizenIDs = append(citizenIDs, cid)
		}
	}

	if len(citizenIDs) == 0 {
		return nil, fmt.Errorf("no citizens found near point %s within %vm", pointID, radius)
	}

	notification := &domain.PushNotification{
		Title: title,
		Body:  body,
		Type:  "POINT_BROADCAST",
		Data: map[string]string{
			"point_id": pointID,
		},
	}

	return uc.sendAndLog(ctx, citizenIDs, notification)
}

// sendAndLog retrieves FCM tokens, sends the push alerts, and audits delivery
func (uc *BroadcastCitizenNotificationUseCase) sendAndLog(ctx context.Context, citizenIDs []string, notification *domain.PushNotification) (map[string]domain.SendResult, error) {
	if len(citizenIDs) == 0 {
		return nil, fmt.Errorf("no citizens to notify")
	}

	userTokens := make(map[string]string)
	for _, cid := range citizenIDs {
		userKey := fmt.Sprintf("user:%s", cid)
		token, err := uc.rdb.HGet(ctx, userKey, "fcm_token").Result()
		if err == nil && token != "" {
			userTokens[cid] = token
		}
	}

	if len(userTokens) == 0 {
		return nil, fmt.Errorf("no active FCM tokens found for target citizens")
	}

	results, err := uc.sender.Send(ctx, userTokens, notification)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	nowUnix := float64(now.Unix())
	nowStr := now.Format(time.RFC3339)

	for citizenID, res := range results {
		// Individual Inbox
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

		pipe := uc.rdb.Pipeline()
		pipe.LPush(ctx, inboxKey, string(inboxData))
		pipe.LTrim(ctx, inboxKey, 0, 49)
		pipe.Expire(ctx, inboxKey, 30*24*time.Hour)

		// Failed ZSET
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

	return results, nil
}
