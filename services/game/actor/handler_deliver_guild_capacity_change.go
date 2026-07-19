package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildCapacityChangeHandler struct{}

func (DeliverGuildCapacityChangeHandler) New() *DeliverGuildCapacityChangeHandler {
	return &DeliverGuildCapacityChangeHandler{}
}

func (h *DeliverGuildCapacityChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildCapacityChange) {
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
	ch.Listener.OnGuildCapacityChange(ch, msg.GuildID, msg.Capacity)
}
