package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type SetFamilyHandler struct{}

func (SetFamilyHandler) New() *SetFamilyHandler {
	return &SetFamilyHandler{}
}

func (h *SetFamilyHandler) Handle(_ actor.Context, a *GameLogicActor, msg *game.SetFamily) {
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
