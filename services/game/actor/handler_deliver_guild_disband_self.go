package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildDisbandSelfHandler struct{}

func (DeliverGuildDisbandSelfHandler) New() *DeliverGuildDisbandSelfHandler {
	return &DeliverGuildDisbandSelfHandler{}
}

func (h *DeliverGuildDisbandSelfHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildDisbandSelf) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildDisbandSelf(ch, msg.GuildID)
}
