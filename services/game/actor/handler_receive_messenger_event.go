package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type ReceiveMessengerEventHandler struct{}

func (ReceiveMessengerEventHandler) New() *ReceiveMessengerEventHandler {
	return &ReceiveMessengerEventHandler{}
}

func (h *ReceiveMessengerEventHandler) Handle(_ actor.Context, a *GameLogicActor, msg *ReceiveMessengerEvent) {
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
