package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverBuddyChannelUpdateHandler struct{}

func (DeliverBuddyChannelUpdateHandler) New() *DeliverBuddyChannelUpdateHandler {
	return &DeliverBuddyChannelUpdateHandler{}
}

func (h *DeliverBuddyChannelUpdateHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverBuddyChannelUpdate) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.RecipientCharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	ch.BuddyList().SetChannel(msg.BuddyCharacterID, msg.Channel)
	ch.Listener.OnBuddyChannelUpdate(ch, msg.BuddyCharacterID, msg.Channel)
}
