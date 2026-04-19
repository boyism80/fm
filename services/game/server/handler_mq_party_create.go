package server

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqCreated struct{ gs *GameServer }

func (partyMqCreated) New(gs *GameServer) *partyMqCreated { return &partyMqCreated{gs: gs} }
func (*partyMqCreated) EventType() string                 { return "created" }
func (h *partyMqCreated) Handle(_ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("party consumer: apply type=created party_id=%d: %v", evt.PartyID, err)
	}
	return nil
}
