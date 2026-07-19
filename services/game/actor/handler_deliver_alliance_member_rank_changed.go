package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceMemberRankChangedHandler struct{}

func (DeliverAllianceMemberRankChangedHandler) New() *DeliverAllianceMemberRankChangedHandler {
	return &DeliverAllianceMemberRankChangedHandler{}
}

func (h *DeliverAllianceMemberRankChangedHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceMemberRankChanged) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceMemberRankChanged(ch, msg.Info, msg.Guilds)
}
