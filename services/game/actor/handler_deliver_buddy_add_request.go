package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverBuddyAddRequestHandler struct{}

func (DeliverBuddyAddRequestHandler) New() *DeliverBuddyAddRequestHandler {
	return &DeliverBuddyAddRequestHandler{}
}

func (h *DeliverBuddyAddRequestHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverBuddyAddRequest) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnBuddyAddRequest(ch, msg.FromCharacterID, msg.FromName)
}
