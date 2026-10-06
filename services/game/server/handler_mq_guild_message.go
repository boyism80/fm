package server

import (
	"encoding/json"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type guildMqMessage struct{ gs *GameServer }

func (guildMqMessage) New(gs *GameServer) *guildMqMessage {
	return &guildMqMessage{gs: gs}
}

func (*guildMqMessage) EventType() string {
	return "message"
}

func (h *guildMqMessage) Handle(_ actor.Context, raw json.RawMessage) error {
	gs := h.gs
	if gs == nil {
		return nil
	}
	var payload struct {
		GuildID     uint32 `json:"guild_id"`
		MessageType uint32 `json:"message_type"`
		Message     string `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		log.Printf("guild consumer: message JSON: %v", err)
		return nil
	}
	if payload.GuildID == 0 || payload.Message == "" {
		return nil
	}
	gs.guild.BroadcastMessage(payload.GuildID, constant.ServerMessageType(payload.MessageType), payload.Message)
	return nil
}
