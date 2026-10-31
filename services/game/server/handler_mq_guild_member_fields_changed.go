package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type guildMqMemberFieldsChanged struct{ gs *GameServer }

func (guildMqMemberFieldsChanged) New(gs *GameServer) *guildMqMemberFieldsChanged {
	return &guildMqMemberFieldsChanged{gs: gs}
}

func (*guildMqMemberFieldsChanged) EventType() string {
	return "member_fields_changed"
}

func (h *guildMqMemberFieldsChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	var extra struct {
		CharacterID uint32 `json:"character_id"`
		Level       uint32 `json:"level"`
		ClassID     uint32 `json:"class_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil
	}
	h.gs.systems.Get[*game.GuildSystem]().UpdateMemberFields(evt.GuildID, extra.CharacterID, extra.Level, extra.ClassID)
	return nil
}
