package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
)

type partyMqDisbanded struct{ gs *GameServer }

func (partyMqDisbanded) New(gs *GameServer) *partyMqDisbanded {
	return &partyMqDisbanded{gs: gs}
}

func (*partyMqDisbanded) EventType() string {
	return "disbanded"
}

func (h *partyMqDisbanded) Handle(ctx actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevParty := pc.Get(evt.PartyID)
	pc.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		if raw == nil {
			return nil, nil
		}
		var extra struct {
			CharacterID uint32 `json:"character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 || prevParty == nil {
			return nil, nil
		}
		pc.BroadcastDisbanded(prevParty, extra.CharacterID)
		return nil, nil
	})
	return nil
}
