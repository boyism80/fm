package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqDisbanded struct{ gs *GameServer }

func (partyMqDisbanded) New(gs *GameServer) *partyMqDisbanded {
	return &partyMqDisbanded{gs: gs}
}
func (*partyMqDisbanded) EventType() string { return "disbanded" }
func (h *partyMqDisbanded) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
	pc.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) { return nil, nil }, func(interface{}) error {
			if raw == nil {
				return nil
			}
			var extra struct {
				CharacterID uint32 `json:"character_id"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 || prevSnapshot == nil {
				return nil
			}
			pc.DeliverPartyDisbandUpdate(prevSnapshot, extra.CharacterID)
			return nil
		}).
		Run()
	return nil
}
