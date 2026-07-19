package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type DeliverAllianceInviteHandler struct{}

func (DeliverAllianceInviteHandler) New() *DeliverAllianceInviteHandler {
	return &DeliverAllianceInviteHandler{}
}

func (h *DeliverAllianceInviteHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceInvite) {
	if a.Map == nil || msg == nil || a.GameWorld == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	guildID, inGuild := ch.GetGuildID()
	if !inGuild || guildID != msg.TargetGuildID {
		return
	}
	if g := a.GameWorld.GetGuildSystem().Get(guildID); g != nil {
		if _, joined := g.GetAllianceID(); joined {
			return
		}
	}
	if !a.GameWorld.GetGuildSystem().TrySetAllianceInvite(msg.TargetGuildID, msg.AllianceID, msg.ExpiresAt) {
		a.GameWorld.GetDispatchSystem().SendTo(msg.InviterCharacterID, &DeliverMessage{
			CharacterID: msg.InviterCharacterID,
			MessageType: constant.MsgPinkText,
			Message:     constant.GuildInviteTargetBusyMessage,
		})
		return
	}
	ch.Listener.OnAllianceInvite(ch, msg.InviterGuildID, msg.InviterName, msg.AllianceName)
}
