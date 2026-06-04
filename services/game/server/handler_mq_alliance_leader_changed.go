package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqLeaderChanged struct{ gs *GameServer }

func (allianceMqLeaderChanged) New(gs *GameServer) *allianceMqLeaderChanged {
	return &allianceMqLeaderChanged{gs: gs}
}

func (*allianceMqLeaderChanged) EventType() string {
	return "leader_changed"
}

func (h *allianceMqLeaderChanged) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("alliance consumer: leader_changed alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	var extra struct {
		OldLeaderCharacterID uint32 `json:"old_leader_character_id"`
		NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.NewLeaderCharacterID == 0 {
		log.Printf("alliance consumer: leader_changed alliance_id=%d missing leader ids", evt.AllianceID)
		return nil
	}
	gs.alliance.Update(alliancePb)
	gs.alliance.BroadcastLeaderChanged(alliancePb, extra.OldLeaderCharacterID, extra.NewLeaderCharacterID)
	log.Printf("alliance consumer: applied leader_changed alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}
