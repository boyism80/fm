package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/services/cashshop/entity"
	gconstant "github.com/boyism80/fm/services/game/constant"
	gentity "github.com/boyism80/fm/services/game/entity"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
	"google.golang.org/protobuf/proto"
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
		// TODO: ring system (equip effect, partner display on spawn, ring persistence) before selling couple and friendship rings
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBuyFailed, Failure: constant.CashShopFailureRing}, types.SEND_POLICY_ENCRYPT)
	default:
		return fmt.Errorf("cash shop operation: unsupported action %d", req.Action)
	}
	return nil
}

func (h *Operation) start(ctx *core.ClientContext, character *entity.Character, failed constant.CashShopResultKind) *async.Promise {
	character.Busy = true
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.OnError(func(err error) {
		character.Busy = false
		log.Printf("cash shop operation (async): character %d: %v", character.ID(), err)
		_ = ctx.Client.Send(&response.CashShopResult{Kind: failed, Failure: constant.CashShopFailureUnknown}, types.SEND_POLICY_ENCRYPT)
	})
	return promise
}

func (h *Operation) findCommodity(sn uint32) (*wz.Commodity, constant.CashShopFailure) {
	commodity := h.cs.resources.Commodities[sn]
	if commodity == nil || commodity.OnSale == false {
		return nil, constant.CashShopFailureNotPurchasableNow
	}
	return commodity, constant.CashShopFailureUnknown
}

func (h *Operation) buy(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
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

	h.purchase(ctx, character, req.Currency, commodity.Price, []*internal.CashItem{cashItem}, constant.CashShopResultBuyFailed, func() {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBought, Item: character.LockerItemDTO(cashItem)}, types.SEND_POLICY_ENCRYPT)
	})
}

func (h *Operation) buyPackage(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
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
			bought = append(bought, character.LockerItemDTO(item))
		}
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPackageBought, Items: bought}, types.SEND_POLICY_ENCRYPT)
	})
}

