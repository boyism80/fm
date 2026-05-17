package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/response"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type buddyMqListUpdate struct {
	gs *GameServer
}

func (buddyMqListUpdate) New(gs *GameServer) *buddyMqListUpdate {
	return &buddyMqListUpdate{gs: gs}
}

func (*buddyMqListUpdate) EventType() string {
	return "list_update"
}

func (h *buddyMqListUpdate) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		SyncAction         uint8    `json:"sync_action"`
		NotifyCharacterIDs []uint32 `json:"notify_character_ids"`
		Entries            []struct {
			CharacterID  uint32 `json:"character_id"`
			Name         string `json:"name"`
			GroupName    string `json:"group_name"`
			Pending      bool   `json:"pending"`
			ChannelIndex int32  `json:"channel_index"`
		} `json:"entries"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("buddy consumer: list_update JSON: %v", err)
		return nil
	}
	if len(payload.NotifyCharacterIDs) == 0 || len(payload.Entries) == 0 {
		return nil
	}
	action := pconst.BuddyListSyncAction(payload.SyncAction)
	if action == 0 {
		action = pconst.BuddyListSyncUpdate
	}
	entries := make([]response.BuddyEntry, 0, len(payload.Entries))
	for _, row := range payload.Entries {
		if row.CharacterID == 0 {
			continue
		}
		ch := row.ChannelIndex
		if ch < 0 {
			ch = -1
		}
		entries = append(entries, response.BuddyEntry{
			CharacterID: row.CharacterID,
			Name:        row.Name,
			Pending:     row.Pending,
			Channel:     ch,
			Group:       row.GroupName,
		})
	}
	if len(entries) == 0 {
		return nil
	}
	for _, recipientID := range payload.NotifyCharacterIDs {
		if recipientID == 0 {
			continue
		}
		gs.DeliverBuddyListUpdate(recipientID, uint8(action), entries)
	}
	return nil
}

func (gs *GameServer) DeliverBuddyListUpdate(recipientCharacterID uint32, syncAction uint8, entries []response.BuddyEntry) {
	if gs == nil || recipientCharacterID == 0 || len(entries) == 0 {
		return
	}
	if gs.characterRuntime != nil && !gs.characterRuntime.Exists(recipientCharacterID) {
		return
	}
	gs.EnsureSend(nil, recipientCharacterID, &g_actor.DeliverBuddyListUpdate{
		RecipientCharacterID: recipientCharacterID,
		SyncAction:           syncAction,
		Entries:              entries,
	})
}
