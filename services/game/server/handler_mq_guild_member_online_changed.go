package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqMemberOnlineChanged struct{ gs *GameServer }

func (guildMqMemberOnlineChanged) New(gs *GameServer) *guildMqMemberOnlineChanged {
	return &guildMqMemberOnlineChanged{gs: gs}
}

func (*guildMqMemberOnlineChanged) EventType() string {
	return "member_online_changed"
}

func (h *guildMqMemberOnlineChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.guild == nil {
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
	gs.guild.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) {
			return nil, nil
		}, func(interface{}) error {
			g := gs.guild.Get(evt.GuildID)
			if g == nil {
				return nil
			}
			guildID := g.GetGuildId()
			subjectID := extra.CharacterID
			online := extra.Online
			allianceID := g.AllianceID
			for _, m := range g.GetMembers() {
				if m == nil {
					continue
				}
				memberID := m.GetCharacterId()
				if memberID == 0 || memberID == subjectID {
					continue
				}
				if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
					continue
				}
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberOnlineChange{
					CharacterID: memberID,
					GuildID:     guildID,
					SubjectID:   subjectID,
					Online:      online,
				})
			}
			if allianceID > 0 {
				gs.deliverAllianceMemberOnlineChange(allianceID, guildID, subjectID, online)
			}
			return nil
		}).
		Run()
	return nil
}
