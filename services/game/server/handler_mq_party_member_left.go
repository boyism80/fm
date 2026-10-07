package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
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
	pc := h.gs.party
	evt, ok := decodePartyEventEnvelope(raw)
	if !ok {
		return nil
	}
	prevParty := pc.Get(evt.PartyID)
	pc.UpdateAsync(ctx, evt).Do(func() error {
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
		h.gs.EnsureSend(nil, extra.CharacterID, &g_actor.ClearPartyByPartyID{
			CharacterID: extra.CharacterID,
			PartyID:     evt.PartyID,
		})
		party := pc.Get(evt.PartyID)
		pc.BroadcastMemberLeft(prevParty, party, extra.CharacterID, extra.ExpelledByCharacter != 0)
		if extra.LeaderChanged && extra.NewLeaderCharacterID != 0 && party != nil {
			pc.BroadcastLeaderChanged(party, extra.NewLeaderCharacterID, true)
		}
		return nil
	})
	return nil
}
