package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqMemberRankChanged struct{ gs *GameServer }

func (allianceMqMemberRankChanged) New(gs *GameServer) *allianceMqMemberRankChanged {
	return &allianceMqMemberRankChanged{gs: gs}
}

func (*allianceMqMemberRankChanged) EventType() string {
	return "member_rank_changed"
}

func (h *allianceMqMemberRankChanged) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("alliance consumer: member_rank_changed alliance_id=%d missing alliance_pb", evt.AllianceID)
		return nil
	}
	var extra struct {
		CharacterID     uint32 `json:"character_id"`
		NewAllianceRank uint32 `json:"new_alliance_rank"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
		log.Printf("alliance consumer: member_rank_changed alliance_id=%d missing character_id", evt.AllianceID)
		return nil
	}
	allianceRank := extra.NewAllianceRank
	if allianceRank == 0 {
		for _, g := range alliancePb.GetGuilds() {
			if g == nil {
				continue
			}
			for _, m := range g.GetMembers() {
				if m != nil && m.GetCharacterId() == extra.CharacterID && m.AllianceRank != nil {
					allianceRank = m.GetAllianceRank()
					break
				}
			}
			if allianceRank != 0 {
				break
			}
		}
	}
	if allianceRank == 0 {
		log.Printf("alliance consumer: member_rank_changed alliance_id=%d missing rank for character=%d", evt.AllianceID, extra.CharacterID)
		return nil
	}
	gs.alliance.Update(alliancePb)
	gs.alliance.BroadcastMemberRankChanged(alliancePb, extra.CharacterID, allianceRank)
	log.Printf("alliance consumer: applied member_rank_changed alliance_id=%d character=%d revision=%d",
		evt.AllianceID, extra.CharacterID, evt.Revision)
	return nil
}
