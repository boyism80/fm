package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqMemberRankChanged struct{ gs *GameServer }

func (guildMqMemberRankChanged) New(gs *GameServer) *guildMqMemberRankChanged {
	return &guildMqMemberRankChanged{gs: gs}
}

func (*guildMqMemberRankChanged) EventType() string {
	return "member_rank_changed"
}

func (h *guildMqMemberRankChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
			var extra struct {
				CharacterID uint32 `json:"character_id"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
				return nil
			}
			g := gs.guild.Get(evt.GuildID)
			if g == nil {
				return nil
			}
			targetCharacterID := extra.CharacterID
			targetRank := uint8(0)
			found := false
			for _, m := range g.GetMembers() {
				if m == nil || m.GetCharacterId() != targetCharacterID {
					continue
				}
				targetRank = entity.GuildMemberRankWire(m.GetRank())
				found = true
				break
			}
			if !found {
				return nil
			}
			guildID := g.GetGuildId()
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
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberRankChange{
					CharacterID: memberID,
					GuildID:     guildID,
					TargetID:    targetCharacterID,
					GuildRank:   targetRank,
				})
			}
			return nil
		}).
		Run()
	return nil
}
