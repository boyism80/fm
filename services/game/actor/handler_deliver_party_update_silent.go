package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateSilentHandler struct{}

func (DeliverPartyUpdateSilentHandler) New() *DeliverPartyUpdateSilentHandler {
	return &DeliverPartyUpdateSilentHandler{}
}

func (h *DeliverPartyUpdateSilentHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverPartyUpdateSilent) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateSilent(ch, msg.ForChannel, msg.PartyID, msg.LeaderID, msg.Members)
}
