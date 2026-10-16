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
		if cur := ch.Party().ID(); cur != nil && msg.PartyID != nil && *cur == *msg.PartyID {
			ch.Party().SetAt(nil, msg.UpdatedAt)
		}
		return
	}
	if msg.PartyID == nil {
		ch.Party().SetAt(nil, msg.UpdatedAt)
		return
	}
	id := *msg.PartyID
	ch.Party().SetAt(&id, msg.UpdatedAt)
}
