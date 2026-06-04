package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type buddyMqAddRequest struct {
	gs *GameServer
}

func (buddyMqAddRequest) New(gs *GameServer) *buddyMqAddRequest {
	return &buddyMqAddRequest{gs: gs}
}

func (*buddyMqAddRequest) EventType() string {
	return "add_request"
}

func (h *buddyMqAddRequest) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		NotifyCharacterIDs []uint32 `json:"notify_character_ids"`
		FromCharacterID    uint32   `json:"from_character_id"`
		FromName           string   `json:"from_name"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("buddy consumer: add_request JSON: %v", err)
		return nil
	}
	if payload.FromCharacterID == 0 {
		return nil
	}
	ids := payload.NotifyCharacterIDs
	if len(ids) == 0 {
		return nil
	}
	for _, recipientID := range ids {
		if recipientID == 0 {
			continue
		}
		if gs.characterRuntime != nil && !gs.characterRuntime.Exists(recipientID) {
			continue
		}
		gs.EnsureSend(nil, recipientID, &g_actor.DeliverBuddyAddRequest{
			RecipientCharacterID: recipientID,
			FromCharacterID:      payload.FromCharacterID,
			FromName:             payload.FromName,
		})
	}
	return nil
}
