package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type guildMqMemberOnlineChanged struct{ gs *GameServer }

func (guildMqMemberOnlineChanged) New(gs *GameServer) *guildMqMemberOnlineChanged {
	return &guildMqMemberOnlineChanged{gs: gs}
}

func (*guildMqMemberOnlineChanged) EventType() string {
	return "member_online_changed"
}

func (h *guildMqMemberOnlineChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
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
		Online      bool   `json:"online"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
		return nil
	}
	subjectID := extra.CharacterID
	online := extra.Online
	gs.systems.Get[*game.GuildSystem]().SyncGuildEventAsync(ctx, evt, func(guildID uint32) {
		gs.systems.Get[*game.GuildSystem]().BroadcastMemberOnlineChanged(guildID, subjectID, online)
	})
	return nil
}
