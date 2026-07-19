package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyStatusMessageHandler struct{}

func (DeliverPartyStatusMessageHandler) New() *DeliverPartyStatusMessageHandler {
	return &DeliverPartyStatusMessageHandler{}
}

func (h *DeliverPartyStatusMessageHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyStatusMessage) {
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
	ch.Listener.OnPartyStatusMessage(ch, msg.Code, msg.Name)
}
