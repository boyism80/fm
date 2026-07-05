package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	amqp "github.com/rabbitmq/amqp091-go"
)

type globalMqServerDatetime struct{}

func (globalMqServerDatetime) New(_ *GameServer) *globalMqServerDatetime {
	return &globalMqServerDatetime{}
}

func (*globalMqServerDatetime) EventType() string {
	return "server_datetime_changed"
}

func (*globalMqServerDatetime) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	var payload struct {
		Reset    bool   `json:"reset"`
		Datetime string `json:"datetime"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	return clock.ApplyDateTime(payload.Reset, payload.Datetime)
}
