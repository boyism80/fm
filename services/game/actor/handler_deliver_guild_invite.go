package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type DeliverGuildInviteHandler struct{}

func (DeliverGuildInviteHandler) New() *DeliverGuildInviteHandler {
	return &DeliverGuildInviteHandler{}
}

func (h *DeliverGuildInviteHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildInvite) {
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
			a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverGuildMessage{
				CharacterID: msg.InviterCharacterID,
				Code:        pconst.GuildResponseAlreadyInGuild,
			})
		}
	case entity.ErrGuildInviteBusy:
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverMessage{
				CharacterID: msg.InviterCharacterID,
				MessageType: constant.MsgPinkText,
				Message:     constant.GuildInviteTargetBusyMessage,
			})
		}
	}
}
