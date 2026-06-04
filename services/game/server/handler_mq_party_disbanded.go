package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqDisbanded struct {
	partyMqHandler
}

func (partyMqDisbanded) New(gs *GameServer) *partyMqDisbanded {
	return &partyMqDisbanded{partyMqHandler: partyMqHandler{gs: gs}}
}

func (*partyMqDisbanded) EventType() string {
	return "disbanded"
}

func (h *partyMqDisbanded) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevParty := pc.Get(evt.PartyID)
	pc.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) {
			return nil, nil
		}, func(interface{}) error {
			if raw == nil {
				return nil
			}
			var extra struct {
				CharacterID uint32 `json:"character_id"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 || prevParty == nil {
				return nil
			}
			h.sendDisbandUpdate(prevParty, extra.CharacterID)
			return nil
		}).
		Run()
	return nil
}
