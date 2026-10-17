package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type SyncCharacterPartyStateHandler struct{}

func (SyncCharacterPartyStateHandler) New() *SyncCharacterPartyStateHandler {
	return &SyncCharacterPartyStateHandler{}
}

func (h *SyncCharacterPartyStateHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *game.SyncCharacterPartyState) {
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
	if msg.Clear {
		ch.Party().Leave(msg.PartyID, msg.Revision)
		return
	}
	ch.Party().Join(msg.PartyID, msg.Revision)
}
