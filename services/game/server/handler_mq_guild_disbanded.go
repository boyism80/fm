package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqDisbanded struct{ gs *GameServer }

func (guildMqDisbanded) New(gs *GameServer) *guildMqDisbanded {
	return &guildMqDisbanded{gs: gs}
}

func (*guildMqDisbanded) EventType() string {
	return "disbanded"
}

func (h *guildMqDisbanded) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
			if prevGuild == nil {
				return nil
			}
			guildID := prevGuild.GetGuildId()
			memberIDs := entity.GuildMemberCharacterIDs(prevGuild.GetMembers())
			var extra struct {
				MemberCharacterIDs []uint32 `json:"member_character_ids"`
			}
			if err := json.Unmarshal(raw, &extra); err == nil && len(extra.MemberCharacterIDs) > 0 {
				memberIDs = extra.MemberCharacterIDs
			}
			for _, memberID := range memberIDs {
				if memberID == 0 {
					continue
				}
				if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
					continue
				}
				gs.EnsureSend(nil, memberID, &g_actor.DeliverGuildDisbandSelf{
					CharacterID: memberID,
					GuildID:     guildID,
				})
			}
			return nil
		}).
		Run()
	return nil
}
