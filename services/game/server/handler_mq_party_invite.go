package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type partyMqPartyInvite struct{ gs *GameServer }

func (partyMqPartyInvite) New(gs *GameServer) *partyMqPartyInvite {
	return &partyMqPartyInvite{gs: gs}
}

func (*partyMqPartyInvite) EventType() string {
	return "party_invite"
}

func (h *partyMqPartyInvite) Handle(_ actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	var payload struct {
		TargetCharacterID uint32 `json:"target_character_id"`
		PartyID           uint32 `json:"party_id"`
		InviterName       string `json:"inviter_name"`
		PartySearch       bool   `json:"party_search"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.TargetCharacterID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.DispatchSystem]().Visit(payload.TargetCharacterID, func(_ actor.Context, ch *game.Character) {
		if ch == nil {
			return
		}
		ch.Listener.OnPartyInvite(ch, payload.PartyID, payload.InviterName, payload.PartySearch)
	})
	return nil
}
