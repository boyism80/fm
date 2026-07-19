package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildEmblemChangeHandler struct{}

func (DeliverGuildEmblemChangeHandler) New() *DeliverGuildEmblemChangeHandler {
	return &DeliverGuildEmblemChangeHandler{}
}

func (h *DeliverGuildEmblemChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildEmblemChange) {
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
	ch.Listener.OnGuildEmblemChange(ch, msg.GuildID, msg.LogoBG, msg.LogoBGColor, msg.Logo, msg.LogoColor)
}
