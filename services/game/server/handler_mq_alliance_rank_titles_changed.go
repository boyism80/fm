package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
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
	if !gs.systems.Get[*game.AllianceSystem]().Update(alliancePb) {
		return nil
	}
	gs.systems.Get[*game.AllianceSystem]().BroadcastInfoUpdate(alliancePb)
	return nil
}
