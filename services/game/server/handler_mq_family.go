package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
)

type familyMq struct{ gs *GameServer }

func (familyMq) New(gs *GameServer) *familyMq {
	return &familyMq{gs: gs}
}

func (*familyMq) EventType() string {
	return "family"
}

func (h *familyMq) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32               `json:"character_id"`
		Event       internal.FamilyEvent `json:"event"`
		SenderID    uint32               `json:"sender_id"`
		Name        string               `json:"name"`
		MapName     string               `json:"map_name"`
		Accepted    bool                 `json:"accepted"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.Deliver(nil, payload.CharacterID, &g_actor.DeliverFamily{
		CharacterID: payload.CharacterID,
		Event:       payload.Event,
		SenderID:    payload.SenderID,
		Name:        payload.Name,
		MapName:     payload.MapName,
		Accepted:    payload.Accepted,
	})
	return nil
}
