package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceMemberFieldsChangeHandler struct{}

func (DeliverAllianceMemberFieldsChangeHandler) New() *DeliverAllianceMemberFieldsChangeHandler {
	return &DeliverAllianceMemberFieldsChangeHandler{}
}

func (h *DeliverAllianceMemberFieldsChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceMemberFieldsChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberFieldsChange(ch, msg.AllianceID, msg.GuildID, msg.SubjectID, msg.Level, msg.ClassID)
}
