package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
)

type partyMqMemberJoined struct{ gs *GameServer }

func (partyMqMemberJoined) New(gs *GameServer) *partyMqMemberJoined {
	return &partyMqMemberJoined{gs: gs}
}

func (*partyMqMemberJoined) EventType() string {
	return "member_joined"
}

func (h *partyMqMemberJoined) Handle(ctx actor.Context, raw json.RawMessage) error {
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
			CharacterID uint32 `json:"character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
			return nil
		}
		party := pc.Get(evt.PartyID)
		if party != nil {
			pc.BroadcastMemberJoined(party, extra.CharacterID)
		}
		return nil
	})
	return nil
}
