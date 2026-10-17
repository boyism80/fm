package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type allianceMqCapacityChanged struct{ gs *GameServer }

func (allianceMqCapacityChanged) New(gs *GameServer) *allianceMqCapacityChanged {
	return &allianceMqCapacityChanged{gs: gs}
}

func (*allianceMqCapacityChanged) EventType() string {
	return "capacity_changed"
}

func (h *allianceMqCapacityChanged) Handle(_ actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	evt, ok := decodeAllianceEventEnvelope(raw)
	if !ok {
		return nil
	}
	alliancePb, ok := allianceFromMQPayload(raw)
	if !ok {
		log.Printf("alliance consumer: capacity_changed alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	if !gs.systems.Get[*game.AllianceSystem]().Update(alliancePb) {
		return nil
	}
	gs.systems.Get[*game.AllianceSystem]().BroadcastInfoUpdate(alliancePb)
	return nil
}
