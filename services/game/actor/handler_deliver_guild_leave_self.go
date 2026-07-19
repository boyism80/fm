package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildLeaveSelfHandler struct{}

func (DeliverGuildLeaveSelfHandler) New() *DeliverGuildLeaveSelfHandler {
	return &DeliverGuildLeaveSelfHandler{}
}

func (h *DeliverGuildLeaveSelfHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildLeaveSelf) {
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
	ch.Listener.OnGuildLeaveSelf(ch)
}
