package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	amqp "github.com/rabbitmq/amqp091-go"

	g_actor "github.com/boyism80/fm/services/game/actor"
)

type partyMqMultiChat struct{ gs *GameServer }

func (partyMqMultiChat) New(gs *GameServer) *partyMqMultiChat {
	return &partyMqMultiChat{gs: gs}
}

func (*partyMqMultiChat) EventType() string { return "multi_chat" }

func (h *partyMqMultiChat) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
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
		log.Printf("party consumer: multi_chat JSON: %v", err)
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

func (h *partyMqMultiChat) handlePartyMultiChat(gs *GameServer, memberID uint32, senderCharacterID uint32, mode pconst.MultiChatMode, senderName string, message string) {
	party := gs.GetPartyByID(memberID)
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
		gs.EnsureSend(nil, cid, &g_actor.DeliverMultiChat{
			CharacterID: cid,
			Mode:        mode,
			SenderName:  senderName,
			Message:     message,
		})
	}
}
