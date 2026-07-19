package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type PartyMemberLeftHandler struct{}

func (PartyMemberLeftHandler) New() *PartyMemberLeftHandler {
	return &PartyMemberLeftHandler{}
}

func (h *PartyMemberLeftHandler) Handle(ctx actor.Context, a *MapActor, msg *PartyMemberLeft) {
	if a.Map == nil || msg == nil {
		return
	}
	a.Map.ApplyPartyLeaveDoorSync(msg.LeaverID)
	ch := a.Map.GetPlayer(msg.LeaverID)
	if ch == nil {
		return
	}
	if sm := ch.StateMachine(); sm != nil {
		sm.CallHook("on_left_party", sm, ch)
	}
}
