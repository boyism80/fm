package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateSilentHandler struct{}

func (DeliverPartyUpdateSilentHandler) New() *DeliverPartyUpdateSilentHandler {
	return &DeliverPartyUpdateSilentHandler{}
}

func (h *DeliverPartyUpdateSilentHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateSilent) {
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
	ch.Listener.OnPartyUpdateSilent(ch, msg.ForChannel, msg.PartyID, msg.LeaderID, msg.Members)
}
