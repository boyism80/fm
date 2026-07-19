package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverAllianceLeaderChangedHandler struct{}

func (DeliverAllianceLeaderChangedHandler) New() *DeliverAllianceLeaderChangedHandler {
	return &DeliverAllianceLeaderChangedHandler{}
}

func (h *DeliverAllianceLeaderChangedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverAllianceLeaderChanged) {
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
	ch.Listener.OnAllianceLeaderChanged(ch, msg.AllianceID, msg.OldLeaderID, msg.NewLeaderID, msg.Info, msg.Guilds)
}
