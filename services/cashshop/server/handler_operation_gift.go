package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/entity"
	"github.com/boyism80/fm/types"
)

func (h *Operation) gift(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultGiftFailed, Failure: failure, CommoditySN: req.CommoditySN}, types.SEND_POLICY_ENCRYPT)
	}

	commodity, failure := h.findCommodity(req.CommoditySN)
	if commodity == nil {
		fail(failure)
		return
	}
	if req.Recipient == "" {
		fail(constant.CashShopFailureWrongName)
		return
	}
	if character.NXCash < commodity.Price {
		fail(constant.CashShopFailureNotEnoughCash)
		return
	}

	var items []*internal.CashItem
	if req.Action == request.CashShopActionGiftPackage {
		packaged, err := h.packageItems(character, commodity)
		if err != nil {
			log.Printf("cash shop gift package: character %d commodity %d: %v", character.ID(), commodity.SN, err)
			fail(constant.CashShopFailureUnknown)
			return
		}
		items = packaged
	} else {
		item, err := character.NewCashItem(commodity, h.cs)
		if err != nil {
			fail(constant.CashShopFailureUnknown)
			return
		}
		items = []*internal.CashItem{item}
	}

	promise := h.start(ctx, character, constant.CashShopResultGiftFailed)
	async.ThenRPC(promise, func(c context.Context) (*internal.GiftCashItemReply, error) {
		return h.cs.internalClient.GiftCashItem(c, &internal.GiftCashItemRequest{
			WorldId:       h.cs.worldID(),
			AccountId:     character.AccountID(),
			RecipientName: req.Recipient,
			Message:       req.Message,
			Price:         commodity.Price,
			Gender:        uint32(commodity.Gender),
			Items:         items,
		})
	}, func(reply *internal.GiftCashItemReply) error {
		character.Busy = false
		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			_ = ctx.Client.Send(&response.CashShopResult{
				Kind:      constant.CashShopResultGiftSent,
				Recipient: req.Recipient,
				ItemID:    commodity.ItemID,
				Count:     max(commodity.Count, 1),
			}, types.SEND_POLICY_ENCRYPT)
		case internal.CashShopResult_CASH_SHOP_RESULT_NOT_ENOUGH_CASH:
			fail(constant.CashShopFailureNotEnoughCash)
		case internal.CashShopResult_CASH_SHOP_RESULT_NOT_FOUND:
			fail(constant.CashShopFailureWrongName)
		case internal.CashShopResult_CASH_SHOP_RESULT_SAME_ACCOUNT:
			fail(constant.CashShopFailureSameAccount)
		case internal.CashShopResult_CASH_SHOP_RESULT_GENDER:
			fail(constant.CashShopFailureGender)
		case internal.CashShopResult_CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL:
			fail(constant.CashShopFailureRecipientFull)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}
