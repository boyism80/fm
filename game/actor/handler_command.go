package actor

import (
	"log"
	"strconv"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/common/handler"
	"github.com/boyism80/fm/common/types"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/game/protocol/resp"
)

func RegisterGameClientCommandHandler(ctx actor.Context, client *GameClientActor, h *handler.CommandHandler) {
	handler.RegisterCommandHandler("아이템", ctx, client, h, onCreateItem)
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
