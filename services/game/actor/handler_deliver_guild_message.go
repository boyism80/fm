package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMessageHandler struct{}

func (DeliverGuildMessageHandler) New() *DeliverGuildMessageHandler {
	return &DeliverGuildMessageHandler{}
}

func (h *DeliverGuildMessageHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildMessage) {
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
	ch.Listener.OnGuildMessage(ch, msg.Code)
}
