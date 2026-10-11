package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverNoteArrivedHandler struct{}

func (DeliverNoteArrivedHandler) New() *DeliverNoteArrivedHandler {
	return &DeliverNoteArrivedHandler{}
}

func (h *DeliverNoteArrivedHandler) Handle(_ actor.Context, a *GameLogicActor, msg *DeliverNoteArrived) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.LoggedOut() {
		return
	}
	ch.Listener.OnNoteArrived(ch)
}
