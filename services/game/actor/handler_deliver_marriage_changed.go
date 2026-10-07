package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverMarriageChangedHandler struct{}

func (DeliverMarriageChangedHandler) New() *DeliverMarriageChangedHandler {
	return &DeliverMarriageChangedHandler{}
}

func (h *DeliverMarriageChangedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverMarriageChanged) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.RefreshMarriage(ctx, msg.Event)
}
