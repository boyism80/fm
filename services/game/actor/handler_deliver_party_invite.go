package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyInviteHandler struct{}

func (DeliverPartyInviteHandler) New() *DeliverPartyInviteHandler {
	return &DeliverPartyInviteHandler{}
}

func (h *DeliverPartyInviteHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyInvite) {
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
	ch.Listener.OnPartyInvite(ch, msg.PartyID, msg.InviterName, msg.PartySearch)
}
