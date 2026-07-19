package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceDisbandHandler struct{}

func (DeliverAllianceDisbandHandler) New() *DeliverAllianceDisbandHandler {
	return &DeliverAllianceDisbandHandler{}
}

func (h *DeliverAllianceDisbandHandler) Handle(ctx actor.Context, a *MapActor, msg *DeliverAllianceDisband) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	ch.Listener.OnAllianceDisband(ch, msg.AllianceID)
}
