package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqRankTitlesChanged struct{ gs *GameServer }

func (guildMqRankTitlesChanged) New(gs *GameServer) *guildMqRankTitlesChanged {
	return &guildMqRankTitlesChanged{gs: gs}
}

func (*guildMqRankTitlesChanged) EventType() string {
	return "rank_titles_changed"
}

func (h *guildMqRankTitlesChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.guild.ApplyEventAsync(ctx, evt, gs.guild.BroadcastRankTitlesChanged).Run()
	return nil
}
