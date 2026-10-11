package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

type noteMqArrived struct{ gs *GameServer }

func (noteMqArrived) New(gs *GameServer) *noteMqArrived {
	return &noteMqArrived{gs: gs}
}

func (*noteMqArrived) EventType() string {
	return "note_arrived"
}

func (h *noteMqArrived) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32 `json:"character_id"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.EnsureSend(nil, payload.CharacterID, &g_actor.DeliverNoteArrived{
		CharacterID: payload.CharacterID,
	})
	return nil
}
