package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type DeliverPartyUpdateJoinHandler struct{}

func (DeliverPartyUpdateJoinHandler) New() *DeliverPartyUpdateJoinHandler {
	return &DeliverPartyUpdateJoinHandler{}
}

func (h *DeliverPartyUpdateJoinHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *DeliverPartyUpdateJoin) {
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
	ch.Listener.OnPartyUpdateJoin(ch, msg.ForChannel, msg.PartyID, msg.JoinName, msg.LeaderID, msg.Members)
	if ch.GetID() != msg.JoinCharacterID {
		return
	}
	m.SyncPartyHP(ch, msg.PartyID)
}
