package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverMessageHandler struct{}

func (DeliverMessageHandler) New() *DeliverMessageHandler {
	return &DeliverMessageHandler{}
}

func (h *DeliverMessageHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverMessage) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnMessage(ch, msg.MessageType, msg.Message)
}
