package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildRankTitleChangeHandler struct{}

func (DeliverGuildRankTitleChangeHandler) New() *DeliverGuildRankTitleChangeHandler {
	return &DeliverGuildRankTitleChangeHandler{}
}

func (h *DeliverGuildRankTitleChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildRankTitleChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildRankTitleChange(ch, msg.GuildID, msg.RankTitles)
}
