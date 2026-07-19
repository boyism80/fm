package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceLeaderChangedHandler struct{}

func (DeliverAllianceLeaderChangedHandler) New() *DeliverAllianceLeaderChangedHandler {
	return &DeliverAllianceLeaderChangedHandler{}
}

func (h *DeliverAllianceLeaderChangedHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceLeaderChanged) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceLeaderChanged(ch, msg.AllianceID, msg.OldLeaderID, msg.NewLeaderID, msg.Info, msg.Guilds)
}
