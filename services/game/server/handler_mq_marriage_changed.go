package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type marriageMqChanged struct{ gs *GameServer }

func (marriageMqChanged) New(gs *GameServer) *marriageMqChanged {
	return &marriageMqChanged{gs: gs}
}

func (*marriageMqChanged) EventType() string {
	return "marriage_changed"
}

func (h *marriageMqChanged) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32 `json:"character_id"`
		Event       string `json:"event"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.systems.Get[*game.DispatchSystem]().SendToCharacter(nil, payload.CharacterID, &game.RefreshWedding{
		CharacterID: payload.CharacterID,
		Event:       payload.Event,
	})
	return nil
}
