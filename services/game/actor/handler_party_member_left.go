package actor

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/boyism80/fm/services/game/entity"
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
	m.SyncDoors([]uint32{msg.LeaverID})
	ch := m.GetPlayer(msg.LeaverID)
	if ch == nil {
		return
	}
	if sm := ch.StateMachine(); sm != nil {
		finished := sm.LeavePlayer(ctx, ch, true, entity.StateMachineLeaveParty)
		if !finished {
			sm.CallHook("on_left_party", ch)
		}
	}
}
