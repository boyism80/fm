package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverMultiChatHandler struct{}

func (DeliverMultiChatHandler) New() *DeliverMultiChatHandler {
	return &DeliverMultiChatHandler{}
}

func (h *DeliverMultiChatHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverMultiChat) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnMultiChat(ch, msg.Mode, msg.SenderName, msg.Message)
}
