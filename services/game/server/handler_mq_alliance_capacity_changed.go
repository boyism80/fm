package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqCapacityChanged struct{ gs *GameServer }

func (allianceMqCapacityChanged) New(gs *GameServer) *allianceMqCapacityChanged {
	return &allianceMqCapacityChanged{gs: gs}
}

func (*allianceMqCapacityChanged) EventType() string {
	return "capacity_changed"
}

func (h *allianceMqCapacityChanged) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
	gs.alliance.Update(alliancePb)
	gs.alliance.BroadcastCapacityChanged(alliancePb)
	log.Printf("alliance consumer: applied capacity_changed alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}
