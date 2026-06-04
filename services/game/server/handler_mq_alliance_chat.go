package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	g_actor "github.com/boyism80/fm/services/game/actor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type allianceMqChat struct {
	gs *GameServer
}

func (allianceMqChat) New(gs *GameServer) *allianceMqChat {
	return &allianceMqChat{gs: gs}
}

func (*allianceMqChat) EventType() string {
	return "multi_chat"
}

func (h *allianceMqChat) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		AllianceID        uint32 `json:"alliance_id"`
		SenderCharacterID uint32 `json:"sender_character_id"`
		ChatMode          uint32 `json:"chat_mode"`
		SenderName        string `json:"sender_name"`
		Message           string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("alliance consumer: multi_chat JSON: %v", err)
		return nil
	}
	if payload.AllianceID == 0 || payload.SenderCharacterID == 0 {
		return nil
	}
	if payload.ChatMode != uint32(pconst.MultiChatModeAlliance) {
		return nil
	}
	mode := pconst.MultiChatModeAlliance
	guildIDs := gs.alliance.GuildIDs(payload.AllianceID)
	if len(guildIDs) == 0 {
		return nil
	}
	seen := make(map[uint32]struct{})
	for _, guildID := range guildIDs {
		g := gs.guild.Get(guildID)
		if g == nil {
			continue
		}
		for _, m := range g.GetMembers() {
			if m == nil {
				continue
			}
			memberID := m.GetCharacterId()
			if memberID == 0 || memberID == payload.SenderCharacterID {
				continue
			}
			if _, dup := seen[memberID]; dup {
				continue
			}
			seen[memberID] = struct{}{}
			if gs.characterRuntime == nil || !gs.characterRuntime.Exists(memberID) {
				continue
			}
			gs.EnsureSend(nil, memberID, &g_actor.DeliverMultiChat{
				CharacterID: memberID,
				Mode:        mode,
				SenderName:  payload.SenderName,
				Message:     payload.Message,
			})
		}
	}
	return nil
}
