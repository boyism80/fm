package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateLogOnOffHandler struct{}

func (DeliverPartyUpdateLogOnOffHandler) New() *DeliverPartyUpdateLogOnOffHandler {
	return &DeliverPartyUpdateLogOnOffHandler{}
}

func (h *DeliverPartyUpdateLogOnOffHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateLogOnOff) {
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
	ch.Listener.OnPartyUpdateLogOnOff(ch, msg.ForChannel, msg.PartyID, msg.LeaderID, msg.Members)
}
