package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildCapacityChangeHandler struct{}

func (DeliverGuildCapacityChangeHandler) New() *DeliverGuildCapacityChangeHandler {
	return &DeliverGuildCapacityChangeHandler{}
}

func (h *DeliverGuildCapacityChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildCapacityChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildCapacityChange(ch, msg.GuildID, msg.Capacity)
}
