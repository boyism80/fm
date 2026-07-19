package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildNewMemberHandler struct{}

func (DeliverGuildNewMemberHandler) New() *DeliverGuildNewMemberHandler {
	return &DeliverGuildNewMemberHandler{}
}

func (h *DeliverGuildNewMemberHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildNewMember) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildNewMember(ch, msg.GuildID, msg.Member)
}
