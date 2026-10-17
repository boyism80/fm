package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type partyMqLogOnOff struct{ gs *GameServer }

func (partyMqLogOnOff) New(gs *GameServer) *partyMqLogOnOff {
	return &partyMqLogOnOff{gs: gs}
}

func (*partyMqLogOnOff) EventType() string {
	return "log_onoff"
}

func (h *partyMqLogOnOff) Handle(ctx actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.systems.Get[*game.PartySystem]()
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	if !pc.Update(decodePartySnapshot(evt.PartyPB)) {
		return nil
	}
	if party := pc.Get(evt.PartyID); party != nil {
		pc.BroadcastLogOnOff(party)
	}
	return nil
}
