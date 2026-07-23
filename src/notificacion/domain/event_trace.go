package domain

import (
	"context"
	"time"
)

type EventVersion string

const EventVersionV1 EventVersion = "v1"

type StateCode string

const StateArrival StateCode = "ARRIVAL"

type Action string

const ActionNotifyAdminOnly Action = "NOTIFY_ADMIN_ONLY"

type TruckStateEvent struct {
	EventID      string
	EventType    string
	EventVersion EventVersion
	TruckID      int32
	StateCode    StateCode
	OccurredAt   time.Time
	Lat          float64
	Lon          float64
}

type EventTraceRecord struct {
	EventID            string       `json:"event_id"`
	EventHash          string       `json:"event_hash"`
	EventType          string       `json:"event_type"`
	EventVersion       EventVersion `json:"event_version"`
	TruckID            int32        `json:"truck_id"`
	StateCode          StateCode    `json:"state_code"`
	ResolvedAction     Action       `json:"resolved_action"`
	AdminNotified      bool         `json:"admin_notified"`
	CitizenFanoutCount int          `json:"citizen_fanout_count"`
	Result             string       `json:"result"`
	CreatedAt          time.Time    `json:"created_at"`
}

type IEventTraceRepository interface {
	TryAcquireDeduplication(ctx context.Context, hash string, event *TruckStateEvent) (bool, error)
	SaveTrace(ctx context.Context, trace *EventTraceRecord) error
	GetByEventID(ctx context.Context, eventID string) (*EventTraceRecord, error)
	ListByTruckID(ctx context.Context, truckID int32, limit int) ([]EventTraceRecord, error)
}
