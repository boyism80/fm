package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type SpouseMovedHandler struct{}

func (SpouseMovedHandler) New() *SpouseMovedHandler {
	return &SpouseMovedHandler{}
}

func (h *SpouseMovedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.SpouseMoved) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.Wedding.Marriage == nil || ch.Wedding.Marriage.PartnerID(ch.GetID()) != msg.SpouseID {
		return
	}
	ch.Listener.OnSpouseMap(ch, msg.MapID, msg.SpouseID)
	if msg.Reply {
		ch.Wedding.NotifySpouseMap(ctx, false)
	}
}
