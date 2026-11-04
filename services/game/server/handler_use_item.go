package server

import (
	"errors"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type UseItem struct {
	gs *GameServer
}

func (UseItem) New(gs *GameServer) *UseItem {
	return &UseItem{gs: gs}
}

func (h *UseItem) Handle(ctx *core.ClientContext, req *request.UseItem) error {
	ch, err := h.gs.character(ctx)
	if err != nil {
		return nil
	}
	if ch.UsingCashItem() {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	mapInstance := ch.GetMap()
	if mapInstance != nil && mapInstance.Wz().Limits(constant.FieldLimitPotion) {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	item := ch.Inventory().GetItem(constant.InventoryTypeConsume, int16(req.Slot))
	if item == nil {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	if item.GetCount() < 1 {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	if item.GetModel().GetID() != req.ItemID {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	consumeEnt, ok := item.(*game.Consume)
	if !ok || consumeEnt == nil {
		ch.Listener().OnUpdateStats(ch, nil, true)
		return nil
	}

	if errors.Is(ch.UseConsumeInSlot(int16(req.Slot), consumeEnt), game.ErrConsumeCooldown) {
		ch.Listener().OnMessage(ch, constant.MsgPinkText, "아직 아이템을 사용할 수 없습니다.")
	}
	ch.Listener().OnUpdateStats(ch, nil, true)

	return nil
}
