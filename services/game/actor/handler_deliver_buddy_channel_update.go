package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverBuddyChannelUpdateHandler struct{}

func (DeliverBuddyChannelUpdateHandler) New() *DeliverBuddyChannelUpdateHandler {
	return &DeliverBuddyChannelUpdateHandler{}
}

func (h *DeliverBuddyChannelUpdateHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverBuddyChannelUpdate) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.RecipientCharacterID)
	if ch == nil {
		return
	}
	ch.BuddyList().SetChannel(msg.BuddyCharacterID, msg.Channel)
	ch.Listener.OnBuddyChannelUpdate(ch, msg.BuddyCharacterID, msg.Channel)
}
