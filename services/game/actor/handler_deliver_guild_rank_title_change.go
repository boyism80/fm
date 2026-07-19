package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildRankTitleChangeHandler struct{}

func (DeliverGuildRankTitleChangeHandler) New() *DeliverGuildRankTitleChangeHandler {
	return &DeliverGuildRankTitleChangeHandler{}
}

func (h *DeliverGuildRankTitleChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildRankTitleChange) {
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
	ch.Listener.OnGuildRankTitleChange(ch, msg.GuildID, msg.RankTitles)
}
