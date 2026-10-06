package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type guildMqCapacityChanged struct{ gs *GameServer }

func (guildMqCapacityChanged) New(gs *GameServer) *guildMqCapacityChanged {
	return &guildMqCapacityChanged{gs: gs}
}

func (*guildMqCapacityChanged) EventType() string {
	return "capacity_changed"
}

func (h *guildMqCapacityChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	var payload struct {
		GPAmount int32 `json:"gp_amount"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("guild consumer: capacity_changed guild_id=%d invalid payload: %v", evt.GuildID, err)
		return nil
	}
	gs.guild.SyncGuildEventAsync(ctx, evt, func(guildID uint32) {
		gs.guild.BroadcastCapacityChanged(guildID)
		gs.guild.BroadcastGPChanged(guildID, payload.GPAmount)
	})
	return nil
}
