package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberOnlineChangeHandler struct{}

func (DeliverGuildMemberOnlineChangeHandler) New() *DeliverGuildMemberOnlineChangeHandler {
	return &DeliverGuildMemberOnlineChangeHandler{}
}

func (h *DeliverGuildMemberOnlineChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildMemberOnlineChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberOnlineChange(ch, msg.GuildID, msg.SubjectID, msg.Online)
}
