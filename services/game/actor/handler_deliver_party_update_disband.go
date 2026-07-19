package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateDisbandHandler struct{}

func (DeliverPartyUpdateDisbandHandler) New() *DeliverPartyUpdateDisbandHandler {
	return &DeliverPartyUpdateDisbandHandler{}
}

func (h *DeliverPartyUpdateDisbandHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateDisband) {
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
	ch.Listener.OnPartyUpdateDisband(ch, msg.PartyID, msg.LeaderID)
}
