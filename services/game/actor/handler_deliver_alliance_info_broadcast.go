package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceInfoBroadcastHandler struct{}

func (DeliverAllianceInfoBroadcastHandler) New() *DeliverAllianceInfoBroadcastHandler {
	return &DeliverAllianceInfoBroadcastHandler{}
}

func (h *DeliverAllianceInfoBroadcastHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceInfoBroadcast) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceInfoBroadcast(ch, msg.Info)
}
