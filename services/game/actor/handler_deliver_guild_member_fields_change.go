package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberFieldsChangeHandler struct{}

func (DeliverGuildMemberFieldsChangeHandler) New() *DeliverGuildMemberFieldsChangeHandler {
	return &DeliverGuildMemberFieldsChangeHandler{}
}

func (h *DeliverGuildMemberFieldsChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildMemberFieldsChange) {
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
	ch.Listener.OnGuildMemberFieldsChange(ch, msg.GuildID, msg.SubjectID, msg.Level, msg.ClassID)
}
