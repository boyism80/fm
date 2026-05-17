package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type buddyMqChat struct {
	gs *GameServer
}

func (buddyMqChat) New(gs *GameServer) *buddyMqChat {
	return &buddyMqChat{gs: gs}
}

func (*buddyMqChat) EventType() string {
	return "multi_chat"
}

func (h *buddyMqChat) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		SenderCharacterID  uint32   `json:"sender_character_id"`
		ChatMode           uint32   `json:"chat_mode"`
		SenderName         string   `json:"sender_name"`
		Message            string   `json:"message"`
		NotifyCharacterIDs []uint32 `json:"notify_character_ids"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("buddy consumer: multi_chat JSON: %v", err)
		return nil
	}
	if payload.SenderCharacterID == 0 || len(payload.NotifyCharacterIDs) == 0 {
		return nil
	}
	if payload.ChatMode > uint32(pconst.MultiChatModeAlliance) {
		return nil
	}
	mode := pconst.MultiChatMode(payload.ChatMode)
	if mode != pconst.MultiChatModeBuddy {
		return nil
	}
	for _, recipientID := range payload.NotifyCharacterIDs {
		if recipientID == 0 || recipientID == payload.SenderCharacterID {
			continue
		}
		if gs.characterRuntime != nil && !gs.characterRuntime.Exists(recipientID) {
			continue
		}
		gs.EnsureSend(nil, recipientID, &g_actor.DeliverMultiChat{
			CharacterID: recipientID,
			Mode:        mode,
			SenderName:  payload.SenderName,
			Message:     payload.Message,
		})
	}
	return nil
}
