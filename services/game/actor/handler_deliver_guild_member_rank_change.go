package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberRankChangeHandler struct{}

func (DeliverGuildMemberRankChangeHandler) New() *DeliverGuildMemberRankChangeHandler {
	return &DeliverGuildMemberRankChangeHandler{}
}

func (h *DeliverGuildMemberRankChangeHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverGuildMemberRankChange) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnGuildMemberRankChange(ch, msg.GuildID, msg.TargetID, msg.GuildRank)
}
