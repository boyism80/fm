package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type PartyMemberLeftHandler struct{}

func (PartyMemberLeftHandler) New() *PartyMemberLeftHandler {
	return &PartyMemberLeftHandler{}
}

func (h *PartyMemberLeftHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *PartyMemberLeft) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.LeaverID)
	if m == nil {
		return
	}
	m.ApplyPartyLeaveDoorSync(msg.LeaverID)
	ch := m.GetPlayer(msg.LeaverID)
	if ch == nil {
		return
	}
	if sm := ch.StateMachine(); sm != nil {
		finished := sm.LeavePlayer(ctx, ch, true)
		if !finished {
			sm.CallHook("on_left_party", sm, ch)
		}
	}
}
