package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
)

type DeliverGuildInviteHandler struct{}

func (DeliverGuildInviteHandler) New() *DeliverGuildInviteHandler {
	return &DeliverGuildInviteHandler{}
}

func (h *DeliverGuildInviteHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildInvite) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	if _, inGuild := ch.GetGuildID(); inGuild {
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverGuildMessage{
				CharacterID: msg.InviterCharacterID,
				Code:        pconst.GuildResponseAlreadyInGuild,
			})
		}
		return
	}
	now := clock.Now()
	for id, expiresAt := range ch.GuildInvites {
		if !now.Before(expiresAt) {
			delete(ch.GuildInvites, id)
		}
	}
	if len(ch.GuildInvites) > 0 {
		if a.GameWorld != nil {
			a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverMessage{
				CharacterID: msg.InviterCharacterID,
				MessageType: constant.MsgPinkText,
				Message:     constant.GuildInviteTargetBusyMessage,
			})
		}
		return
	}
	ch.GuildInvites[msg.GuildID] = now.Add(constant.GuildInviteDuration)
	ch.Listener.OnGuildInvite(ch, msg.GuildID, msg.InviterName)
}
