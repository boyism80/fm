package server

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqDisbanded struct{ gs *GameServer }

func (partyMqDisbanded) New(gs *GameServer) *partyMqDisbanded {
	return &partyMqDisbanded{gs: gs}
}
func (*partyMqDisbanded) EventType() string { return "disbanded" }
func (h *partyMqDisbanded) Handle(_ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevSnapshot := pc.CachedSnapshot(evt.PartyID)
	if err := pc.apply(evt); err != nil {
		log.Printf("party consumer: apply type=disbanded party_id=%d: %v", evt.PartyID, err)
		return nil
	}
	if raw == nil {
		return nil
	}
	var extra struct {
		CharacterID uint32 `json:"character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 || prevSnapshot == nil {
		return nil
	}
	gs.DeliverPartyDisbandUpdate(prevSnapshot, extra.CharacterID)
	return nil
}
