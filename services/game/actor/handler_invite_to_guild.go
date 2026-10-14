package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type InviteToGuildHandler struct{}

func (InviteToGuildHandler) New() *InviteToGuildHandler {
	return &InviteToGuildHandler{}
}

func (h *InviteToGuildHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *InviteToGuild) {
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
	case entity.ErrGuildAlreadyJoined:
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().CallCharacter(msg.InviterCharacterID, func(_ actor.Context, inviter *entity.Character) {
				if inviter == nil {
					return
				}
				inviter.Listener.OnGuildMessage(inviter, pconst.GuildResponseAlreadyInGuild)
			})
		}
	case entity.ErrGuildInviteBusy:
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().CallCharacter(msg.InviterCharacterID, func(_ actor.Context, inviter *entity.Character) {
				if inviter == nil {
					return
				}
				inviter.Listener.OnMessage(inviter, constant.MsgPinkText, constant.GuildInviteTargetBusyMessage)
			})
		}
	}
}
