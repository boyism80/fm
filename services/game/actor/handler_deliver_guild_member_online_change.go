package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberOnlineChangeHandler struct{}

func (DeliverGuildMemberOnlineChangeHandler) New() *DeliverGuildMemberOnlineChangeHandler {
	return &DeliverGuildMemberOnlineChangeHandler{}
}

func (h *DeliverGuildMemberOnlineChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildMemberOnlineChange) {
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
	ch.Listener.OnGuildMemberOnlineChange(ch, msg.GuildID, msg.SubjectID, msg.Online)
}
