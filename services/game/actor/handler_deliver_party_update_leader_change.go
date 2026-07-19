package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateLeaderChangeHandler struct{}

func (DeliverPartyUpdateLeaderChangeHandler) New() *DeliverPartyUpdateLeaderChangeHandler {
	return &DeliverPartyUpdateLeaderChangeHandler{}
}

func (h *DeliverPartyUpdateLeaderChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverPartyUpdateLeaderChange) {
	if a.Map == nil || msg == nil || msg.NewLeaderCharacterID == 0 {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateLeaderChange(ch, msg.NewLeaderCharacterID, msg.ByDisconnect)
}
