package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberFieldsChangeHandler struct{}

func (DeliverGuildMemberFieldsChangeHandler) New() *DeliverGuildMemberFieldsChangeHandler {
	return &DeliverGuildMemberFieldsChangeHandler{}
}

func (h *DeliverGuildMemberFieldsChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildMemberFieldsChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberFieldsChange(ch, msg.GuildID, msg.SubjectID, msg.Level, msg.ClassID)
}
