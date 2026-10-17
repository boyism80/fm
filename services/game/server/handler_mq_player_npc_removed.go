package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type playerNpcMqRemoved struct{ gs *GameServer }

func (playerNpcMqRemoved) New(gs *GameServer) *playerNpcMqRemoved {
	return &playerNpcMqRemoved{gs: gs}
}

func (*playerNpcMqRemoved) EventType() string {
	return "player_npc_removed"
}

func (h *playerNpcMqRemoved) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		NpcID uint32 `json:"npc_id"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	h.gs.systems.Get[*game.PlayerNpcSystem]().Remove(payload.NpcID)
	return nil
}
