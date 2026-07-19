package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type PartyDisbandHandler struct{}

func (PartyDisbandHandler) New() *PartyDisbandHandler {
	return &PartyDisbandHandler{}
}

func (h *PartyDisbandHandler) Handle(ctx actor.Context, a *MapActor, msg *PartyDisband) {
	if a.Map == nil || msg == nil {
		return
	}
	a.Map.ApplyPartyDisbandDoorSync(msg.FormerMemberIDs)
	seen := make(map[*entity.StateMachine]struct{})
	for _, id := range msg.FormerMemberIDs {
		ch := a.Map.GetPlayer(id)
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
		sm.CallHook("on_disband_party", sm)
	}
}
