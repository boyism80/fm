package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	gconstant "github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cash"
	"github.com/boyism80/fm/services/cash/client"
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

func (h *Coupon) start(ctx *core.ClientContext, character *cash.Character) *async.Task {
	character.Busy = true
	return async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).
		OnError(func(err error) {
			log.Printf("cash shop coupon (async): character %d: %v", character.ID(), err)
			h.fail(ctx, constant.CashShopFailureUnknown)
		}).
		Finally(func() {
			character.Busy = false
		})
}

func (h *Coupon) redeem(ctx *core.ClientContext, character *cash.Character, req *request.CashShopCoupon) {
	if req.Recipient != "" {
		h.fail(ctx, constant.CashShopFailureCouponNoGift)
		return
	}
	if req.Code == "" {
		h.fail(ctx, constant.CashShopFailureCouponWrong)
		return
	}

	promise := h.start(ctx, character)
	promise.ThenRPC(func(c context.Context) (*internal.FindCashCouponReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.FindCashCoupon(c, &internal.FindCashCouponRequest{WorldId: uint32(h.cs.config.WorldID), Code: req.Code, AccountId: character.AccountID()})
	}, func(reply *internal.FindCashCouponReply) error {
		switch reply.GetResult() {
		case internal.CashShopResult_CASH_SHOP_RESULT_OK:
		case internal.CashShopResult_CASH_SHOP_RESULT_COUPON_USED:
			h.fail(ctx, constant.CashShopFailureCouponUsed)
			return nil
		default:
			h.fail(ctx, constant.CashShopFailureCouponWrong)
			return nil
		}

		var item *cash.LockerItem
		if reply.GetKind() == internal.CashCouponKind_CASH_COUPON_KIND_ITEM {
			commodity := h.cs.GetResources().Commodities[reply.GetValue()]
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
		h.claim(ctx, promise, character, req.Code, reply.GetKind(), reply.GetValue(), item)
		return nil
	})
}

// claim continues the task redeem started, so Busy stays set until the coupon is spent.
func (h *Coupon) claim(ctx *core.ClientContext, promise *async.Task, character *cash.Character, code string, kind internal.CashCouponKind, value uint32, item *cash.LockerItem) {
	promise.ThenRPC(func(c context.Context) (*internal.RedeemCashCouponReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.RedeemCashCoupon(c, &internal.RedeemCashCouponRequest{
			WorldId:     uint32(h.cs.config.WorldID),
			AccountId:   character.AccountID(),
			CharacterId: character.ID(),
			Code:        code,
			Item:        item.BuildRPC(),
			RequestId:   h.cs.NewUniqueID(),
		})
	}, func(reply *internal.RedeemCashCouponReply) error {
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

		character.SetBalance(reply.GetNxCash(), reply.GetMaplePoint())
		character.SetMeso(reply.GetMeso())
		redeemed := &response.CashShopResult{Kind: constant.CashShopResultCouponRedeemed}
		switch kind {
		case internal.CashCouponKind_CASH_COUPON_KIND_MAPLE_POINT:
			redeemed.MaplePoint = value
		case internal.CashCouponKind_CASH_COUPON_KIND_MESO:
			redeemed.Meso = value
		case internal.CashCouponKind_CASH_COUPON_KIND_ITEM:
			character.AddLocker(item)
			redeemed.Items = []*dto.CashShopItem{character.BuildLockerItemS2C(item)}
		}
		_ = ctx.Client.Send(redeemed, types.SEND_POLICY_ENCRYPT)
		if kind == internal.CashCouponKind_CASH_COUPON_KIND_MESO {
			_ = ctx.Client.Send(&response.UpdateStats{Stats: map[gconstant.Stat]int32{gconstant.StatMeso: reply.GetMeso()}}, types.SEND_POLICY_ENCRYPT)
		}
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
}
