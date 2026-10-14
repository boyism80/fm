package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/domain"
)

type SyncCharacterPartyStateHandler struct{}

func (SyncCharacterPartyStateHandler) New() *SyncCharacterPartyStateHandler {
	return &SyncCharacterPartyStateHandler{}
}

func (h *SyncCharacterPartyStateHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *domain.SyncCharacterPartyState) {
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
	if msg.PartyID == nil {
		ch.Party.SetID(nil)
	} else {
		id := *msg.PartyID
		ch.Party.SetID(&id)
	}
}
