package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqPartyInviteDenied struct {
	partyMqHandler
}

func (partyMqPartyInviteDenied) New(gs *GameServer) *partyMqPartyInviteDenied {
	return &partyMqPartyInviteDenied{partyMqHandler: partyMqHandler{gs: gs}}
}

func (*partyMqPartyInviteDenied) EventType() string {
	return "party_invite_denied"
}

func (h *partyMqPartyInviteDenied) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	var payload struct {
		InviterCharacterID  uint32 `json:"inviter_character_id"`
		DeniedCharacterName string `json:"denied_character_name"`
		Action              uint8  `json:"action"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.InviterCharacterID == 0 {
		return nil
	}
	h.sendDenyStatus(payload.InviterCharacterID, payload.Action, payload.DeniedCharacterName)
	return nil
}
