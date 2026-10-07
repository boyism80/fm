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
	"github.com/boyism80/fm/types"
)

type Coupon struct {
	cs *CashShopServer
}

func (Coupon) New(cs *CashShopServer) *Coupon {
	return &Coupon{cs: cs}
}

func (h *Coupon) Handle(ctx *core.ClientContext, req *request.CashShopCoupon) error {
	csClient, ok := ctx.Client.(*client.CashShopClient)
	if !ok {
		return fmt.Errorf("cash shop coupon: invalid client type")
	}
	character := csClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("cash shop coupon: character not loaded")
	}
	if character.Busy {
		return fmt.Errorf("cash shop coupon: character %d is busy", character.ID())
	}

	h.redeem(ctx, character, req)
	return nil
}

func (h *Coupon) fail(ctx *core.ClientContext, failure constant.CashShopFailure) {
	_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultCouponFailed, Failure: failure}, types.SEND_POLICY_ENCRYPT)
}

func (h *Coupon) start(ctx *core.ClientContext, character *entity.Character) *async.Promise {
	character.Busy = true
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.OnError(func(err error) {
		character.Busy = false
		log.Printf("cash shop coupon (async): character %d: %v", character.ID(), err)
		h.fail(ctx, constant.CashShopFailureUnknown)
	})
	return promise
}

func (h *Coupon) redeem(ctx *core.ClientContext, character *entity.Character, req *request.CashShopCoupon) {
	if req.Recipient != "" {
		h.fail(ctx, constant.CashShopFailureCouponNoGift)
		return
	}
	if req.Code == "" {
		h.fail(ctx, constant.CashShopFailureCouponWrong)
		return
	}

	promise := h.start(ctx, character)
	async.ThenRPC(promise, func(c context.Context) (*internal.FindCashCouponReply, error) {
		return h.cs.internalClient.FindCashCoupon(c, &internal.FindCashCouponRequest{WorldId: h.cs.worldID(), Code: req.Code})
	}, func(reply *internal.FindCashCouponReply) error {
		character.Busy = false
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
		case internal.CashShopResult_CASH_SHOP_RESULT_COUPON_USED:
			h.fail(ctx, constant.CashShopFailureCouponUsed)
			return nil
		default:
			h.fail(ctx, constant.CashShopFailureCouponWrong)
			return nil
		}

		var item *internal.CashItem
		if reply.GetKind() == internal.CashCouponKind_CASH_COUPON_KIND_ITEM {
			commodity := h.cs.resources.Commodities[reply.GetValue()]
			if commodity == nil {
				h.fail(ctx, constant.CashShopFailureUnknown)
				return nil
			}
			created, err := character.NewCashItem(commodity, h.cs)
			if err != nil {
				h.fail(ctx, constant.CashShopFailureUnknown)
				return nil
			}
			item = created
		}
		h.claim(ctx, character, req.Code, reply.GetKind(), reply.GetValue(), item)
		return nil
	})
}

func (h *Coupon) claim(ctx *core.ClientContext, character *entity.Character, code string, kind internal.CashCouponKind, value uint32, item *internal.CashItem) {
	promise := h.start(ctx, character)
	async.ThenRPC(promise, func(c context.Context) (*internal.RedeemCashCouponReply, error) {
		return h.cs.internalClient.RedeemCashCoupon(c, &internal.RedeemCashCouponRequest{
			WorldId:     h.cs.worldID(),
			AccountId:   character.AccountID(),
			CharacterId: character.ID(),
			Code:        code,
			Item:        item,
		})
	}, func(reply *internal.RedeemCashCouponReply) error {
		character.Busy = false
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
		case internal.CashShopResult_CASH_SHOP_RESULT_COUPON_USED:
			h.fail(ctx, constant.CashShopFailureCouponUsed)
			return nil
		case internal.CashShopResult_CASH_SHOP_RESULT_LOCKER_FULL:
			h.fail(ctx, constant.CashShopFailureLockerFull)
			return nil
		default:
			h.fail(ctx, constant.CashShopFailureUnknown)
			return nil
		}

		character.NXCash = reply.GetNxCash()
		character.MaplePoint = reply.GetMaplePoint()
		character.Game.Character.Meso = reply.GetMeso()
		redeemed := &response.CashShopResult{Kind: constant.CashShopResultCouponRedeemed}
		switch kind {
		case internal.CashCouponKind_CASH_COUPON_KIND_MAPLE_POINT:
			redeemed.MaplePoint = value
		case internal.CashCouponKind_CASH_COUPON_KIND_MESO:
			redeemed.Meso = value
		case internal.CashCouponKind_CASH_COUPON_KIND_ITEM:
			character.Locker = append(character.Locker, item)
			redeemed.Items = []*dto.CashShopItem{character.LockerItemDTO(item)}
		}
		_ = ctx.Client.Send(redeemed, types.SEND_POLICY_ENCRYPT)
		if kind == internal.CashCouponKind_CASH_COUPON_KIND_MESO {
			_ = ctx.Client.Send(&response.UpdateStats{Stats: map[gconstant.Stat]int32{gconstant.StatMeso: reply.GetMeso()}}, types.SEND_POLICY_ENCRYPT)
		}
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}
