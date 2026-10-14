package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type InviteToAllianceHandler struct{}

func (InviteToAllianceHandler) New() *InviteToAllianceHandler {
	return &InviteToAllianceHandler{}
}

func (h *InviteToAllianceHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.InviteToAlliance) {
	if msg == nil || a.GameWorld == nil {
		return
	}
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	guildID, inGuild := ch.Guild.ID()
	if !inGuild || guildID != msg.TargetGuildID {
		return
	}
	if g := a.GameWorld.Systems().Get[*domain.GuildSystem]().Get(guildID); g != nil {
		if _, joined := g.GetAllianceID(); joined {
			return
		}
	}
	if !a.GameWorld.Systems().Get[*domain.GuildSystem]().TrySetAllianceInvite(msg.TargetGuildID, msg.AllianceID, msg.ExpiresAt) {
		a.GameWorld.Systems().Get[*domain.DispatchSystem]().CallCharacter(msg.InviterCharacterID, func(_ actor.Context, inviter *domain.Character) {
			if inviter == nil {
				return
			}
			inviter.Listener.OnMessage(inviter, constant.MsgPinkText, constant.GuildInviteTargetBusyMessage)
		})
		return
	}
	ch.Listener.OnAllianceInvite(ch, msg.InviterGuildID, msg.InviterName, msg.AllianceName)
}
