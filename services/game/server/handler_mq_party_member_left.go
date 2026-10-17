package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game"
)

type partyMqMemberLeft struct{ gs *GameServer }

func (partyMqMemberLeft) New(gs *GameServer) *partyMqMemberLeft {
	return &partyMqMemberLeft{gs: gs}
}

func (*partyMqMemberLeft) EventType() string {
	return "member_left"
}

func (h *partyMqMemberLeft) Handle(ctx actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	pc := h.gs.systems.Get[*game.PartySystem]()
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	var extra struct {
		CharacterID          uint32 `json:"character_id"`
		ExpelledByCharacter  uint32 `json:"expelled_by_character_id"`
		LeaderChanged        bool   `json:"leader_changed"`
		NewLeaderCharacterID uint32 `json:"new_leader_character_id"`
	}
	if err := json.Unmarshal(raw, &extra); err != nil {
		return nil
	}
	partyPb := decodePartySnapshot(evt.PartyPB)
	if partyPb == nil {
		return nil
	}
	h.gs.systems.Get[*game.DispatchSystem]().SendToCharacter(nil, extra.CharacterID, &game.ClearPartyByPartyID{
		CharacterID: extra.CharacterID,
		PartyID:     evt.PartyID,
		Revision:    partyPb.GetRevision(),
	})

	prevParty := pc.Get(evt.PartyID)
	if !pc.Update(partyPb) {
		return nil
	}
	party := pc.Get(evt.PartyID)
	pc.BroadcastMemberLeft(prevParty, party, extra.CharacterID, extra.ExpelledByCharacter != 0)
	if extra.LeaderChanged && extra.NewLeaderCharacterID != 0 && party != nil {
		pc.BroadcastLeaderChanged(party, extra.NewLeaderCharacterID, true)
	}
	return nil
}
