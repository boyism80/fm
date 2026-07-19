package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateLeaderChangeHandler struct{}

func (DeliverPartyUpdateLeaderChangeHandler) New() *DeliverPartyUpdateLeaderChangeHandler {
	return &DeliverPartyUpdateLeaderChangeHandler{}
}

func (h *DeliverPartyUpdateLeaderChangeHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateLeaderChange) {
	if msg == nil || msg.NewLeaderCharacterID == 0 {
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
	ch.Listener.OnPartyUpdateLeaderChange(ch, msg.NewLeaderCharacterID, msg.ByDisconnect)
}
