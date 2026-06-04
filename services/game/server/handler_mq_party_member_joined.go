package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqMemberJoined struct {
	partyMqHandler
}

func (partyMqMemberJoined) New(gs *GameServer) *partyMqMemberJoined {
	return &partyMqMemberJoined{partyMqHandler: partyMqHandler{gs: gs}}
}

func (*partyMqMemberJoined) EventType() string {
	return "member_joined"
}

func (h *partyMqMemberJoined) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
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
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
				return nil
			}
			party := pc.Get(evt.PartyID)
			if party != nil {
				h.sendJoinUpdate(party, extra.CharacterID)
			}
			return nil
		}).
		Run()
	return nil
}
