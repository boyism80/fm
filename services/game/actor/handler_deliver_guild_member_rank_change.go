package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildMemberRankChangeHandler struct{}

func (DeliverGuildMemberRankChangeHandler) New() *DeliverGuildMemberRankChangeHandler {
	return &DeliverGuildMemberRankChangeHandler{}
}

func (h *DeliverGuildMemberRankChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildMemberRankChange) {
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
	ch.Listener.OnGuildMemberRankChange(ch, msg.GuildID, msg.TargetID, msg.GuildRank)
}
