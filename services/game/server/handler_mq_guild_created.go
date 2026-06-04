package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqCreated struct{ gs *GameServer }

func (guildMqCreated) New(gs *GameServer) *guildMqCreated {
	return &guildMqCreated{gs: gs}
}

func (*guildMqCreated) EventType() string {
	return "created"
}

func (h *guildMqCreated) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.guild.UpdateAsync(ctx, evt).Run()
	return nil
}
