package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqCreated struct{ gs *GameServer }

func (partyMqCreated) New(gs *GameServer) *partyMqCreated { return &partyMqCreated{gs: gs} }
func (*partyMqCreated) EventType() string                 { return "created" }
func (h *partyMqCreated) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	pc := gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	pc.UpdateAsync(ctx, evt).Run()
	return nil
}
