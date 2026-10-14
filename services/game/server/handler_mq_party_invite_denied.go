package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/domain"
)

type partyMqPartyInviteDenied struct{ gs *GameServer }

func (partyMqPartyInviteDenied) New(gs *GameServer) *partyMqPartyInviteDenied {
	return &partyMqPartyInviteDenied{gs: gs}
}

func (*partyMqPartyInviteDenied) EventType() string {
	return "party_invite_denied"
}

func (h *partyMqPartyInviteDenied) Handle(_ actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	var payload struct {
		InviterCharacterID  uint32 `json:"inviter_character_id"`
		DeniedCharacterName string `json:"denied_character_name"`
		Action              uint8  `json:"action"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.InviterCharacterID == 0 {
		return nil
	}
	h.gs.systems.Get[*domain.DispatchSystem]().CallCharacter(payload.InviterCharacterID, func(_ actor.Context, ch *domain.Character) {
		if ch == nil {
			return
		}
		ch.Listener.OnPartyStatusMessage(ch, pconst.PartyStatusCode(payload.Action), payload.DeniedCharacterName)
	})
	return nil
}
