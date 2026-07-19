package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildDisbandSelfHandler struct{}

func (DeliverGuildDisbandSelfHandler) New() *DeliverGuildDisbandSelfHandler {
	return &DeliverGuildDisbandSelfHandler{}
}

func (h *DeliverGuildDisbandSelfHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildDisbandSelf) {
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
	ch.Listener.OnGuildDisbandSelf(ch, msg.GuildID)
}
