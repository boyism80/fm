package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqCapacityChanged struct{ gs *GameServer }

func (guildMqCapacityChanged) New(gs *GameServer) *guildMqCapacityChanged {
	return &guildMqCapacityChanged{gs: gs}
}

func (*guildMqCapacityChanged) EventType() string {
	return "capacity_changed"
}

func (h *guildMqCapacityChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.guild.ApplyEventAsync(ctx, evt, gs.guild.BroadcastCapacityChanged)
	return nil
}
