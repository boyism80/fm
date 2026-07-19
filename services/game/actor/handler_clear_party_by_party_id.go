package actor

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type ClearPartyByPartyIDHandler struct{}

func (ClearPartyByPartyIDHandler) New() *ClearPartyByPartyIDHandler {
	return &ClearPartyByPartyIDHandler{}
}

func (h *ClearPartyByPartyIDHandler) Handle(ctx actor.Context, a *MapActor, msg *ClearPartyByPartyID) {
	if a.Map == nil || msg == nil {
		return
	}
	for _, obj := range a.Map.GetAllPlayers() {
		ch, ok := obj.(*entity.Character)
		if !ok || ch == nil {
			continue
		}
		if cur := ch.GetPartyID(); cur != nil && *cur == msg.PartyID {
			ch.SetPartyID(nil)
		}
	}
}
