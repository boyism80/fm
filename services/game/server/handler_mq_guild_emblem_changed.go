package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqEmblemChanged struct{ gs *GameServer }

func (guildMqEmblemChanged) New(gs *GameServer) *guildMqEmblemChanged {
	return &guildMqEmblemChanged{gs: gs}
}

func (*guildMqEmblemChanged) EventType() string {
	return "emblem_changed"
}

func (h *guildMqEmblemChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.guild == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	gs.guild.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) {
			return nil, nil
		}, func(interface{}) error {
			g := gs.guild.Get(evt.GuildID)
			if g == nil || g.Logo == nil {
				return nil
			}
			guildID := g.GetGuildId()
			logo := g.Logo
			for _, m := range g.GetMembers() {
				if m == nil {
					continue
				}
				memberID := m.GetCharacterId()
				if memberID == 0 {
					continue
				}
				if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
					continue
				}
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildEmblemChange{
					CharacterID: memberID,
					GuildID:     guildID,
					LogoBG:      uint16(logo.LogoBG),
					LogoBGColor: uint8(logo.LogoBGColor),
					Logo:        uint16(logo.Logo),
					LogoColor:   uint8(logo.LogoColor),
				})
			}
			return nil
		}).
		Run()
	return nil
}
