package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqRankTitlesChanged struct{ gs *GameServer }

func (allianceMqRankTitlesChanged) New(gs *GameServer) *allianceMqRankTitlesChanged {
	return &allianceMqRankTitlesChanged{gs: gs}
}

func (*allianceMqRankTitlesChanged) EventType() string {
	return "rank_titles_changed"
}

func (h *allianceMqRankTitlesChanged) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("alliance consumer: rank_titles_changed alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	gs.alliance.Update(alliancePb)
	gs.alliance.BroadcastRankTitlesChanged(alliancePb)
	log.Printf("alliance consumer: applied rank_titles_changed alliance_id=%d revision=%d", evt.AllianceID, evt.Revision)
	return nil
}
