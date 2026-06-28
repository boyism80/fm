package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqLeaderChanged struct {
	partyMqHandler
}

func (partyMqLeaderChanged) New(gs *GameServer) *partyMqLeaderChanged {
	return &partyMqLeaderChanged{partyMqHandler: partyMqHandler{gs: gs}}
}

func (*partyMqLeaderChanged) EventType() string {
	return "leader_changed"
}

func (h *partyMqLeaderChanged) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	pc.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		if raw == nil {
			return nil, nil
		}
		var extra struct {
			NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || extra.NewLeaderCharacterID == 0 {
			return nil, nil
		}
		party := pc.Get(evt.PartyID)
		if party != nil {
			h.sendLeaderChange(party, extra.NewLeaderCharacterID, false)
		}
		return nil, nil
	})
	return nil
}
