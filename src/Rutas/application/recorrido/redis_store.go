package recorrido

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	fieldRutaID           = "ruta_id"
	fieldChoferID         = "chofer_id"
	fieldCamionID         = "camion_id"
	fieldIniciada         = "iniciada"
	fieldPuntoActualIndex = "punto_actual_index"
	fieldPausado          = "pausado"
	fieldStateCode        = "state_code"
	fieldUpdatedAt        = "updated_at"
)

// Estados operativos que pausan el avance del recorrido.
var estadosPausa = map[string]bool{
	"2": true, // Vaciando tolva
	"3": true, // Repostando
	"4": true, // Volviendo a base
}

type Snapshot struct {
	RutaID           int32  `json:"ruta_id"`
	ChoferID         int32  `json:"chofer_id"`
	CamionID         int32  `json:"camion_id"`
	Iniciada         bool   `json:"iniciada"`
	PuntoActualIndex int    `json:"punto_actual_index"`
	Pausado          bool   `json:"pausado"`
	StateCode        string `json:"state_code"`
}

type RedisStore struct {
	rdb *redis.Client
}

func NewRedisStore(rdb *redis.Client) *RedisStore {
	return &RedisStore{rdb: rdb}
}

func camionKey(camionID int32) string {
	return fmt.Sprintf("recorrido:camion:%d", camionID)
}

func choferKey(choferID int32) string {
	return fmt.Sprintf("recorrido:chofer:%d", choferID)
}

func rutaKey(rutaID int32) string {
	return fmt.Sprintf("recorrido:ruta:%d", rutaID)
}

func (s *RedisStore) Iniciar(ctx context.Context, choferID, camionID, rutaID int32) error {
	now := time.Now().Format(time.RFC3339)
	data := map[string]interface{}{
		fieldRutaID:           rutaID,
		fieldChoferID:         choferID,
		fieldCamionID:         camionID,
		fieldIniciada:         "true",
		fieldPuntoActualIndex: "0",
		fieldPausado:          "false",
		fieldStateCode:        "1",
		fieldUpdatedAt:        now,
	}

	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, camionKey(camionID), data)
	pipe.Set(ctx, choferKey(choferID), camionID, 24*time.Hour)
	pipe.Set(ctx, rutaKey(rutaID), camionID, 24*time.Hour)
	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisStore) FinalizarByChofer(ctx context.Context, choferID int32) error {
	camionID, err := s.resolveCamionByChofer(ctx, choferID)
	if err != nil {
		return err
	}
	if camionID == 0 {
		return nil
	}
	return s.FinalizarByCamion(ctx, camionID, choferID)
}

func (s *RedisStore) FinalizarByCamion(ctx context.Context, camionID, choferID int32) error {
	key := camionKey(camionID)
	rutaID, _ := s.readIntField(ctx, key, fieldRutaID)

	pipe := s.rdb.Pipeline()
	pipe.Del(ctx, key)
	if rutaID > 0 {
		pipe.Del(ctx, rutaKey(int32(rutaID)))
	}
	if choferID > 0 {
		pipe.Del(ctx, choferKey(choferID))
	} else if id, err := s.readIntField(ctx, key, fieldChoferID); err == nil && id > 0 {
		pipe.Del(ctx, choferKey(int32(id)))
	}
	_, err := pipe.Exec(ctx)
	return err
}

func (s *RedisStore) AvanzarByChofer(ctx context.Context, choferID int32) error {
	camionID, err := s.resolveCamionByChofer(ctx, choferID)
	if err != nil || camionID == 0 {
		return fmt.Errorf("recorrido activo no encontrado para el chofer")
	}

	key := camionKey(camionID)
	pausado, _ := s.readBoolField(ctx, key, fieldPausado)
	if pausado {
		return fmt.Errorf("recorrido pausado: no se avanza mientras el camión vacía, reposta o vuelve a base")
	}

	idx, _ := s.readIntField(ctx, key, fieldPuntoActualIndex)
	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, key, fieldPuntoActualIndex, idx+1)
	pipe.HSet(ctx, key, fieldUpdatedAt, time.Now().Format(time.RFC3339))
	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisStore) GetActivoByChofer(ctx context.Context, choferID int32) (*Snapshot, error) {
	camionID, err := s.resolveCamionByChofer(ctx, choferID)
	if err != nil || camionID == 0 {
		return nil, nil
	}
	return s.GetByCamion(ctx, camionID)
}

func (s *RedisStore) GetByRuta(ctx context.Context, rutaID int32) (*Snapshot, error) {
	val, err := s.rdb.Get(ctx, rutaKey(rutaID)).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	camionID, err := strconv.ParseInt(val, 10, 32)
	if err != nil {
		return nil, err
	}
	return s.GetByCamion(ctx, int32(camionID))
}

func (s *RedisStore) GetByCamion(ctx context.Context, camionID int32) (*Snapshot, error) {
	key := camionKey(camionID)
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return nil, err
	}

	snap := &Snapshot{CamionID: camionID, Iniciada: true}
	if v, err := s.readIntField(ctx, key, fieldRutaID); err == nil {
		snap.RutaID = int32(v)
	}
	if v, err := s.readIntField(ctx, key, fieldChoferID); err == nil {
		snap.ChoferID = int32(v)
	}
	if v, err := s.readIntField(ctx, key, fieldPuntoActualIndex); err == nil {
		snap.PuntoActualIndex = v
	}
	snap.Pausado, _ = s.readBoolField(ctx, key, fieldPausado)
	snap.StateCode, _ = s.rdb.HGet(ctx, key, fieldStateCode).Result()
	return snap, nil
}

func (s *RedisStore) SyncOperationalState(ctx context.Context, camionID int32, stateCode string) error {
	key := camionKey(camionID)
	exists, err := s.rdb.Exists(ctx, key).Result()
	if err != nil || exists == 0 {
		return nil
	}

	pausado := estadosPausa[stateCode]
	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, key, fieldStateCode, stateCode)
	pipe.HSet(ctx, key, fieldPausado, pausado)
	pipe.HSet(ctx, key, fieldUpdatedAt, time.Now().Format(time.RFC3339))

	if stateCode == "5" {
		if choferID, err := s.readIntField(ctx, key, fieldChoferID); err == nil {
			pipe.Del(ctx, choferKey(int32(choferID)))
		}
		if rutaID, err := s.readIntField(ctx, key, fieldRutaID); err == nil && rutaID > 0 {
			pipe.Del(ctx, rutaKey(int32(rutaID)))
		}
		pipe.Del(ctx, key)
	}

	_, err = pipe.Exec(ctx)
	return err
}

func (s *RedisStore) IsPausedByCamion(ctx context.Context, camionID int32) (bool, error) {
	return s.readBoolField(ctx, camionKey(camionID), fieldPausado)
}

func (s *RedisStore) resolveCamionByChofer(ctx context.Context, choferID int32) (int32, error) {
	val, err := s.rdb.Get(ctx, choferKey(choferID)).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	id, err := strconv.ParseInt(val, 10, 32)
	return int32(id), err
}

func (s *RedisStore) readIntField(ctx context.Context, key, field string) (int, error) {
	val, err := s.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(val)
}

func (s *RedisStore) readBoolField(ctx context.Context, key, field string) (bool, error) {
	val, err := s.rdb.HGet(ctx, key, field).Result()
	if err != nil {
		return false, err
	}
	return val == "true" || val == "1", nil
}
