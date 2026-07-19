package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateDisbandHandler struct{}

func (DeliverPartyUpdateDisbandHandler) New() *DeliverPartyUpdateDisbandHandler {
	return &DeliverPartyUpdateDisbandHandler{}
}

func (h *DeliverPartyUpdateDisbandHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverPartyUpdateDisband) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateDisband(ch, msg.PartyID, msg.LeaderID)
}
