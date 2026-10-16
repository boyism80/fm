package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
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
	pc := h.gs.systems.Get[*game.PartySystem]()
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevParty := pc.Get(evt.PartyID)
	pc.UpdateAsync(ctx, evt).Do(func() error {
		if raw == nil {
			return nil
		}
		var extra struct {
			CharacterID uint32 `json:"character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || prevParty == nil {
			return nil
		}
		pc.BroadcastDisbanded(prevParty, extra.CharacterID)
		return nil
	})
	return nil
}
