package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"

	g_actor "github.com/boyism80/fm/services/game/actor"
)

type partyMqPartyChat struct{ gs *GameServer }

func (partyMqPartyChat) New(gs *GameServer) *partyMqPartyChat {
	return &partyMqPartyChat{gs: gs}
}

func (*partyMqPartyChat) EventType() string { return "party_chat" }

func (h *partyMqPartyChat) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		TargetCharacterID uint32 `json:"target_character_id"`
		ChatMode          uint32 `json:"chat_mode"`
		SenderName        string `json:"sender_name"`
		Message           string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("party consumer: party_chat JSON: %v", err)
		return nil
	}
	if payload.TargetCharacterID == 0 {
		return nil
	}
	mode := byte(payload.ChatMode)
	if payload.ChatMode > 255 {
		mode = 255
	}
	gs.EnsureSend(nil, payload.TargetCharacterID, &g_actor.DeliverPartyMultiChat{
		CharacterID: payload.TargetCharacterID,
		Mode:        mode,
		SenderName:  payload.SenderName,
		Message:     payload.Message,
	})
	return nil
}
