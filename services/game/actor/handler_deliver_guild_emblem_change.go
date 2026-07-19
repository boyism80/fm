package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildEmblemChangeHandler struct{}

func (DeliverGuildEmblemChangeHandler) New() *DeliverGuildEmblemChangeHandler {
	return &DeliverGuildEmblemChangeHandler{}
}

func (h *DeliverGuildEmblemChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildEmblemChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildEmblemChange(ch, msg.GuildID, msg.LogoBG, msg.LogoBGColor, msg.Logo, msg.LogoColor)
}
