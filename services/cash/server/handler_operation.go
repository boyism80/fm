package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	gconstant "github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/item"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cash"
	"github.com/boyism80/fm/services/cash/client"
	"github.com/boyism80/fm/types"
)

type Operation struct {
	cs *CashShopServer
}

func (Operation) New(cs *CashShopServer) *Operation {
	return &Operation{cs: cs}
}

func (h *Operation) Handle(ctx *core.ClientContext, req *request.CashShopOperation) error {
	csClient, ok := ctx.Client.(*client.CashShopClient)
	if !ok {
		return fmt.Errorf("cash shop operation: invalid client type")
	}
	character := csClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("cash shop operation: character not loaded")
	}
	if character.Busy {
		return fmt.Errorf("cash shop operation: character %d is busy", character.ID())
	}

	switch req.Action {
	case request.CashShopActionBuy:
		h.buy(ctx, character, req)
	case request.CashShopActionPackage:
		h.buyPackage(ctx, character, req)
	case request.CashShopActionGift, request.CashShopActionGiftPackage:
		h.gift(ctx, character, req)
	case request.CashShopActionWishlist:
		h.setWishlist(ctx, character, req)
	case request.CashShopActionInventorySlots, request.CashShopActionStorageSlots, request.CashShopActionCharacterSlots:
		h.expandSlot(ctx, character, req)
	case request.CashShopActionTakeOut:
		h.takeOut(ctx, character, req)
	case request.CashShopActionPutIn:
		h.putIn(ctx, character, req)
	case request.CashShopActionPayBack:
		h.payBack(ctx, character, req)
	case request.CashShopActionQuestItem:
		h.buyQuestItem(ctx, character, req)
	case request.CashShopActionCoupleRing, request.CashShopActionFriendshipRing:
		h.buyRing(ctx, character, req)
	default:
		return fmt.Errorf("cash shop operation: unsupported action %d", req.Action)
	}
	return nil
}

func (h *Operation) start(ctx *core.ClientContext, character *cash.Character, failed constant.CashShopResultKind) *async.Task {
	character.Busy = true
	return async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).
		OnError(func(err error) {
			log.Printf("cash shop operation (async): character %d: %v", character.ID(), err)
			_ = ctx.Client.Send(&response.CashShopResult{Kind: failed, Failure: constant.CashShopFailureUnknown}, types.SEND_POLICY_ENCRYPT)
		}).
		Finally(func() {
			character.Busy = false
		})
}

func (h *Operation) findCommodity(sn uint32) (*wz.Commodity, constant.CashShopFailure) {
	commodity := h.cs.GetResources().Commodities[sn]
	if commodity == nil || !commodity.OnSale {
		return nil, constant.CashShopFailureNotPurchasableNow
	}
	return commodity, constant.CashShopFailureUnknown
}

func (h *Operation) buy(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBuyFailed, Failure: failure, CommoditySN: req.CommoditySN}, types.SEND_POLICY_ENCRYPT)
	}

	commodity, failure := h.findCommodity(req.CommoditySN)
	if commodity == nil {
		fail(failure)
		return
	}
	if commodity.Gender != 2 && commodity.Gender != character.Gender() {
		fail(constant.CashShopFailureGender)
		return
	}
	cashItem, err := character.NewCashItem(commodity, h.cs)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}

	h.purchase(ctx, character, req.Currency, commodity.Price, []*cash.LockerItem{cashItem}, constant.CashShopResultBuyFailed, func() {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBought, Item: character.BuildLockerItemS2C(cashItem)}, types.SEND_POLICY_ENCRYPT)
	})
}

func (h *Operation) buyPackage(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPackageFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	commodity, failure := h.findCommodity(req.CommoditySN)
	if commodity == nil {
		fail(failure)
		return
	}
	if commodity.Gender != 2 && commodity.Gender != character.Gender() {
		fail(constant.CashShopFailureGender)
		return
	}
	items, err := h.packageItems(character, commodity)
	if err != nil {
		log.Printf("cash shop package: character %d commodity %d: %v", character.ID(), commodity.SN, err)
		fail(constant.CashShopFailureUnknown)
		return
	}

	h.purchase(ctx, character, req.Currency, commodity.Price, items, constant.CashShopResultPackageFailed, func() {
		bought := make([]*dto.CashShopItem, 0, len(items))
		for _, item := range items {
			bought = append(bought, character.BuildLockerItemS2C(item))
		}
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPackageBought, Items: bought}, types.SEND_POLICY_ENCRYPT)
	})
}

