package server

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqMemberJoined struct{ gs *GameServer }

func (partyMqMemberJoined) New(gs *GameServer) *partyMqMemberJoined {
	return &partyMqMemberJoined{gs: gs}
}
func (*partyMqMemberJoined) EventType() string { return "member_joined" }
func (h *partyMqMemberJoined) Handle(_ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	if err := pc.apply(evt); err != nil {
		log.Printf("party consumer: apply type=member_joined party_id=%d: %v", evt.PartyID, err)
		return nil
	}
	if raw == nil {
		return nil
	}
	var extra struct {
		CharacterID uint32 `json:"character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
		return nil
	}
	snapshot := pc.CachedSnapshot(evt.PartyID)
	if snapshot != nil {
		pc.DeliverPartyJoinUpdate(snapshot, extra.CharacterID)
	}
	return nil
}
