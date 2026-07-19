package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMessageHandler struct{}

func (DeliverGuildMessageHandler) New() *DeliverGuildMessageHandler {
	return &DeliverGuildMessageHandler{}
}

func (h *DeliverGuildMessageHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildMessage) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMessage(ch, msg.Code)
}
