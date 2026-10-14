package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type guildMqMemberRankChanged struct{ gs *GameServer }

func (guildMqMemberRankChanged) New(gs *GameServer) *guildMqMemberRankChanged {
	return &guildMqMemberRankChanged{gs: gs}
}

func (*guildMqMemberRankChanged) EventType() string {
	return "member_rank_changed"
}

func (h *guildMqMemberRankChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
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
	gs.systems.Get[*domain.GuildSystem]().SyncGuildEventAsync(ctx, evt, func(guildID uint32) {
		gs.systems.Get[*domain.GuildSystem]().BroadcastMemberRankChanged(guildID, targetID)
	})
	return nil
}
