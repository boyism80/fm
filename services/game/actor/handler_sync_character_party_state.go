package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type SyncCharacterPartyStateHandler struct{}

func (SyncCharacterPartyStateHandler) New() *SyncCharacterPartyStateHandler {
	return &SyncCharacterPartyStateHandler{}
}

func (h *SyncCharacterPartyStateHandler) Handle(ctx actor.Context, a *MapActor, msg *SyncCharacterPartyState) {
	if a.Map == nil || msg == nil {
		return
	}
	ch := a.Map.GetPlayer(msg.CharacterID)
	if ch == nil {
		return
	}
	if msg.PartyID == nil {
		ch.SetPartyID(nil)
	} else {
		id := *msg.PartyID
		ch.SetPartyID(&id)
	}
}
