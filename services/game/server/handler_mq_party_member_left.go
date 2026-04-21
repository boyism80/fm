package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqMemberLeft struct{ gs *GameServer }

func (partyMqMemberLeft) New(gs *GameServer) *partyMqMemberLeft {
	return &partyMqMemberLeft{gs: gs}
}
func (*partyMqMemberLeft) EventType() string { return "member_left" }
func (h *partyMqMemberLeft) Handle(ctx actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil || gs.party == nil {
		return nil
	}
	pc := gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevParty := pc.Get(evt.PartyID)
	pc.UpdateAsync(ctx, evt).
		Then(func() (interface{}, error) { return nil, nil }, func(interface{}) error {
			if raw == nil {
				return nil
			}
			var extra struct {
				CharacterID          uint32 `json:"character_id"`
				ExpelledByCharacter  uint32 `json:"expelled_by_character_id"`
				LeaderChanged        bool   `json:"leader_changed"`
				NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
			}
			if err := json.Unmarshal(raw, &extra); err != nil || extra.CharacterID == 0 {
				return nil
			}
			party := pc.Get(evt.PartyID)
			pc.DeliverPartyLeaveUpdate(prevParty, party, extra.CharacterID, extra.ExpelledByCharacter != 0)
			if extra.LeaderChanged && extra.NewLeaderCharacterID != 0 && party != nil {
				pc.DeliverPartyLeaderChange(party, extra.NewLeaderCharacterID, true)
			}
			return nil
		}).
		Run()
	return nil
}
