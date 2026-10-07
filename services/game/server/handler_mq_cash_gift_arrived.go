package server

import (
	"encoding/json"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	g_actor "github.com/boyism80/fm/services/game/actor"
	gameconst "github.com/boyism80/fm/services/game/constant"
)

type cashGiftMqArrived struct{ gs *GameServer }

func (cashGiftMqArrived) New(gs *GameServer) *cashGiftMqArrived {
	return &cashGiftMqArrived{gs: gs}
}

func (*cashGiftMqArrived) EventType() string {
	return "cash_gift_arrived"
}

func (h *cashGiftMqArrived) Handle(_ actor.Context, raw json.RawMessage) error {
	var payload struct {
		CharacterID uint32 `json:"character_id"`
		SenderName  string `json:"sender_name"`
	}
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil
	}
	if payload.CharacterID == 0 {
		return nil
	}
	h.gs.EnsureSend(nil, payload.CharacterID, &g_actor.DeliverMessage{
		CharacterID: payload.CharacterID,
		MessageType: gameconst.MsgNotice,
		Message:     fmt.Sprintf("%s님이 캐시샵 선물을 보냈습니다. 캐시샵 보관함을 확인하세요.", payload.SenderName),
	})
	return nil
}
