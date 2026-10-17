package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type guildMqMemberJoined struct{ gs *GameServer }

func (guildMqMemberJoined) New(gs *GameServer) *guildMqMemberJoined {
	return &guildMqMemberJoined{gs: gs}
}

func (*guildMqMemberJoined) EventType() string {
	return "member_joined"
}

func (h *guildMqMemberJoined) Handle(ctx actor.Context, raw json.RawMessage) error {
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
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil
	}
	if !gs.systems.Get[*game.GuildSystem]().Update(decodeGuildSnapshot(evt.GuildPB)) {
		return nil
	}
	gs.systems.Get[*game.GuildSystem]().BroadcastMemberJoined(evt.GuildID, extra.CharacterID)
	return nil
}
