package server

import (
	"encoding/json"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	gameconst "github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/services/game"
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
	message := fmt.Sprintf("%s님이 캐시샵 선물을 보냈습니다. 캐시샵 보관함을 확인하세요.", payload.SenderName)
	h.gs.systems.Get[*game.DispatchSystem]().CallCharacter(payload.CharacterID, func(_ actor.Context, ch *game.Character) {
		if ch == nil {
			return
		}
		ch.Listener.OnMessage(ch, gameconst.MsgNotice, message)
	})
	return nil
}
