package server

import (
	"context"
	"unicode/utf8"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/entity"
	gconstant "github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/types"
)

const ringMessageMaxLength = 73

func (h *Operation) buyRing(ctx *core.ClientContext, character *entity.Character, req *request.CashShopOperation) {
	fail := func(failure constant.CashShopFailure) {
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBuyFailed, Failure: failure, CommoditySN: req.CommoditySN}, types.SEND_POLICY_ENCRYPT)
	}

	commodity, failure := h.findCommodity(req.CommoditySN)
	if commodity == nil {
		fail(failure)
		return
	}
	couple := req.Action == request.CashShopActionCoupleRing
	ring := gconstant.IsFriendshipRing(commodity.ItemID)
	if couple {
		ring = gconstant.IsCrushRing(commodity.ItemID)
	}
	if !ring {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if length := utf8.RuneCountInString(req.Message); length == 0 || length > ringMessageMaxLength {
		fail(constant.CashShopFailureUnknown)
		return
	}
	if commodity.Gender != 2 && commodity.Gender != character.Gender() {
		fail(constant.CashShopFailureRing)
		return
	}
	if character.NXCash < commodity.Price {
		fail(constant.CashShopFailureNotEnoughCash)
		return
	}
	item, err := character.NewCashItem(commodity, h.cs)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}
	partnerItem, err := character.NewCashItem(commodity, h.cs)
	if err != nil {
		fail(constant.CashShopFailureUnknown)
		return
	}

	promise := h.start(ctx, character, constant.CashShopResultBuyFailed)
	promise.ThenRPC(func(c context.Context) (*internal.BuyCashRingReply, error) {
		return h.cs.internalClient.BuyCashRing(c, &internal.BuyCashRingRequest{
			WorldId:       h.cs.worldID(),
			AccountId:     character.AccountID(),
			CharacterId:   character.ID(),
			CharacterName: character.Name(),
			Gender:        uint32(character.Gender()),
			Couple:        couple,
			PartnerName:   req.Recipient,
			Message:       req.Message,
			Price:         commodity.Price,
			Item:          item,
			PartnerItem:   partnerItem,
		})
	}, func(reply *internal.BuyCashRingReply) error {
		character.Busy = false
		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
			character.Locker = append(character.Locker, item)
			_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultBought, Item: character.LockerItemDTO(item)}, types.SEND_POLICY_ENCRYPT)
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
			fail(constant.CashShopFailureRing)
		case internal.CashShopResult_CASH_SHOP_RESULT_LOCKER_FULL:
			fail(constant.CashShopFailureLockerFull)
		case internal.CashShopResult_CASH_SHOP_RESULT_RECIPIENT_LOCKER_FULL:
			fail(constant.CashShopFailureRecipientFull)
		default:
			fail(constant.CashShopFailureUnknown)
		}
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}