func (h *Operation) packageItems(character *cash.Character, commodity *wz.Commodity) ([]*cash.LockerItem, error) {
	if len(commodity.Package) == 0 {
		return nil, fmt.Errorf("commodity %d is not a package", commodity.SN)
	}

	items := make([]*cash.LockerItem, 0, len(commodity.Package))
	for _, sn := range commodity.Package {
		content := h.cs.GetResources().Commodities[sn]
		if content == nil {
			return nil, fmt.Errorf("package content %d not found", sn)
		}
		item, err := character.NewCashItem(content, h.cs)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (h *Operation) purchase(ctx *core.ClientContext, character *cash.Character, currency uint8, price uint32, items []*cash.LockerItem, failed constant.CashShopResultKind, bought func()) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: failed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	pay, balance := character.Balance(currency)
	if balance < price {
		fail(constant.CashShopFailureNotEnoughCash)
		return
	}

	promise := h.start(ctx, character, failed)
	boughtItems := make([]*internal.CashItem, 0, len(items))
	for _, item := range items {
		boughtItems = append(boughtItems, item.BuildRPC())
	}
	promise.ThenRPC(func(c context.Context) (*internal.BuyCashItemReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.BuyCashItem(c, &internal.BuyCashItemRequest{
			WorldId:   uint32(h.cs.config.WorldID),
			AccountId: character.AccountID(),
			Currency:  pay,
			Price:     price,
			Items:     boughtItems,
		})
	}, func(reply *internal.BuyCashItemReply) error {
		character.SetBalance(reply.GetNxCash(), reply.GetMaplePoint())
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.AddLocker(items...)
			bought()
		case internal.CashShopResult_CASH_SHOP_RESULT_NOT_ENOUGH_CASH:
			fail(constant.CashShopFailureNotEnoughCash)
		case internal.CashShopResult_CASH_SHOP_RESULT_LOCKER_FULL:
			fail(constant.CashShopFailureLockerFull)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) buyQuestItem(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultQuestItemFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	commodity, failure := h.findCommodity(req.CommoditySN)
	if commodity == nil {
		fail(failure)
		return
	}
	wzItem := h.cs.GetResources().Items[commodity.ItemID]
	if wzItem == nil || wzItem.GetTraits()&gconstant.ItemTraitQuest == 0 {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if wzItem.GetTraits()&gconstant.ItemTraitUnique != 0 && character.HasItem(commodity.ItemID) {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if character.Meso() < int32(commodity.Price) {
		fail(constant.CashShopFailureNotEnoughMeso)
		return
	}
	invType := wzItem.GetInventoryType()
	slot, ok := character.FreeSlot(invType)
	if !ok {
		fail(constant.CashShopFailureLockerFull)
		return
	}
	pb, err := item.NewInventory(commodity.ItemID, max(commodity.Count, 1), character.ID(), int32(slot), h.cs.GetResources(), h.cs.NewUniqueID)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	pb.InventoryType = uint32(invType)
	if pb.UniqueId == nil {
		serial := h.cs.NewUniqueID()
		pb.UniqueId = &serial
	}
	shown, err := item.BuildS2C(pb, h.cs.GetResources(), h.cs.NewUniqueID)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}

	promise := h.start(ctx, character, constant.CashShopResultQuestItemFailed)
	promise.ThenRPC(func(c context.Context) (*internal.BuyCashQuestItemReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.BuyCashQuestItem(c, &internal.BuyCashQuestItemRequest{
			WorldId:     uint32(h.cs.config.WorldID),
			CharacterId: character.ID(),
			Price:       commodity.Price,
			Item:        pb,
		})
	}, func(reply *internal.BuyCashQuestItemReply) error {
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.SetMeso(reply.GetMeso())
			character.AddInventory(cash.ItemFromProto(pb))
			_ = ctx.Client.Send(&response.AddInventorySlot{InventoryType: invType, Slot: int16(slot), Item: shown}, types.SEND_POLICY_ENCRYPT)
			_ = ctx.Client.Send(&response.UpdateStats{Stats: map[gconstant.Stat]int32{gconstant.StatMeso: reply.GetMeso()}}, types.SEND_POLICY_ENCRYPT)
			_ = ctx.Client.Send(&response.CashShopResult{
				Kind:    constant.CashShopResultQuestItemBought,
				Granted: []*dto.CashShopGrantedItem{{Count: uint16(pb.GetCount()), Slot: uint16(slot), ItemID: commodity.ItemID}},
			}, types.SEND_POLICY_ENCRYPT)
		case internal.CashShopResult_CASH_SHOP_RESULT_NOT_ENOUGH_MESO:
			fail(constant.CashShopFailureNotEnoughMeso)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		return nil
	})
}

func (h *Operation) payBack(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPayBackFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	cashItem := character.FindLocker(req.Serial)
	if cashItem == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	model := h.cs.GetResources().Items[cashItem.ItemID()]
	if model == nil || model.GetInventoryType() != gconstant.InventoryTypeEquipment || cashItem.Expiration() > 0 {
		fail(constant.CashShopFailureUnknown)
		return
	}
	commodity := h.cs.GetResources().Commodities[cashItem.CommoditySN()]
	if commodity == nil || commodity.PayBackRate == 0 {
		fail(constant.CashShopFailureUnknown)
		return
	}
	refund := commodity.Price * commodity.PayBackRate / 100

	promise := h.start(ctx, character, constant.CashShopResultPayBackFailed)
	promise.ThenRPC(func(c context.Context) (*internal.PayBackCashItemReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.PayBackCashItem(c, &internal.PayBackCashItemRequest{
			WorldId:    uint32(h.cs.config.WorldID),
			AccountId:  character.AccountID(),
			Serial:     req.Serial,
			MaplePoint: refund,
		})
	}, func(reply *internal.PayBackCashItemReply) error {
		if reply.GetResult() != internal.CashShopResult_CASH_SHOP_RESULT_OK {
			fail(constant.CashShopFailureUnknown)
			return nil
		}
		character.RemoveLocker(req.Serial)
		character.SetBalance(reply.GetNxCash(), reply.GetMaplePoint())
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPaidBack, Serial: req.Serial, MaplePoint: refund}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) setWishlist(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	wishlist := make([]uint32, 0, len(req.Wishlist))
	for _, sn := range req.Wishlist {
		if sn != 0 {
			wishlist = append(wishlist, sn)
		}
	}

	promise := h.start(ctx, character, constant.CashShopResultWishlistFailed)
	promise.ThenRPC(func(c context.Context) (*internal.SetCashWishlistReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.SetCashWishlist(c, &internal.SetCashWishlistRequest{
			WorldId:      uint32(h.cs.config.WorldID),
			CharacterId:  character.ID(),
			CommoditySns: wishlist,
		})
	}, func(reply *internal.SetCashWishlistReply) error {
		character.Wishlist = wishlist
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultWishlistUpdated, Wishlist: wishlist}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) takeOut(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultTakeOutFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	cashItem := character.FindLocker(req.Serial)
	if cashItem == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	model := h.cs.GetResources().Items[cashItem.ItemID()]
	if model == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if model.GetTraits()&gconstant.ItemTraitUnique != 0 && character.HasItem(model.GetID()) {
		fail(constant.CashShopFailureUnknown)
		return
	}
	invType := model.GetInventoryType()
	slot, ok := character.FreeSlot(invType)
	if !ok {
		fail(constant.CashShopFailureNotEnoughSlots)
		return
	}
	record, ok := character.PrepareTakeOut(req.Serial, slot, invType, h.cs)
	if !ok {
		fail(constant.CashShopFailureUnknown)
		return
	}
	pb := record.BuildRPC()
	shown, err := item.BuildS2C(pb, h.cs.GetResources(), h.cs.NewUniqueID)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}

	promise := h.start(ctx, character, constant.CashShopResultTakeOutFailed)
	promise.ThenRPC(func(c context.Context) (*internal.TakeOutCashItemReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.TakeOutCashItem(c, &internal.TakeOutCashItemRequest{
			WorldId:     uint32(h.cs.config.WorldID),
			AccountId:   character.AccountID(),
			CharacterId: character.ID(),
			Item:        pb,
		})
	}, func(reply *internal.TakeOutCashItemReply) error {
		if reply.GetResult() != internal.CashShopResult_CASH_SHOP_RESULT_OK {
			fail(constant.CashShopFailureUnknown)
			return nil
		}
		character.TakeOut(req.Serial, record)
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultTakenOut, Slot: slot, TakenOut: shown}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) putIn(ctx *core.ClientContext, character *cash.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPutInFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	cashItem := character.PreparePutIn(req.Serial, gconstant.InventoryType(req.InventoryType))
	if cashItem == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}

	promise := h.start(ctx, character, constant.CashShopResultPutInFailed)
	promise.ThenRPC(func(c context.Context) (*internal.PutInCashItemReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.PutInCashItem(c, &internal.PutInCashItemRequest{
			WorldId:     uint32(h.cs.config.WorldID),
			AccountId:   character.AccountID(),
			CharacterId: character.ID(),
			Item:        cashItem.BuildRPC(),
		})
	}, func(reply *internal.PutInCashItemReply) error {
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.PutIn(req.Serial, gconstant.InventoryType(req.InventoryType), cashItem)
			_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPutIn, Item: character.BuildLockerItemS2C(cashItem)}, types.SEND_POLICY_ENCRYPT)
		case internal.CashShopResult_CASH_SHOP_RESULT_LOCKER_FULL:
			fail(constant.CashShopFailureLockerFull)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		return nil
	})
}
