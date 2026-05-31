package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/dto"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqMemberJoined struct{ gs *GameServer }

func (guildMqMemberJoined) New(gs *GameServer) *guildMqMemberJoined {
	return &guildMqMemberJoined{gs: gs}
}

func (*guildMqMemberJoined) EventType() string {
	return "member_joined"
}

func (h *guildMqMemberJoined) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
			guildID := evt.GuildID
			joinerCharacterID := extra.CharacterID
			if guildID == 0 || joinerCharacterID == 0 || gs.characterRuntime == nil {
				return nil
			}
			g := gs.guild.Get(guildID)
			if g == nil {
				return nil
			}
			var joinerStatus dto.GuildMemberStatus
			found := false
			for _, m := range g.GetMembers() {
				if m == nil || m.GetCharacterId() != joinerCharacterID {
					continue
				}
				joinerStatus = entity.GuildMemberToDTO(m)
				found = true
				break
			}
			if !found {
				return nil
			}
			for _, m := range g.GetMembers() {
				if m == nil {
					continue
				}
				memberID := m.GetCharacterId()
				if memberID == 0 || memberID == joinerCharacterID {
					continue
				}
				if !gs.characterRuntime.Exists(memberID) {
					continue
				}
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildNewMember{
					CharacterID: memberID,
					GuildID:     guildID,
					Member:      joinerStatus,
				})
			}
			return nil
		}).
		Run()
	return nil
}
