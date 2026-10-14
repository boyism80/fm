package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/domain"
	"github.com/boyism80/fm/stream"
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
		SourceChannelID uint32   `json:"source_channel_id"`
		MessageType     uint32   `json:"message_type"`
		Message         string   `json:"message"`
		MegaEar         bool     `json:"mega_ear"`
		NextLines       []string `json:"next_lines"`
		Item            []byte   `json:"item"`
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
	switch messageType {
	case constant.MsgThreeMegaphoneLines, constant.MsgItemMegaphone:
		megaphone := &response.Megaphone{
			Type:    messageType,
			Lines:   append([]string{payload.Message}, payload.NextLines...),
			Channel: channel,
			Ear:     payload.MegaEar,
		}
		if len(payload.Item) > 0 {
			megaphone.Item = dto.NewItemFromStream(stream.NewStreamReader(&payload.Item, stream.LittleEndian))
		}
		h.gs.systems.Get[*domain.NoticeSystem]().BroadcastMegaphone(megaphone)
	default:
		h.gs.systems.Get[*domain.NoticeSystem]().BroadcastNotice(messageType, payload.Message, channel, payload.MegaEar)
	}
	return nil
}
