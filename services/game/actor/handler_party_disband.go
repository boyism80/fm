package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type PartyDisbandHandler struct{}

func (PartyDisbandHandler) New() *PartyDisbandHandler {
	return &PartyDisbandHandler{}
}

func (h *PartyDisbandHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.PartyDisband) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	m.SyncDoors(msg.FormerMemberIDs)
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil || ch.StateMachine() == nil {
		return
	}
	ch.StateMachine().CallHook("on_disband_party")
}
