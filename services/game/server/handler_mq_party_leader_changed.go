package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
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
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	pc.UpdateAsync(ctx, evt).Do(func() error {
		if raw == nil {
			return nil
		}
		var extra struct {
			NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || extra.NewLeaderCharacterID == 0 {
			return nil
		}
		party := pc.Get(evt.PartyID)
		if party != nil {
			pc.BroadcastLeaderChanged(party, extra.NewLeaderCharacterID, false)
		}
		return nil
	})
	return nil
}
