package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/services/game/constant"
)

type globalMqNotice struct {
	gs *GameServer
}

func (globalMqNotice) New(gs *GameServer) *globalMqNotice {
	return &globalMqNotice{gs: gs}
}

func (*globalMqNotice) EventType() string {
	return "notice"
}

func (h *globalMqNotice) Handle(_ actor.Context, raw json.RawMessage) error {
	if h.gs == nil {
		return nil
	}
	var payload struct {
		SourceChannelID uint32 `json:"source_channel_id"`
		MessageType     uint32 `json:"message_type"`
		Message         string `json:"message"`
		MegaEar         bool   `json:"mega_ear"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.Message == "" {
		return nil
	}
	if payload.SourceChannelID == h.gs.config.ChannelId {
		return nil
	}
	messageType := constant.ServerMessageType(payload.MessageType)
	channel := int(payload.SourceChannelID) + 1
	h.gs.BroadcastNotice(messageType, payload.Message, channel, payload.MegaEar)
	return nil
}
