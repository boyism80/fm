package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type SetBuddyChannelHandler struct{}

func (SetBuddyChannelHandler) New() *SetBuddyChannelHandler {
	return &SetBuddyChannelHandler{}
}

func (h *SetBuddyChannelHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.SetBuddyChannel) {
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
	ch.Buddies.SetChannel(msg.BuddyCharacterID, msg.Channel)
	ch.Listener.OnBuddyChannelUpdate(ch, msg.BuddyCharacterID, msg.Channel)
}
