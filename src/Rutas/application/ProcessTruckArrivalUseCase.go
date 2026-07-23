package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/vicpoo/API_recolecta/src/Rutas/application/recorrido"
	"github.com/vicpoo/API_recolecta/src/notificacion/domain"
)

type ProcessTruckArrivalUseCase struct {
	rdb            *redis.Client
	rulesRepo      domain.INotificationRuleRepository
	sender         domain.PushNotificationSender
	recorridoStore *recorrido.RedisStore
}

func NewProcessTruckArrivalUseCase(rdb *redis.Client, rulesRepo domain.INotificationRuleRepository, sender domain.PushNotificationSender) *ProcessTruckArrivalUseCase {
	return &ProcessTruckArrivalUseCase{
		rdb:            rdb,
		rulesRepo:      rulesRepo,
		sender:         sender,
		recorridoStore: recorrido.NewRedisStore(rdb),
	}
}

func (uc *ProcessTruckArrivalUseCase) Execute(ctx context.Context, tenantID int, truckID int32, pointID string) error {
	pausado, err := uc.recorridoStore.IsPausedByCamion(ctx, truckID)
	if err == nil && pausado {
		return nil
	}
	key := fmt.Sprintf("truck:%d", truckID)
	nowStr := time.Now().Format(time.RFC3339)

	// 1. Marcar estado como ARRIVAL y registrar last_visited_point en Redis
	pipe := uc.rdb.Pipeline()
	pipe.HSet(ctx, key, "state", "ARRIVAL")
	pipe.HSet(ctx, key, "last_visited_point", pointID)
	pipe.HSet(ctx, key, "updated_at", nowStr)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}

	// 2. Obtener coordenadas de la parada desde Redis
	ptKey := fmt.Sprintf("point:%s", pointID)
	latStr, err := uc.rdb.HGet(ctx, ptKey, "lat").Result()
	lonStr, _ := uc.rdb.HGet(ctx, ptKey, "lon").Result()
	ptName, _ := uc.rdb.HGet(ctx, ptKey, "name").Result()

	if err != nil || latStr == "" || lonStr == "" {
		return fmt.Errorf("coordenadas no encontradas para el punto %s", pointID)
	}

	ptLat, _ := strconv.ParseFloat(latStr, 64)
	ptLon, _ := strconv.ParseFloat(lonStr, 64)

	// 3. Obtener el radio de búsqueda (por defecto 100m) de las reglas de notificaciones
	radius := 100.0
	rule, err := uc.rulesRepo.GetByStateCode(ctx, tenantID, "ARRIVAL")
	if err == nil && rule.Enabled && rule.Threshold > 0 {
		radius = float64(rule.Threshold)
	}

	// 4. Buscar ciudadanos georreferenciados en el radio de la parada
	doms, err := uc.rdb.GeoRadius(ctx, "domicilios:geo", ptLon, ptLat, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   "m",
	}).Result()

	if err != nil || len(doms) == 0 {
		return nil
	}

	todayStr := time.Now().Format("2006-01-02")
	userTokens := make(map[string]string)

	for _, d := range doms {
		domKey := fmt.Sprintf("domicilio:%s", d.Name)
		citizenID, _ := uc.rdb.HGet(ctx, domKey, "ciudadano_id").Result()
		if citizenID == "" {
			continue
		}

		// Evitar notificaciones duplicadas en el mismo punto el mismo día
		dupKey := fmt.Sprintf("notification:sent:%s:%d:%s", citizenID, truckID, todayStr)
		alreadySent, _ := uc.rdb.SIsMember(ctx, dupKey, fmt.Sprintf("arrival:%s", pointID)).Result()
		if alreadySent {
			continue
		}

		userKey := fmt.Sprintf("user:%s", citizenID)
		token, _ := uc.rdb.HGet(ctx, userKey, "fcm_token").Result()
		if token != "" {
			userTokens[citizenID] = token
			uc.rdb.SAdd(ctx, dupKey, fmt.Sprintf("arrival:%s", pointID))
			uc.rdb.Expire(ctx, dupKey, 24*time.Hour)
		}
	}

	// 5. Enviar push notification al grupo del punto actual
	if len(userTokens) > 0 {
		notification := &domain.PushNotification{
			Title: "¡Camión de basura en tu punto!",
			Body:  fmt.Sprintf("El camión ha llegado a la parada: %s. Por favor, saca tus bolsas de basura.", ptName),
			Type:  "ARRIVAL",
			Data: map[string]string{
				"truck_id": fmt.Sprintf("%d", truckID),
				"point_id": pointID,
			},
		}

		results, err := uc.sender.Send(ctx, userTokens, notification)
		if err == nil {
			now := time.Now()
			nowUnix := float64(now.Unix())
			nowStr := now.Format(time.RFC3339)

			for citizenID, res := range results {
				// 1. Registrar en el inbox individual
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

				// 2. Si falló, registrar globalmente en ZSET para auditoría por fecha
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
		}
	}

	return nil
}
