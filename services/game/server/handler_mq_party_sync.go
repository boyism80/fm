package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type partyMqPartySync struct{ gs *GameServer }

func (partyMqPartySync) New(gs *GameServer) *partyMqPartySync {
	return &partyMqPartySync{gs: gs}
}
func (*partyMqPartySync) EventType() string { return "party_sync" }
func (h *partyMqPartySync) Handle(ctx actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	pc := gs.systems.Get[*game.PartySystem]()
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	if !pc.Update(decodePartySnapshot(evt.PartyPB)) {
		return nil
	}
	if party := pc.Get(evt.PartyID); party != nil {
		pc.BroadcastPartySilent(party)
	}
	return nil
}
