package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverMessengerHandler struct{}

func (DeliverMessengerHandler) New() *DeliverMessengerHandler {
	return &DeliverMessengerHandler{}
}

func (h *DeliverMessengerHandler) Handle(_ actor.Context, a *GameLogicActor, msg *DeliverMessenger) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.LoggedOut() {
		return
	}
	ch.Messenger.Receive(msg.Event)
}
