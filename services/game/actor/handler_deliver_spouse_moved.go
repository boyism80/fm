package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverSpouseMovedHandler struct{}

func (DeliverSpouseMovedHandler) New() *DeliverSpouseMovedHandler {
	return &DeliverSpouseMovedHandler{}
}

func (h *DeliverSpouseMovedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverSpouseMoved) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.Marriage == nil || ch.Marriage.PartnerID(ch.GetID()) != msg.SpouseID {
		return
	}
	ch.Listener.OnSpouseMap(ch, msg.MapID, msg.SpouseID)
	if msg.Reply {
		ch.NotifySpouseMap(ctx, false)
	}
}
