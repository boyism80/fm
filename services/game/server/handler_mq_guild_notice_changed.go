package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
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
	gs.systems.Get[*domain.GuildSystem]().SyncGuildEventAsync(ctx, evt, gs.systems.Get[*domain.GuildSystem]().BroadcastNoticeChanged)
	return nil
}
