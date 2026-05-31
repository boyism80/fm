package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqMemberLeft struct{ gs *GameServer }

func (guildMqMemberLeft) New(gs *GameServer) *guildMqMemberLeft {
	return &guildMqMemberLeft{gs: gs}
}

func (*guildMqMemberLeft) EventType() string {
	return "member_left"
}

func (h *guildMqMemberLeft) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.guild == nil {
		return nil
	}
	evt, ok := decodeGuildEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevGuild := gs.guild.Get(evt.GuildID)
	gs.guild.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) {
			return nil, nil
		}, func(interface{}) error {
			var extra struct {
				CharacterID uint32 `json:"character_id"`
				Expelled    bool   `json:"expelled"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
				return nil
			}
			if prevGuild == nil {
				return nil
			}
			leftCharacterID := extra.CharacterID
			leftName := ""
			for _, m := range prevGuild.GetMembers() {
				if m == nil || m.GetCharacterId() != leftCharacterID {
					continue
				}
				leftName = m.GetCharacterName()
				break
			}
			if leftName == "" {
				return nil
			}
			guildID := prevGuild.GetGuildId()
			for _, m := range prevGuild.GetMembers() {
				if m == nil {
					continue
				}
				memberID := m.GetCharacterId()
				if memberID == 0 {
					continue
				}
				if memberID == leftCharacterID {
					if extra.Expelled {
						gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildExpelSelf{
							CharacterID: memberID,
							GuildID:     guildID,
						})
					} else {
						gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildLeaveSelf{
							CharacterID: memberID,
						})
					}
					continue
				}
				if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
					continue
				}
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildMemberLeft{
					CharacterID: memberID,
					GuildID:     guildID,
					TargetID:    leftCharacterID,
					TargetName:  leftName,
					WasExpelled: extra.Expelled,
				})
			}
			return nil
		}).
		Run()
	return nil
}
