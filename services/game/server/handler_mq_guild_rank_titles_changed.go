package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type guildMqRankTitlesChanged struct{ gs *GameServer }

func (guildMqRankTitlesChanged) New(gs *GameServer) *guildMqRankTitlesChanged {
	return &guildMqRankTitlesChanged{gs: gs}
}

func (*guildMqRankTitlesChanged) EventType() string {
	return "rank_titles_changed"
}

func (h *guildMqRankTitlesChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	if !gs.systems.Get[*game.GuildSystem]().Update(decodeGuildSnapshot(evt.GuildPB)) {
		return nil
	}
	gs.systems.Get[*game.GuildSystem]().BroadcastRankTitlesChanged(evt.GuildID)
	return nil
}
