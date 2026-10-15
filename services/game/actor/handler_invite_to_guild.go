package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game"
)

type InviteToGuildHandler struct{}

func (InviteToGuildHandler) New() *InviteToGuildHandler {
	return &InviteToGuildHandler{}
}

func (h *InviteToGuildHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.InviteToGuild) {
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
	case game.ErrGuildAlreadyJoined:
		if a.GameWorld != nil {
			a.GameWorld.Systems().Get[*game.DispatchSystem]().Visit(msg.InviterCharacterID, func(_ actor.Context, inviter *game.Character) {
				if inviter == nil {
					return
				}
				inviter.Listener.OnGuildMessage(inviter, pconst.GuildResponseAlreadyInGuild)
			})
		}
	case game.ErrGuildInviteBusy:
		if a.GameWorld != nil {
			a.GameWorld.Systems().Get[*game.DispatchSystem]().Visit(msg.InviterCharacterID, func(_ actor.Context, inviter *game.Character) {
				if inviter == nil {
					return
				}
				inviter.Listener.OnMessage(inviter, constant.MsgPinkText, constant.GuildInviteTargetBusyMessage)
			})
		}
	}
}
