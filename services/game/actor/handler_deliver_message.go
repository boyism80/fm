package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverMessageHandler struct{}

func (DeliverMessageHandler) New() *DeliverMessageHandler {
	return &DeliverMessageHandler{}
}

func (h *DeliverMessageHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverMessage) {
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
	ch.Listener.OnMessage(ch, msg.MessageType, msg.Message)
}
