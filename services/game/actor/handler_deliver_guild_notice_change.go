package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildNoticeChangeHandler struct{}

func (DeliverGuildNoticeChangeHandler) New() *DeliverGuildNoticeChangeHandler {
	return &DeliverGuildNoticeChangeHandler{}
}

func (h *DeliverGuildNoticeChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildNoticeChange) {
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
	ch.Listener.OnGuildNoticeChange(ch, msg.GuildID, msg.Notice)
}
