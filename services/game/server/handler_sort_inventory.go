package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
)

type SortInventory struct {
	gs *GameServer
}

func (SortInventory) New(gs *GameServer) *SortInventory {
	return &SortInventory{
		gs: gs,
	}
}

func (h *SortInventory) Handle(ctx *core.ClientContext, req *request.SortInventory) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}
	if character.Trading() || character.UsingCashItem() || character.Occupied() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}

	character.Inventory().Merge(req.InventoryType)
	character.Inventory().Sort(req.InventoryType)

	character.Listener().OnEndSortInventory(character, req.InventoryType)
	character.Listener().OnUpdateStats(character, nil, true)

	return nil
}
