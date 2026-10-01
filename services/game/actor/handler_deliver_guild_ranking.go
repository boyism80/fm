package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverGuildRankingHandler struct{}

func (DeliverGuildRankingHandler) New() *DeliverGuildRankingHandler {
	return &DeliverGuildRankingHandler{}
}

func (h *DeliverGuildRankingHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverGuildRanking) {
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
	ch.Listener.OnGuildRanking(ch, msg.NPCID, msg.Entries)
}
