package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildExpelSelfHandler struct{}

func (DeliverGuildExpelSelfHandler) New() *DeliverGuildExpelSelfHandler {
	return &DeliverGuildExpelSelfHandler{}
}

func (h *DeliverGuildExpelSelfHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildExpelSelf) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildExpelledSelf(ch, msg.GuildID)
}
