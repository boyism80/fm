package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
)

type globalMqServerDatetime struct{}

func (globalMqServerDatetime) New(_ *GameServer) *globalMqServerDatetime {
	return &globalMqServerDatetime{}
}

func (*globalMqServerDatetime) EventType() string {
	return "server_datetime_changed"
}

func (*globalMqServerDatetime) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		Reset    bool   `json:"reset"`
		Datetime string `json:"datetime"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	return clock.SyncDateTime(payload.Reset, payload.Datetime)
}
