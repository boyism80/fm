package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqNoticeChanged struct{ gs *GameServer }

func (guildMqNoticeChanged) New(gs *GameServer) *guildMqNoticeChanged {
	return &guildMqNoticeChanged{gs: gs}
}

func (*guildMqNoticeChanged) EventType() string {
	return "notice_changed"
}

func (h *guildMqNoticeChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
			notice := g.Notice
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
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildNoticeChange{
					CharacterID: memberID,
					GuildID:     guildID,
					Notice:      notice,
				})
			}
			return nil
		}).
		Run()
	return nil
}
