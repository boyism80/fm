package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqGPChanged struct{ gs *GameServer }

func (guildMqGPChanged) New(gs *GameServer) *guildMqGPChanged {
	return &guildMqGPChanged{gs: gs}
}

func (*guildMqGPChanged) EventType() string {
	return "gp_changed"
}

func (h *guildMqGPChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	var payload struct {
		Amount int32 `json:"amount"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("guild consumer: gp_changed guild_id=%d invalid payload: %v", evt.GuildID, err)
		return nil
	}
	gs.guild.SyncGuildEventAsync(ctx, evt, func(guildID uint32) {
		gs.guild.BroadcastGPChanged(guildID, payload.Amount)
	})
	return nil
}
