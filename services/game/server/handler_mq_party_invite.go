package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqPartyInvite struct{ gs *GameServer }

func (partyMqPartyInvite) New(gs *GameServer) *partyMqPartyInvite {
	return &partyMqPartyInvite{gs: gs}
}
func (*partyMqPartyInvite) EventType() string { return "party_invite" }
func (h *partyMqPartyInvite) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	var payload struct {
		TargetCharacterID uint32 `json:"target_character_id"`
		PartyID           uint32 `json:"party_id"`
		InviterName       string `json:"inviter_name"`
		PartySearch       bool   `json:"party_search"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("party consumer: party_invite JSON: %v", err)
		return nil
	}
	if payload.TargetCharacterID == 0 {
		return nil
	}
	pc.DeliverPartyInviteToCharacter(payload.TargetCharacterID, payload.PartyID, payload.InviterName, payload.PartySearch)
	return nil
}
