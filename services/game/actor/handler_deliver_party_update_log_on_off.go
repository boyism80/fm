package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateLogOnOffHandler struct{}

func (DeliverPartyUpdateLogOnOffHandler) New() *DeliverPartyUpdateLogOnOffHandler {
	return &DeliverPartyUpdateLogOnOffHandler{}
}

func (h *DeliverPartyUpdateLogOnOffHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverPartyUpdateLogOnOff) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyUpdateLogOnOff(ch, msg.ForChannel, msg.PartyID, msg.LeaderID, msg.Members)
}
