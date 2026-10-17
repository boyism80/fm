package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type guildMqMemberLeft struct{ gs *GameServer }

func (guildMqMemberLeft) New(gs *GameServer) *guildMqMemberLeft {
	return &guildMqMemberLeft{gs: gs}
}

func (*guildMqMemberLeft) EventType() string {
	return "member_left"
}

func (h *guildMqMemberLeft) Handle(ctx actor.Context, raw json.RawMessage) error {
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
		Expelled    bool   `json:"expelled"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil
	}
	prevGuild := gs.systems.Get[*game.GuildSystem]().Get(evt.GuildID)
	if !gs.systems.Get[*game.GuildSystem]().Update(decodeGuildSnapshot(evt.GuildPB)) {
		return nil
	}
	gs.systems.Get[*game.GuildSystem]().BroadcastMemberLeft(prevGuild, extra.CharacterID, extra.Expelled)
	return nil
}
