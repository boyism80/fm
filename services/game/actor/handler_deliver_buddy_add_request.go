package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverBuddyAddRequestHandler struct{}

func (DeliverBuddyAddRequestHandler) New() *DeliverBuddyAddRequestHandler {
	return &DeliverBuddyAddRequestHandler{}
}

func (h *DeliverBuddyAddRequestHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverBuddyAddRequest) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.RecipientCharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnBuddyAddRequest(ch, msg.FromCharacterID, msg.FromName)
}
