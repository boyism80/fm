package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceMemberFieldsChangeHandler struct{}

func (DeliverAllianceMemberFieldsChangeHandler) New() *DeliverAllianceMemberFieldsChangeHandler {
	return &DeliverAllianceMemberFieldsChangeHandler{}
}

func (h *DeliverAllianceMemberFieldsChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceMemberFieldsChange) {
	if msg == nil {
		return
	}
	m := a.GetCharacter(msg.CharacterID)
	if m == nil {
		return
	}
	ch := m.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberFieldsChange(ch, msg.AllianceID, msg.GuildID, msg.SubjectID, msg.Level, msg.ClassID)
}
