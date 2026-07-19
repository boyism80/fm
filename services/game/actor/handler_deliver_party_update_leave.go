package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateLeaveHandler struct{}

func (DeliverPartyUpdateLeaveHandler) New() *DeliverPartyUpdateLeaveHandler {
	return &DeliverPartyUpdateLeaveHandler{}
}

func (h *DeliverPartyUpdateLeaveHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateLeave) {
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
	if msg.Expelled {
		ch.Listener.OnPartyUpdateExpel(ch, msg.ForChannel, msg.PartyID, msg.TargetID, msg.TargetName, msg.LeaderID, msg.Members)
		return
	}
	ch.Listener.OnPartyUpdateLeave(ch, msg.ForChannel, msg.PartyID, msg.TargetID, msg.TargetName, msg.LeaderID, msg.Members)
}
