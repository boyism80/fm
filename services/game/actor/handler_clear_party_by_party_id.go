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
	if msg == nil {
		return
	}
	for _, m := range a.Maps() {
		for _, obj := range m.GetAllPlayers() {
			ch, ok := obj.(*entity.Character)
			if !ok || ch == nil {
				continue
			}
			if cur := ch.GetPartyID(); cur != nil && *cur == msg.PartyID {
				ch.SetPartyID(nil)
			}
		}
	}
}
