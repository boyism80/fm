package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildLeaveSelfHandler struct{}

func (DeliverGuildLeaveSelfHandler) New() *DeliverGuildLeaveSelfHandler {
	return &DeliverGuildLeaveSelfHandler{}
}

func (h *DeliverGuildLeaveSelfHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildLeaveSelf) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildLeaveSelf(ch)
}
