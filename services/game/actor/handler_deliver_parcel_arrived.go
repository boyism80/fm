package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverParcelArrivedHandler struct{}

func (DeliverParcelArrivedHandler) New() *DeliverParcelArrivedHandler {
	return &DeliverParcelArrivedHandler{}
}

func (h *DeliverParcelArrivedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverParcelArrived) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnDueyArrival(ch, msg.SenderName, msg.Quick, 1)
}
