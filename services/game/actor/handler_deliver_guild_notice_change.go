package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildNoticeChangeHandler struct{}

func (DeliverGuildNoticeChangeHandler) New() *DeliverGuildNoticeChangeHandler {
	return &DeliverGuildNoticeChangeHandler{}
}

func (h *DeliverGuildNoticeChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildNoticeChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildNoticeChange(ch, msg.GuildID, msg.Notice)
}
