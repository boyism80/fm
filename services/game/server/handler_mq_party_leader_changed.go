package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqLeaderChanged struct{ gs *GameServer }

func (partyMqLeaderChanged) New(gs *GameServer) *partyMqLeaderChanged {
	return &partyMqLeaderChanged{gs: gs}
}
func (*partyMqLeaderChanged) EventType() string { return "leader_changed" }
func (h *partyMqLeaderChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	pc.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) { return nil, nil }, func(interface{}) error {
			if raw == nil {
				return nil
			}
			var extra struct {
				NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.NewLeaderCharacterID == 0 {
				return nil
			}
			snapshot := pc.CachedSnapshot(evt.PartyID)
			if snapshot != nil {
				pc.DeliverPartyLeaderChange(snapshot, extra.NewLeaderCharacterID, false)
			}
			return nil
		}).
		Run()
	return nil
}
