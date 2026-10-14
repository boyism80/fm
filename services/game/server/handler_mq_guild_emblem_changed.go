package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type guildMqEmblemChanged struct{ gs *GameServer }

func (guildMqEmblemChanged) New(gs *GameServer) *guildMqEmblemChanged {
	return &guildMqEmblemChanged{gs: gs}
}

func (*guildMqEmblemChanged) EventType() string {
	return "emblem_changed"
}

func (h *guildMqEmblemChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.systems.Get[*domain.GuildSystem]().SyncGuildEventAsync(ctx, evt, gs.systems.Get[*domain.GuildSystem]().BroadcastEmblemChanged)
	return nil
}
