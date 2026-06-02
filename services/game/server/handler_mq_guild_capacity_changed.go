package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqCapacityChanged struct{ gs *GameServer }

func (guildMqCapacityChanged) New(gs *GameServer) *guildMqCapacityChanged {
	return &guildMqCapacityChanged{gs: gs}
}

func (*guildMqCapacityChanged) EventType() string {
	return "capacity_changed"
}

func (h *guildMqCapacityChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
			if g == nil {
				return nil
			}
			guildID := g.GetGuildId()
			capacity := uint8(g.Capacity)
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
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildCapacityChange{
					CharacterID: memberID,
					GuildID:     guildID,
					Capacity:    capacity,
				})
			}
			return nil
		}).
		Run()
	return nil
}
