package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

type marriageMqSpouseMoved struct{ gs *GameServer }

func (marriageMqSpouseMoved) New(gs *GameServer) *marriageMqSpouseMoved {
	return &marriageMqSpouseMoved{gs: gs}
}

func (*marriageMqSpouseMoved) EventType() string {
	return "spouse_moved"
}

func (h *marriageMqSpouseMoved) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32 `json:"character_id"`
		SpouseID    uint32 `json:"spouse_id"`
		MapID       uint32 `json:"map_id"`
		Reply       bool   `json:"reply"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.EnsureSend(nil, payload.CharacterID, &g_actor.DeliverSpouseMoved{
		CharacterID: payload.CharacterID,
		SpouseID:    payload.SpouseID,
		MapID:       payload.MapID,
		Reply:       payload.Reply,
	})
	return nil
}
