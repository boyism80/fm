package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqMemberLeft struct {
	partyMqHandler
}

func (partyMqMemberLeft) New(gs *GameServer) *partyMqMemberLeft {
	return &partyMqMemberLeft{partyMqHandler: partyMqHandler{gs: gs}}
}

func (*partyMqMemberLeft) EventType() string {
	return "member_left"
}

func (h *partyMqMemberLeft) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevParty := pc.Get(evt.PartyID)
	pc.UpdateAsync(ctx, evt).Then(func(interface{}) (interface{}, error) {
		if raw == nil {
			return nil, nil
		}
		var extra struct {
			CharacterID          uint32 `json:"character_id"`
			ExpelledByCharacter  uint32 `json:"expelled_by_character_id"`
			LeaderChanged        bool   `json:"leader_changed"`
			NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
		}
		if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
			return nil, nil
		}
		party := pc.Get(evt.PartyID)
		h.sendLeaveUpdate(prevParty, party, extra.CharacterID, extra.ExpelledByCharacter != 0)
		if extra.LeaderChanged && extra.NewLeaderCharacterID != 0 && party != nil {
			h.sendLeaderChange(party, extra.NewLeaderCharacterID, true)
		}
		return nil, nil
	})
	return nil
}
