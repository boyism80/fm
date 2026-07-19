package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceNoticeChangedHandler struct{}

func (DeliverAllianceNoticeChangedHandler) New() *DeliverAllianceNoticeChangedHandler {
	return &DeliverAllianceNoticeChangedHandler{}
}

func (h *DeliverAllianceNoticeChangedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceNoticeChanged) {
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
	ch.Listener.OnAllianceNoticeChanged(ch, msg.Info)
}
