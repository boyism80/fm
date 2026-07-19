package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceMemberRankChangedHandler struct{}

func (DeliverAllianceMemberRankChangedHandler) New() *DeliverAllianceMemberRankChangedHandler {
	return &DeliverAllianceMemberRankChangedHandler{}
}

func (h *DeliverAllianceMemberRankChangedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceMemberRankChanged) {
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
	ch.Listener.OnAllianceMemberRankChanged(ch, msg.Info, msg.Guilds)
}
