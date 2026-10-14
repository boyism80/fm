package actor

import (
	"github.com/asynkron/protoactor-go/actor"

	"github.com/boyism80/fm/services/game/domain"
)

type PartyMemberLeftHandler struct{}

func (PartyMemberLeftHandler) New() *PartyMemberLeftHandler {
	return &PartyMemberLeftHandler{}
}

func (h *PartyMemberLeftHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.PartyMemberLeft) {
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
	ch.StateMachine().RequestLeave(ch, true, domain.StateMachineLeaveParty)
}
