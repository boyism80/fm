package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqMemberRankChanged struct{ gs *GameServer }

func (guildMqMemberRankChanged) New(gs *GameServer) *guildMqMemberRankChanged {
	return &guildMqMemberRankChanged{gs: gs}
}

func (*guildMqMemberRankChanged) EventType() string {
	return "member_rank_changed"
}

func (h *guildMqMemberRankChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	var extra struct {
		CharacterID uint32 `json:"character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
		return nil
	}
	targetID := extra.CharacterID
	gs.guild.ApplyEventAsync(ctx, evt, func(guildID uint32) {
		gs.guild.BroadcastMemberRankChanged(guildID, targetID)
	})
	return nil
}
