package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"

	"github.com/boyism80/fm/services/game"
)

type partyMqChat struct{ gs *GameServer }

func (partyMqChat) New(gs *GameServer) *partyMqChat {
	return &partyMqChat{gs: gs}
}

func (*partyMqChat) EventType() string { return "multi_chat" }

func (h *partyMqChat) Handle(_ actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		MemberID          uint32 `json:"member_id"`
		SenderCharacterID uint32 `json:"sender_character_id"`
		ChatMode          uint32 `json:"chat_mode"`
		SenderName        string `json:"sender_name"`
		Message           string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.MemberID == 0 {
		return nil
	}
	if payload.ChatMode > uint32(pconst.MultiChatModeAlliance) {
		return nil
	}
	mode := pconst.MultiChatMode(payload.ChatMode)
	switch mode {
	case pconst.MultiChatModeParty:
		h.handlePartyMultiChat(gs, payload.MemberID, payload.SenderCharacterID, mode, payload.SenderName, payload.Message)
	default:
		return nil
	}
	return nil
}

func (h *partyMqChat) handlePartyMultiChat(gs *GameServer, memberID uint32, senderCharacterID uint32, mode pconst.MultiChatMode, senderName string, message string) {
	party := gs.systems.Get[*game.PartySystem]().Get(memberID)
	if party == nil {
		return
	}
	for _, member := range party.GetMembers() {
		if member == nil {
			continue
		}
		cid := member.GetCharacterId()
		if cid == 0 || cid == senderCharacterID {
			continue
		}
		gs.systems.Get[*game.DispatchSystem]().Visit(cid, func(_ actor.Context, ch *game.Character) {
			if ch == nil {
				return
			}
			ch.Listener().OnMultiChat(ch, mode, senderName, message)
		})
	}
}
