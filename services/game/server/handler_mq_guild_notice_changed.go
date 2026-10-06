package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
)

type guildMqNoticeChanged struct{ gs *GameServer }

func (guildMqNoticeChanged) New(gs *GameServer) *guildMqNoticeChanged {
	return &guildMqNoticeChanged{gs: gs}
}

func (*guildMqNoticeChanged) EventType() string {
	return "notice_changed"
}

func (h *guildMqNoticeChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.guild.SyncGuildEventAsync(ctx, evt, gs.guild.BroadcastNoticeChanged)
	return nil
}
