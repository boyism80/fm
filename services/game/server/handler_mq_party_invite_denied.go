package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqPartyInviteDenied struct{ gs *GameServer }

func (partyMqPartyInviteDenied) New(gs *GameServer) *partyMqPartyInviteDenied {
	return &partyMqPartyInviteDenied{gs: gs}
}
func (*partyMqPartyInviteDenied) EventType() string { return "party_invite_denied" }
func (h *partyMqPartyInviteDenied) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	var payload struct {
		InviterCharacterID  uint32 `json:"inviter_character_id"`
		DeniedCharacterName string `json:"denied_character_name"`
		Action              uint8  `json:"action"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("party consumer: party_invite_denied JSON: %v", err)
		return nil
	}
	if payload.InviterCharacterID == 0 {
		return nil
	}
	pc.DeliverPartyDenyStatusToCharacter(payload.InviterCharacterID, payload.Action, payload.DeniedCharacterName)
	return nil
}
