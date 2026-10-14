package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type PartyMemberJoinedHandler struct{}

func (PartyMemberJoinedHandler) New() *PartyMemberJoinedHandler {
	return &PartyMemberJoinedHandler{}
}

func (h *PartyMemberJoinedHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.PartyMemberJoined) {
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
	m.SyncPartyHP(ch)
}
