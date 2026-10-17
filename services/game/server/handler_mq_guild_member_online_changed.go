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
		CharacterID  uint32 `json:"character_id"`
		Online       bool   `json:"online"`
		ChannelIndex int32  `json:"channel_index"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil
	}
	gs.systems.Get[*game.GuildSystem]().SetMemberChannel(evt.GuildID, extra.CharacterID, extra.ChannelIndex)
	gs.systems.Get[*game.GuildSystem]().BroadcastMemberOnlineChanged(evt.GuildID, extra.CharacterID, extra.Online)
	return nil
}
