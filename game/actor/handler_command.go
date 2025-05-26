package actor

import (
	"log"
	"math"
	"strconv"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
)

func RegisterGameClientCommandHandler(ctx actor.Context, client *GameClientActor, h *handler.CommandHandler) {
	handler.RegisterCommandHandler("아이템", ctx, client, h, onCreateItem)
	handler.RegisterCommandHandler("메소초기화", ctx, client, h, onClearMeso)
	handler.RegisterCommandHandler("메소얻기", ctx, client, h, onGainMeso)
	handler.RegisterCommandHandler("풀메소", ctx, client, h, onFullMeso)
}

func onCreateItem(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onCreateItem: missing itemId")
		return
	}

	itemId, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onCreateItem: invalid itemId:", params[0])
		return
	}

	count := 1
	if len(params) >= 2 {
		if parsedCount, err := strconv.Atoi(params[1]); err == nil {
			count = parsedCount
		} else {
			log.Println("onCreateItem: invalid count, using default 1:", params[1])
		}
	}

	item, err := entity.NewItem(client.ctx, uint32(itemId), uint16(count))
	if err != nil {
		log.Println(err)
		return
	}

	inventoryType := item.GetInventoryType()
	inventory := client.ch.Inventory[inventoryType]
	nextSlot, ok := inventory.NextSlot()
	if !ok {
		return
	}

	inventory.Items[int16(nextSlot)] = item
	client.Send(&resp.AddItem{
		IsDrop:        false,
		Slot:          nextSlot,
		InventoryType: inventoryType,
		Item:          item,
	}, types.SEND_POLICY_ENCRYPT)
}

func onClearMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	client.ch.Meso = 0
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onGainMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	if len(params) < 1 {
		log.Println("onGainMeso: missing itemId")
		return
	}

	value, err := strconv.Atoi(params[0])
	if err != nil {
		log.Println("onGainMeso: invalid itemId:", params[0])
		return
	}
	cap := math.MaxInt32 - client.ch.Meso
	meso := min(cap, int32(value))
	client.ch.Meso += meso
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}

func onFullMeso(ctx actor.Context, client *GameClientActor, params ...string) {
	client.ch.Meso = math.MaxInt32
	client.Send(&resp.UpdateStats{
		Stats: map[constant.Stat]int32{
			constant.StatMeso: client.ch.Meso,
		},
	}, types.SEND_POLICY_ENCRYPT)
}
