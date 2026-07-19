package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyInviteHandler struct{}

func (DeliverPartyInviteHandler) New() *DeliverPartyInviteHandler {
	return &DeliverPartyInviteHandler{}
}

func (h *DeliverPartyInviteHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverPartyInvite) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnPartyInvite(ch, msg.PartyID, msg.InviterName, msg.PartySearch)
}
