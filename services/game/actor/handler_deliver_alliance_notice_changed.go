package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceNoticeChangedHandler struct{}

func (DeliverAllianceNoticeChangedHandler) New() *DeliverAllianceNoticeChangedHandler {
	return &DeliverAllianceNoticeChangedHandler{}
}

func (h *DeliverAllianceNoticeChangedHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceNoticeChanged) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceNoticeChanged(ch, msg.Info)
}
