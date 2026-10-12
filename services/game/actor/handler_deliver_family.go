package actor

import "github.com/asynkron/protoactor-go/actor"

type DeliverFamilyHandler struct{}

func (DeliverFamilyHandler) New() *DeliverFamilyHandler {
	return &DeliverFamilyHandler{}
}

func (h *DeliverFamilyHandler) Handle(_ actor.Context, a *GameLogicActor, msg *DeliverFamily) {
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.LoggedOut() {
		return
	}

	ch.Family.Set(msg.Family)
}
