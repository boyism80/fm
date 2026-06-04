package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqEmblemChanged struct{ gs *GameServer }

func (guildMqEmblemChanged) New(gs *GameServer) *guildMqEmblemChanged {
	return &guildMqEmblemChanged{gs: gs}
}

func (*guildMqEmblemChanged) EventType() string {
	return "emblem_changed"
}

func (h *guildMqEmblemChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.guild.ApplyEventAsync(ctx, evt, gs.guild.BroadcastEmblemChanged).Run()
	return nil
}
