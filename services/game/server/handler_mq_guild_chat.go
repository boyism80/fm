package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	pconst "github.com/boyism80/fm/protocol/constant"
	amqp "github.com/rabbitmq/amqp091-go"
)

type guildMqChat struct {
	gs *GameServer
}

func (guildMqChat) New(gs *GameServer) *guildMqChat {
	return &guildMqChat{gs: gs}
}

func (*guildMqChat) EventType() string {
	return "multi_chat"
}

func (h *guildMqChat) Handle(_ actor.Context, _ amqp.Delivery, _ string, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		GuildID           uint32 `json:"guild_id"`
		SenderCharacterID uint32 `json:"sender_character_id"`
		ChatMode          uint32 `json:"chat_mode"`
		SenderName        string `json:"sender_name"`
		Message           string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("guild consumer: multi_chat JSON: %v", err)
		return nil
	}
	if payload.GuildID == 0 || payload.SenderCharacterID == 0 {
		return nil
	}
	if payload.ChatMode != uint32(pconst.MultiChatModeGuild) {
		return nil
	}
	gs.guild.BroadcastMultiChat(
		payload.GuildID,
		payload.SenderCharacterID,
		payload.SenderName,
		payload.Message,
	)
	return nil
}
