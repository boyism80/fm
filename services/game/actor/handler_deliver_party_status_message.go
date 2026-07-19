package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyStatusMessageHandler struct{}

func (DeliverPartyStatusMessageHandler) New() *DeliverPartyStatusMessageHandler {
	return &DeliverPartyStatusMessageHandler{}
}

func (h *DeliverPartyStatusMessageHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverPartyStatusMessage) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyStatusMessage(ch, msg.Code, msg.Name)
}
