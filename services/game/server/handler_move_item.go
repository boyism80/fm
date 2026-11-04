package server

import (
	"errors"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
)

type MoveItem struct {
	gs *GameServer
}

func (MoveItem) New(gs *GameServer) *MoveItem {
	return &MoveItem{
		gs: gs,
	}
}

func (h *MoveItem) Handle(ctx *core.ClientContext, req *request.MoveItem) error {
	character, err := h.gs.character(ctx)
	if err != nil {
		return err
	}
	if character.Trading() || character.UsingCashItem() || character.Occupied() {
		character.Listener().OnUpdateStats(character, nil, true)
		return nil
	}

	switch {
	case req.Dest == 0:
		if !character.DropItem(ctx.ActorContext, req.InventoryType, req.Source, req.Count) {
			character.Listener().OnUpdateStats(character, nil, true)
		}
	case req.Source < 0:
		if err := character.Unequip(ctx.ActorContext, constant.EquipmentPartsType(req.Source), req.Dest); err != nil {
			character.Listener().OnUpdateStats(character, nil, true)
		}
	case req.Dest < 0:
		if err := character.Equip(ctx.ActorContext, req.Source, constant.EquipmentPartsType(req.Dest)); err != nil {
			if errors.Is(err, game.ErrInventoryFull) {
				character.Listener().OnItemGainFailed(character, constant.ItemGainFailedTypeFull)
			}
			character.Listener().OnUpdateStats(character, nil, true)
		}
	default:
		if !character.Inventory().Move(req.InventoryType, req.Source, req.Dest) {
			character.Listener().OnUpdateStats(character, nil, true)
		}
	}
	return nil
}
