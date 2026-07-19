package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildExpelSelfHandler struct{}

func (DeliverGuildExpelSelfHandler) New() *DeliverGuildExpelSelfHandler {
	return &DeliverGuildExpelSelfHandler{}
}

func (h *DeliverGuildExpelSelfHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildExpelSelf) {
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
	ch.Listener.OnGuildExpelledSelf(ch, msg.GuildID)
}
