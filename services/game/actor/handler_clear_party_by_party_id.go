package actor

import (
	"github.com/asynkron/protoactor-go/actor"
)

type ClearPartyByPartyIDHandler struct{}

func (ClearPartyByPartyIDHandler) New() *ClearPartyByPartyIDHandler {
	return &ClearPartyByPartyIDHandler{}
}

func (h *ClearPartyByPartyIDHandler) Handle(ctx actor.Context, a *GameLogicActor, msg *ClearPartyByPartyID) {
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
	if cur := ch.Party.ID(); cur != nil && *cur == msg.PartyID {
		ch.Party.SetID(nil)
	}
}
