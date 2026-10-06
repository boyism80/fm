package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

type allianceMqRankTitlesChanged struct{ gs *GameServer }

func (allianceMqRankTitlesChanged) New(gs *GameServer) *allianceMqRankTitlesChanged {
	return &allianceMqRankTitlesChanged{gs: gs}
}

func (*allianceMqRankTitlesChanged) EventType() string {
	return "rank_titles_changed"
}

func (h *allianceMqRankTitlesChanged) Handle(_ actor.Context, raw json.RawMessage) error {
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
	gs.alliance.BroadcastInfoUpdate(alliancePb)
	return nil
}
