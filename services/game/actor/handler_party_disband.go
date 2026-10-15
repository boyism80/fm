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
	seen := make(map[*game.StateMachine]struct{})
	for _, m := range a.Maps() {
		m.SyncDoors(msg.FormerMemberIDs)
		for _, id := range msg.FormerMemberIDs {
			ch := m.GetPlayer(id)
			if ch == nil {
				continue
			}
			sm := ch.StateMachine()
			if sm == nil {
				continue
			}
			if _, ok := seen[sm]; ok {
				continue
			}
			seen[sm] = struct{}{}
			sm.CallHook("on_disband_party")
		}
	}
}
