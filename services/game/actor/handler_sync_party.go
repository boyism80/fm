package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type SyncPartyHandler struct{}

func (SyncPartyHandler) New() *SyncPartyHandler {
	return &SyncPartyHandler{}
}

func (h *SyncPartyHandler) Handle(ctx actor.Context, a *MapActor, msg *SyncParty) {
	if a.Map == nil || msg == nil || msg.Party == nil {
		return
	}
	partyID := msg.Party.GetPartyId()
	memberSet := make(map[uint32]struct{}, len(msg.Party.GetMembers()))
	for _, member := range msg.Party.GetMembers() {
		if member == nil {
			continue
		}
		memberSet[member.GetCharacterId()] = struct{}{}
	}
	for _, obj := range a.Map.GetAllPlayers() {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil {
			continue
		}
		if _, exists := memberSet[ch.GetID()]; exists {
			id := partyID
			ch.SetPartyID(&id)
			continue
		}
		if cur := ch.GetPartyID(); cur != nil && *cur == partyID {
			ch.SetPartyID(nil)
		}
	}
}
