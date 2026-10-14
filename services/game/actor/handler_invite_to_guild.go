package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type InviteToGuildHandler struct{}

func (InviteToGuildHandler) New() *InviteToGuildHandler {
	return &InviteToGuildHandler{}
}

func (h *InviteToGuildHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.InviteToGuild) {
	if msg == nil {
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
	switch ch.Guild.Invite(msg.GuildID, msg.InviterName) {
	case domain.ErrGuildAlreadyJoined:
		if a.GameWorld != nil {
			a.GameWorld.Systems().Get[*domain.DispatchSystem]().CallCharacter(msg.InviterCharacterID, func(_ actor.Context, inviter *domain.Character) {
				if inviter == nil {
					return
				}
				inviter.Listener.OnGuildMessage(inviter, pconst.GuildResponseAlreadyInGuild)
			})
		}
	case domain.ErrGuildInviteBusy:
		if a.GameWorld != nil {
			a.GameWorld.Systems().Get[*domain.DispatchSystem]().CallCharacter(msg.InviterCharacterID, func(_ actor.Context, inviter *domain.Character) {
				if inviter == nil {
					return
				}
				inviter.Listener.OnMessage(inviter, constant.MsgPinkText, constant.GuildInviteTargetBusyMessage)
			})
		}
	}
}
