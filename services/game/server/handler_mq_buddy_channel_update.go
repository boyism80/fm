package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type buddyMqChannelUpdate struct {
	gs *GameServer
}

func (buddyMqChannelUpdate) New(gs *GameServer) *buddyMqChannelUpdate {
	return &buddyMqChannelUpdate{gs: gs}
}

func (*buddyMqChannelUpdate) EventType() string {
	return "channel_update"
}

func (h *buddyMqChannelUpdate) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		CharacterID        uint32   `json:"character_id"`
		ChannelIndex       int32    `json:"channel_index"`
		NotifyCharacterIDs []uint32 `json:"notify_character_ids"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("buddy consumer: channel_update JSON: %v", err)
		return nil
	}
	if payload.CharacterID == 0 || len(payload.NotifyCharacterIDs) == 0 {
		return nil
	}
	for _, recipientID := range payload.NotifyCharacterIDs {
		if recipientID == 0 {
			continue
		}
		if gs.characterRuntime != nil && !gs.characterRuntime.Exists(recipientID) {
			continue
		}
		gs.EnsureSend(nil, recipientID, &g_actor.DeliverBuddyChannelUpdate{
			RecipientCharacterID: recipientID,
			BuddyCharacterID:     payload.CharacterID,
			Channel:              payload.ChannelIndex,
		})
	}
	return nil
}
