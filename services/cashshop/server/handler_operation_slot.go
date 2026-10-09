package server

import (
	"context"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/entity"
	gconstant "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

const (
	slotExpansionPrice       = 3800
	slotExpansionAmount      = 4
	slotCouponAmount         = 8
	inventorySlotMax         = 96
	storageSlotMax           = 48
	characterSlotMax         = 15
	storageSlotCouponItem    = 9110000
	characterSlotCouponItem  = 5430000
	inventorySlotCouponFirst = 9111000
	inventorySlotCouponLast  = 9114000
)

func (h *Operation) expandSlot(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	var failed, expanded constant.CashShopResultKind
	expansion := &internal.ExpandCashSlotRequest{
		WorldId:     uint32(h.cs.config.WorldID),
		AccountId:   character.AccountID(),
		CharacterId: character.ID(),
		Price:       slotExpansionPrice,
		Amount:      slotExpansionAmount,
	}
	switch req.Action {
	case request.CashShopActionInventorySlots:
		failed, expanded = constant.CashShopResultInventorySlotsFailed, constant.CashShopResultInventorySlots
		expansion.Kind = internal.CashSlotKind_CASH_SLOT_KIND_INVENTORY
		expansion.InventoryType = uint32(req.InventoryType)
		expansion.Max = inventorySlotMax
	case request.CashShopActionStorageSlots:
		failed, expanded = constant.CashShopResultStorageSlotsFailed, constant.CashShopResultStorageSlots
		expansion.Kind = internal.CashSlotKind_CASH_SLOT_KIND_STORAGE
		expansion.Max = storageSlotMax
	case request.CashShopActionCharacterSlots:
		failed, expanded = constant.CashShopResultCharacterSlotsFailed, constant.CashShopResultCharacterSlots
		expansion.Kind = internal.CashSlotKind_CASH_SLOT_KIND_CHARACTER
		expansion.Max = characterSlotMax
	}
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: failed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	if req.ByCoupon || req.Action == request.CashShopActionCharacterSlots {
		commodity, failure := h.findCommodity(req.CommoditySN)
		if commodity == nil {
			fail(failure)
			return
		}
		expansion.Price = commodity.Price
		switch {
		case req.Action == request.CashShopActionCharacterSlots && commodity.ItemID == characterSlotCouponItem:
			expansion.Amount = 1
		case req.Action == request.CashShopActionStorageSlots && commodity.ItemID == storageSlotCouponItem:
			expansion.Amount = slotCouponAmount
		case req.Action == request.CashShopActionInventorySlots && commodity.ItemID >= inventorySlotCouponFirst && commodity.ItemID <= inventorySlotCouponLast:
			expansion.Amount = slotCouponAmount
			expansion.InventoryType = commodity.ItemID / 1000 % 10
		default:
			fail(constant.CashShopFailureUnknown)
			return
		}
	}
	if expansion.Kind == internal.CashSlotKind_CASH_SLOT_KIND_INVENTORY &&
		(expansion.InventoryType < uint32(gconstant.InventoryTypeEquipment) || expansion.InventoryType > uint32(gconstant.InventoryTypeETC)) {
		fail(constant.CashShopFailureUnknown)
		return
	}
	pay, balance := character.Balance(req.Currency)
	if balance < expansion.Price {
		fail(constant.CashShopFailureNotEnoughCash)
		return
	}
	expansion.Currency = pay

	promise := h.start(ctx, character, failed)
	promise.ThenRPC(func(c context.Context) (*internal.ExpandCashSlotReply, error) {
		return h.cs.internalClient.ExpandCashSlot(c, expansion)
	}, func(reply *internal.ExpandCashSlotReply) error {
		character.Busy = false
		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			slots := uint16(reply.GetSlots())
			switch expansion.Kind {
			case internal.CashSlotKind_CASH_SLOT_KIND_INVENTORY:
				character.SetSlotLimit(gconstant.InventoryType(expansion.InventoryType), uint8(slots))
			case internal.CashSlotKind_CASH_SLOT_KIND_STORAGE:
				character.StorageSlots = slots
			case internal.CashSlotKind_CASH_SLOT_KIND_CHARACTER:
				character.CharacterSlots = slots
			}
			_ = ctx.Client.Send(&response.CashShopResult{Kind: expanded, InventoryType: uint8(expansion.InventoryType), Slots: slots}, types.SEND_POLICY_ENCRYPT)
		case internal.CashShopResult_CASH_SHOP_RESULT_NOT_ENOUGH_CASH:
			fail(constant.CashShopFailureNotEnoughCash)
		case internal.CashShopResult_CASH_SHOP_RESULT_SLOT_LIMIT:
			fail(constant.CashShopFailureNotEnoughSlots)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}
