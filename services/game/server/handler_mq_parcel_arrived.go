package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

type parcelMqArrived struct{ gs *GameServer }

func (parcelMqArrived) New(gs *GameServer) *parcelMqArrived {
	return &parcelMqArrived{gs: gs}
}

func (*parcelMqArrived) EventType() string {
	return "parcel_arrived"
}

func (h *parcelMqArrived) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32 `json:"character_id"`
		SenderName  string `json:"sender_name"`
		Quick       bool   `json:"quick"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.EnsureSend(nil, payload.CharacterID, &g_actor.DeliverParcelArrived{
		CharacterID: payload.CharacterID,
		SenderName:  payload.SenderName,
		Quick:       payload.Quick,
	})
	return nil
}
