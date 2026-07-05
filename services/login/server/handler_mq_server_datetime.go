package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	amqp "github.com/rabbitmq/amqp091-go"
)

type loginGlobalMqServerDatetime struct{}

func (loginGlobalMqServerDatetime) New(_ *LoginServer) *loginGlobalMqServerDatetime {
	return &loginGlobalMqServerDatetime{}
}

func (*loginGlobalMqServerDatetime) EventType() string {
	return "server_datetime_changed"
}

func (*loginGlobalMqServerDatetime) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	var payload struct {
		Reset    bool   `json:"reset"`
		Datetime string `json:"datetime"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if err := clock.ApplyDateTime(payload.Reset, payload.Datetime); err != nil {
		log.Printf("login mq server_datetime_changed: %v", err)
	}
	return nil
}
