package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceDisbandHandler struct{}

func (DeliverAllianceDisbandHandler) New() *DeliverAllianceDisbandHandler {
	return &DeliverAllianceDisbandHandler{}
}

func (h *DeliverAllianceDisbandHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceDisband) {
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
	ch.Listener.OnAllianceDisband(ch, msg.AllianceID)
}
