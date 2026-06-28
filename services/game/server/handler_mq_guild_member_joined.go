package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqMemberJoined struct{ gs *GameServer }

func (guildMqMemberJoined) New(gs *GameServer) *guildMqMemberJoined {
	return &guildMqMemberJoined{gs: gs}
}

func (*guildMqMemberJoined) EventType() string {
	return "member_joined"
}

func (h *guildMqMemberJoined) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
	joinerID := extra.CharacterID
	gs.guild.ApplyEventAsync(ctx, evt, func(guildID uint32) {
		gs.guild.BroadcastMemberJoined(guildID, joinerID)
	})
	return nil
}
