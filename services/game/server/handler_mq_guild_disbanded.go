package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqDisbanded struct{ gs *GameServer }

func (guildMqDisbanded) New(gs *GameServer) *guildMqDisbanded {
	return &guildMqDisbanded{gs: gs}
}

func (*guildMqDisbanded) EventType() string {
	return "disbanded"
}

func (h *guildMqDisbanded) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevGuild := gs.guild.Get(evt.GuildID)
	var extra struct {
		MemberCharacterIDs []uint32 `json:"member_character_ids"`
	}
	_ = json.Unmarshal(raw, &extra)
	memberIDs := extra.MemberCharacterIDs
	gs.guild.ApplyEventAsync(ctx, evt, func(uint32) {
		gs.guild.BroadcastDisbanded(prevGuild, memberIDs)
	}).Run()
	return nil
}
