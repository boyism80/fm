package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqRankTitlesChanged struct{ gs *GameServer }

func (guildMqRankTitlesChanged) New(gs *GameServer) *guildMqRankTitlesChanged {
	return &guildMqRankTitlesChanged{gs: gs}
}

func (*guildMqRankTitlesChanged) EventType() string {
	return "rank_titles_changed"
}

func (h *guildMqRankTitlesChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
			rankTitles := g.RankTitles
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
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildRankTitleChange{
					CharacterID: memberID,
					GuildID:     guildID,
					RankTitles:  rankTitles,
				})
			}
			return nil
		}).
		Run()
	return nil
}