func (h *Operation) packageItems(character *entity.Character, commodity *wz.Commodity) ([]*internal.CashItem, error) {
	if len(commodity.Package) == 0 {
		return nil, fmt.Errorf("commodity %d is not a package", commodity.SN)
	}

	items := make([]*internal.CashItem, 0, len(commodity.Package))
	for _, sn := range commodity.Package {
		content := h.cs.resources.Commodities[sn]
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

func (h *Operation) purchase(ctx *core.ClientContext, character *entity.Character, currency uint8, price uint32, items []*internal.CashItem, failed constant.CashShopResultKind, bought func()) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: failed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	pay, balance := character.Balance(currency)
	if balance < price {
		fail(constant.CashShopFailureNotEnoughCash)
		return
	}

	promise := h.start(ctx, character, failed)
	async.ThenRPC(promise, func(c context.Context) (*internal.BuyCashItemReply, error) {
		return h.cs.internalClient.BuyCashItem(c, &internal.BuyCashItemRequest{
			WorldId:   h.cs.worldID(),
			AccountId: character.AccountID(),
			Currency:  pay,
			Price:     price,
			Items:     items,
		})
	}, func(reply *internal.BuyCashItemReply) error {
		character.Busy = false
		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.Locker = append(character.Locker, items...)
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

func (h *Operation) buyQuestItem(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultQuestItemFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	commodity, failure := h.findCommodity(req.CommoditySN)
	if commodity == nil {
		fail(failure)
		return
	}
	wzItem := h.cs.resources.Items[commodity.ItemID]
	if wzItem == nil || wzItem.IsQuest() == false {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if character.Meso() < int32(commodity.Price) {
		fail(constant.CashShopFailureNotEnoughMeso)
		return
	}
	invType := gconstant.GetInventoryTypeByItemID(commodity.ItemID)
	slot, ok := character.FreeSlot(invType)
	if ok == false {
		fail(constant.CashShopFailureLockerFull)
		return
	}
	item, err := gentity.NewItem(commodity.ItemID, max(commodity.Count, 1), h.cs)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	pb := item.ToProto(character.ID(), int32(slot))
	pb.InventoryType = uint32(invType)
	if pb.UniqueId == nil {
		serial := h.cs.NewUniqueID()
		pb.UniqueId = &serial
	}

	promise := h.start(ctx, character, constant.CashShopResultQuestItemFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.BuyCashQuestItemReply, error) {
		return h.cs.internalClient.BuyCashQuestItem(c, &internal.BuyCashQuestItemRequest{
			WorldId:     h.cs.worldID(),
			CharacterId: character.ID(),
			Price:       commodity.Price,
			Item:        pb,
		})
	}, func(reply *internal.BuyCashQuestItemReply) error {
		character.Busy = false
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.Game.Character.Meso = reply.GetMeso()
			character.Game.Inventory = append(character.Game.Inventory, pb)
			_ = ctx.Client.Send(&response.CashShopResult{
				Kind:    constant.CashShopResultQuestItemBought,
				Meso:    commodity.Price,
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

func (h *Operation) payBack(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPayBackFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	cashItem := character.FindLocker(req.Serial)
	if cashItem == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if gconstant.GetInventoryTypeByItemID(cashItem.GetItem().GetItemId()) != gconstant.InventoryTypeEquipment || cashItem.GetItem().GetExpirationUnixMs() > 0 {
		fail(constant.CashShopFailureUnknown)
		return
	}
	commodity := h.cs.resources.Commodities[cashItem.GetCommoditySn()]
	if commodity == nil || commodity.PayBackRate == 0 {
		fail(constant.CashShopFailureUnknown)
		return
	}
	refund := commodity.Price * commodity.PayBackRate / 100

	promise := h.start(ctx, character, constant.CashShopResultPayBackFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.PayBackCashItemReply, error) {
		return h.cs.internalClient.PayBackCashItem(c, &internal.PayBackCashItemRequest{
			WorldId:    h.cs.worldID(),
			AccountId:  character.AccountID(),
			Serial:     req.Serial,
			MaplePoint: refund,
		})
	}, func(reply *internal.PayBackCashItemReply) error {
		character.Busy = false
		if reply.GetResult() != internal.CashShopResult_CASH_SHOP_RESULT_OK {
			fail(constant.CashShopFailureUnknown)
			return nil
		}
		character.RemoveLocker(req.Serial)
		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPaidBack, Serial: req.Serial, MaplePoint: refund}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) setWishlist(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	wishlist := make([]uint32, 0, len(req.Wishlist))
	for _, sn := range req.Wishlist {
		if sn != 0 {
			wishlist = append(wishlist, sn)
		}
	}

	promise := h.start(ctx, character, constant.CashShopResultWishlistFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.SetCashWishlistReply, error) {
		return h.cs.internalClient.SetCashWishlist(c, &internal.SetCashWishlistRequest{
			WorldId:      h.cs.worldID(),
			CharacterId:  character.ID(),
			CommoditySns: wishlist,
		})
	}, func(reply *internal.SetCashWishlistReply) error {
		character.Busy = false
		character.Wishlist = wishlist
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultWishlistUpdated, Wishlist: wishlist}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) takeOut(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultTakeOutFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	cashItem := character.FindLocker(req.Serial)
	if cashItem == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	invType := gconstant.GetInventoryTypeByItemID(cashItem.GetItem().GetItemId())
	slot, ok := character.FreeSlot(invType)
	if ok == false {
		fail(constant.CashShopFailureNotEnoughSlots)
		return
	}
	item, err := gentity.NewItemFromInternalProto(cashItem.GetItem(), h.cs)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	pb := item.ToProto(character.ID(), int32(slot))
	pb.InventoryType = uint32(invType)

	promise := h.start(ctx, character, constant.CashShopResultTakeOutFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.TakeOutCashItemReply, error) {
		return h.cs.internalClient.TakeOutCashItem(c, &internal.TakeOutCashItemRequest{
			WorldId:     h.cs.worldID(),
			AccountId:   character.AccountID(),
			CharacterId: character.ID(),
			Item:        pb,
		})
	}, func(reply *internal.TakeOutCashItemReply) error {
		character.Busy = false
		if reply.GetResult() != internal.CashShopResult_CASH_SHOP_RESULT_OK {
			fail(constant.CashShopFailureUnknown)
			return nil
		}
		character.RemoveLocker(req.Serial)
		character.Game.Inventory = append(character.Game.Inventory, pb)
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultTakenOut, Slot: slot, TakenOut: item.ToDTO()}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}

func (h *Operation) putIn(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPutInFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	index, pb := character.FindInventory(req.Serial, gconstant.InventoryType(req.InventoryType))
	if pb == nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	stored := proto.Clone(pb).(*internal.Inventory)
	stored.Slot = 0
	cashItem := &internal.CashItem{BuyerName: character.Name(), Item: stored}

	promise := h.start(ctx, character, constant.CashShopResultPutInFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.PutInCashItemReply, error) {
		return h.cs.internalClient.PutInCashItem(c, &internal.PutInCashItemRequest{
			WorldId:     h.cs.worldID(),
			AccountId:   character.AccountID(),
			CharacterId: character.ID(),
			Item:        cashItem,
		})
	}, func(reply *internal.PutInCashItemReply) error {
		character.Busy = false
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.Game.Inventory = append(character.Game.Inventory[:index], character.Game.Inventory[index+1:]...)
			character.Locker = append(character.Locker, cashItem)
			_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultPutIn, Item: character.LockerItemDTO(cashItem)}, types.SEND_POLICY_ENCRYPT)
		case internal.CashShopResult_CASH_SHOP_RESULT_LOCKER_FULL:
			fail(constant.CashShopFailureLockerFull)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		return nil
	})
}
