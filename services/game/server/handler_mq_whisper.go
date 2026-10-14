package server

import (
	"encoding/json"

	"github.com/asynkron/protoactor-go/actor"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	g_actor "github.com/boyism80/fm/services/game/actor"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
)

type whisperMq struct{ gs *GameServer }

func (whisperMq) New(gs *GameServer) *whisperMq {
	return &whisperMq{gs: gs}
}

func (*whisperMq) EventType() string {
	return "whisper"
}

func (h *whisperMq) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32                 `json:"character_id"`
		RequestID   string                 `json:"request_id"`
		Mode        internal.WhisperMode   `json:"mode"`
		SenderName  string                 `json:"sender_name"`
		SenderRole  constant.CharacterRole `json:"sender_role"`
		ChannelID   int32                  `json:"channel_id"`
		Message     string                 `json:"message"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.SendToCharacter(nil, payload.CharacterID, &g_actor.ReceiveWhisper{
		CharacterID: payload.CharacterID,
		Whisper: entity.Whisper{
			RequestID:  payload.RequestID,
			Mode:       payload.Mode,
			SenderName: payload.SenderName,
			SenderRole: payload.SenderRole,
			Channel:    payload.ChannelID,
			Message:    payload.Message,
		},
	})
	return nil
}
