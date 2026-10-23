package server

import (
	"errors"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type UseItem struct{}

func (UseItem) New(_ *GameServer) *UseItem {
	return &UseItem{}
}

func (*UseItem) Handle(ctx *core.ClientContext, req *request.UseItem) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		log.Printf("Client is not a GameClient")
		return fmt.Errorf("client is not a GameClient")
	}

	ch := client.GetCharacter()
	if ch == nil {
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
