package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberLeftHandler struct{}

func (DeliverGuildMemberLeftHandler) New() *DeliverGuildMemberLeftHandler {
	return &DeliverGuildMemberLeftHandler{}
}

func (h *DeliverGuildMemberLeftHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildMemberLeft) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberLeft(ch, msg.GuildID, msg.TargetID, msg.TargetName, msg.WasExpelled)
}
