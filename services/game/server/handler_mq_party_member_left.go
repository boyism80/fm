package server

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type partyMqMemberLeft struct{ gs *GameServer }

func (partyMqMemberLeft) New(gs *GameServer) *partyMqMemberLeft {
	return &partyMqMemberLeft{gs: gs}
}
func (*partyMqMemberLeft) EventType() string { return "member_left" }
func (h *partyMqMemberLeft) Handle(_ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("party consumer: apply type=member_left party_id=%d: %v", evt.PartyID, err)
		return nil
	}
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
	snapshot := pc.CachedSnapshot(evt.PartyID)
	pc.DeliverPartyLeaveUpdate(prevSnapshot, snapshot, extra.CharacterID, extra.ExpelledByCharacter != 0)
	if extra.LeaderChanged && extra.NewLeaderCharacterID != 0 && snapshot != nil {
		pc.DeliverPartyLeaderChange(snapshot, extra.NewLeaderCharacterID, true)
	}
	return nil
}
