package server

import (
	"fmt"
	"log"
	"strconv"

	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type CreateItem struct {
	gameServer *GameServer
}

func (*CreateItem) New(gameServer *GameServer) *CreateItem {
	return &CreateItem{
		gameServer: gameServer,
	}
}

func (h *CreateItem) GetCommandName() string {
	return "아이템생성"
}

func (h *CreateItem) GetUsage() string {
	return "<아이템ID/이름> [개수] - 아이템 생성"
}

func (h *CreateItem) Handle(gameClient *client.GameClient, args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing itemId or item name")
	}

	var itemId int
	var err error

	itemId, err = strconv.Atoi(args[0])
	if err != nil {
		itemIdUint, ok := h.gameServer.resources.NameToItem(args[0])
		if !ok {
			return fmt.Errorf("invalid itemId or item name: %s", args[0])
		}
		itemId = int(itemIdUint)
	}

	count := 1
	if len(args) >= 2 {
		if parsedCount, err := strconv.Atoi(args[1]); err == nil {
			count = parsedCount
		} else {
			log.Printf("Command: Invalid count, using default 1: %s", args[1])
		}
	}

	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("character not found")
	}

	_, ok := h.gameServer.resources.Items[uint32(itemId)]
	if !ok {
		return fmt.Errorf("item model not found for id: %d", itemId)
	}

	item, err := entity.NewItem(uint32(itemId), uint16(count), h.gameServer)
	if err != nil {
		return fmt.Errorf("failed to create item: %v", err)
	}

	inventoryType := item.GetInventoryType()
	inventory := character.Inventory[inventoryType]
	if inventory == nil {
		return fmt.Errorf("inventory not found for type: %v", inventoryType)
	}

	nextSlot, ok := inventory.NextSlot()
	if !ok {
		return fmt.Errorf("inventory is full for type: %v", inventoryType)
	}

	inventory.Items[int16(nextSlot)] = item

	itemDTO := entity.ItemToDTO(item)
	gameClient.Send(&response.AddItem{
		IsDrop:        false,
		Slot:          nextSlot,
		InventoryType: inventoryType,
		Item:          itemDTO,
	}, types.SEND_POLICY_ENCRYPT)

	log.Printf("Command: Created item %d, count %d in slot %d for character %d", itemId, count, nextSlot, character.GetID())
	return nil
}

