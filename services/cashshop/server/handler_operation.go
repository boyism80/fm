package server

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/services/cashshop/entity"
	gconstant "github.com/boyism80/fm/services/game/constant"
	gentity "github.com/boyism80/fm/services/game/entity"
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
	case request.CashShopActionWishlist:
		h.setWishlist(ctx, character, req)
	case request.CashShopActionTakeOut:
		h.takeOut(ctx, character, req)
	case request.CashShopActionPutIn:
		h.putIn(ctx, character, req)
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

func (h *Operation) buy(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBuyFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
	}

	commodity := h.cs.resources.Commodities[req.CommoditySN]
	if commodity == nil || commodity.OnSale == false {
		fail(constant.CashShopFailureNotPurchasableNow)
		return
	}
	if commodity.Gender != 2 && commodity.Gender != character.Gender() {
		fail(constant.CashShopFailureGender)
		return
	}

	currency := internal.CashCurrency_CASH_CURRENCY_NX_CASH
	balance := character.NXCash
	if req.Currency == 1 {
		currency = internal.CashCurrency_CASH_CURRENCY_MAPLE_POINT
		balance = character.MaplePoint
	}
	if balance < commodity.Price {
		fail(constant.CashShopFailureNotEnoughCash)
		return
	}

	item, err := gentity.NewItem(commodity.ItemID, max(commodity.Count, 1), h.cs)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	pb := item.ToProto(character.ID(), 0)
	if pb.UniqueId == nil {
		serial := h.cs.NewUniqueID()
		pb.UniqueId = &serial
	}
	if commodity.Period > 0 {
		pb.ExpirationUnixMs = time.Now().AddDate(0, 0, int(commodity.Period)).UnixMilli()
	}
	cashItem := &internal.CashItem{CommoditySn: commodity.SN, BuyerName: character.Name(), Item: pb}

	promise := h.start(ctx, character, constant.CashShopResultBuyFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.BuyCashItemReply, error) {
		return h.cs.internalClient.BuyCashItem(c, &internal.BuyCashItemRequest{
			WorldId:   h.cs.worldID(),
			AccountId: character.AccountID(),
			Currency:  currency,
			Price:     commodity.Price,
			Item:      cashItem,
		})
	}, func(reply *internal.BuyCashItemReply) error {
		character.Busy = false
		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.Locker = append(character.Locker, cashItem)
			_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBought, Item: character.LockerItemDTO(cashItem)}, types.SEND_POLICY_ENCRYPT)
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

func (h *Operation) setWishlist(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	wishlist := make([]uint32, 0, len(req.Wishlist))
	for _, sn := range req.Wishlist {
		if sn != 0 {
			wishlist = append(wishlist, sn)
		}
	}

	promise := h.start(ctx, character, constant.CashShopResultBuyFailed)
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

	index, cashItem := character.FindLocker(req.Serial)
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
		character.Locker = append(character.Locker[:index], character.Locker[index+1:]...)
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
