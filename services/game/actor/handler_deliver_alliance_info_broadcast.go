package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceInfoBroadcastHandler struct{}

func (DeliverAllianceInfoBroadcastHandler) New() *DeliverAllianceInfoBroadcastHandler {
	return &DeliverAllianceInfoBroadcastHandler{}
}

func (h *DeliverAllianceInfoBroadcastHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceInfoBroadcast) {
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
	ch.Listener.OnAllianceInfoBroadcast(ch, msg.Info)
}
