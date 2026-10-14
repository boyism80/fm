package server

import (
	"encoding/base64"
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/entity"
)

type messengerMq struct{ gs *GameServer }

func (messengerMq) New(gs *GameServer) *messengerMq {
	return &messengerMq{gs: gs}
}

func (*messengerMq) EventType() string {
	return "messenger"
}

func (h *messengerMq) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32                   `json:"character_id"`
		Action      internal.MessengerAction `json:"action"`
		RoomID      uint32                   `json:"room_id"`
		Position    uint8                    `json:"position"`
		Name        string                   `json:"name"`
		Text        string                   `json:"text"`
		Member      struct {
			Name      string `json:"name"`
			ChannelID uint16 `json:"channelId"`
			Position  uint8  `json:"position"`
			Look      string `json:"look"`
		} `json:"member"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}

	event := entity.MessengerEvent{
		Action:   payload.Action,
		RoomID:   payload.RoomID,
		Position: payload.Position,
		Name:     payload.Name,
		Text:     payload.Text,
	}
	if event.Action == internal.MessengerAction_MESSENGER_ACTION_ENTER || event.Action == internal.MessengerAction_MESSENGER_ACTION_AVATAR {
		look, err := base64.StdEncoding.DecodeString(payload.Member.Look)
		if err != nil {
			return nil
		}
		event.Look = look
		event.Name = payload.Member.Name
		event.Position = payload.Member.Position
		event.Channel = payload.Member.ChannelID
	}

	h.gs.SendToCharacter(nil, payload.CharacterID, &g_actor.ReceiveMessengerEvent{
		CharacterID: payload.CharacterID,
		Event:       event,
	})
	return nil
}
