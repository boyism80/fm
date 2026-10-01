package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildGPChangeHandler struct{}

func (DeliverGuildGPChangeHandler) New() *DeliverGuildGPChangeHandler {
	return &DeliverGuildGPChangeHandler{}
}

func (h *DeliverGuildGPChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildGPChange) {
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
	ch.Listener.OnGuildGPChange(ch, msg.GuildID, msg.GP, msg.Level, msg.Amount)
}
