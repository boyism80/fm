package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/services/cashshop/entity"
	"github.com/boyism80/fm/types"
)

type Login struct {
	cs *CashShopServer
}

func (Login) New(cs *CashShopServer) *Login {
	return &Login{cs: cs}
}

func (h *Login) Handle(ctx *core.ClientContext, req *request.LoginGame) error {
	csClient, ok := ctx.Client.(*client.CashShopClient)
	if !ok {
		return fmt.Errorf("cash shop login: invalid client type")
	}

	var entered *internal.EnterCashShopReply
	promise := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise = async.ThenRPC(promise, func(c context.Context) (*internal.EnterCashShopReply, error) {
		return h.cs.internalClient.EnterCashShop(c, &internal.EnterCashShopRequest{
			WorldId:     h.cs.worldID(),
			CharacterId: req.PlayerId,
			CashShopId:  h.cs.cashShopID(),
			ClientIp:    ctx.Client.GetRemoteIP(),
		})
	}, func(reply *internal.EnterCashShopReply) error {
		if reply.GetGame().GetFound() == false || reply.GetGame().GetCharacter() == nil {
			return fmt.Errorf("character %d not found", req.PlayerId)
		}
		entered = reply

		character := entity.NewCharacter(reply)
		if csClient.SetCharacter(character) == false {
			return fmt.Errorf("character %d: client disconnected during cash shop login", req.PlayerId)
		}

		_ = ctx.Client.Send(&response.SetCashShop{
			CharacterInfo: response.CharacterInfo{Character: character.ToDTO(h.cs)},
			AccountName:   character.Name(),
		}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopResult{
			Kind:           constant.CashShopResultLocker,
			Locker:         character.LockerDTO(),
			StorageSlots:   character.StorageSlots,
			CharacterSlots: character.CharacterSlots,
		}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultGifts}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultWishlist, Wishlist: character.Wishlist}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("cash shop login (async): %v", err)
		_ = ctx.Client.GetConnection().Close()
		if entered == nil || csClient.GetCharacter() != nil {
			return
		}

		characterID := req.PlayerId
		cashShopID := h.cs.cashShopID()
		logout := async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		logout = async.ThenRPC(logout, func(c context.Context) (*internal.LogoutSessionReply, error) {
			return h.cs.internalClient.LogoutSession(c, &internal.LogoutSessionRequest{
				WorldId:          h.cs.worldID(),
				AccountId:        entered.GetGame().GetCharacter().GetAccountId(),
				CharacterId:      &characterID,
				CashShopId:       &cashShopID,
				DisconnectSource: internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_CASH_SHOP_SERVER,
			})
		}, func(*internal.LogoutSessionReply) error {
			return nil
		})
		logout.OnError(func(err error) {
			log.Printf("cash shop login: logout session of character %d after login failure: %v", characterID, err)
		})
	})
	return nil
}
