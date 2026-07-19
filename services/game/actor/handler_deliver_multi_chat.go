package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverMultiChatHandler struct{}

func (DeliverMultiChatHandler) New() *DeliverMultiChatHandler {
	return &DeliverMultiChatHandler{}
}

func (h *DeliverMultiChatHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverMultiChat) {
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
	ch.Listener.OnMultiChat(ch, msg.Mode, msg.SenderName, msg.Message)
}
