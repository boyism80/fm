package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type partyMqLeaderChanged struct{ gs *GameServer }

func (partyMqLeaderChanged) New(gs *GameServer) *partyMqLeaderChanged {
	return &partyMqLeaderChanged{gs: gs}
}

func (*partyMqLeaderChanged) EventType() string {
	return "leader_changed"
}

func (h *partyMqLeaderChanged) Handle(ctx actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.systems.Get[*game.PartySystem]()
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	var extra struct {
		NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil
	}
	if !pc.Update(decodePartySnapshot(evt.PartyPB)) {
		return nil
	}
	if party := pc.Get(evt.PartyID); party != nil {
		pc.BroadcastLeaderChanged(party, extra.NewLeaderCharacterID, false)
	}
	return nil
}
