package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	"github.com/boyism80/fm/core/rpc"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cash"
	"github.com/boyism80/fm/services/cash/client"
	"github.com/boyism80/fm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.EnterCashShopReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.EnterCashShop(c, &internal.EnterCashShopRequest{
			WorldId:     uint32(h.cs.config.WorldID),
			CharacterId: req.PlayerId,
			CashShopId:  uint32(h.cs.config.CashShopID),
			ClientIp:    ctx.Client.GetRemoteIP(),
		})
	}, func(reply *internal.EnterCashShopReply) error {
		if !reply.GetGame().GetFound() || reply.GetGame().GetCharacter() == nil {
			return fmt.Errorf("character %d not found", req.PlayerId)
		}
		entered = reply

		character := cash.NewCharacter(reply, h.cs.GetResources())
		if !csClient.SetCharacter(character) {
			return fmt.Errorf("character %d: client disconnected during cash shop login", req.PlayerId)
		}

		_ = ctx.Client.Send(&response.SetCashShop{
			CharacterInfo: response.CharacterInfo{Character: character.BuildS2C(h.cs)},
			AccountName:   character.Name(),
		}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopResult{
			Kind:           constant.CashShopResultLocker,
			Locker:         character.BuildLockerS2C(),
			StorageSlots:   character.StorageSlots,
			CharacterSlots: character.CharacterSlots,
		}, types.SEND_POLICY_ENCRYPT)
		for _, serial := range character.Expired {
			_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultExpired, Serial: serial}, types.SEND_POLICY_ENCRYPT)
			character.RemoveLocker(serial)
		}
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultGifts, Gifts: character.BuildGiftsS2C()}, types.SEND_POLICY_ENCRYPT)
		_ = ctx.Client.Send(&response.CashShopResult{Kind: constant.CashShopResultWishlist, Wishlist: character.Wishlist}, types.SEND_POLICY_ENCRYPT)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("cash shop login (async): %v", err)
		_ = ctx.Client.GetConnection().Close()
		if csClient.GetCharacter() != nil {
			return
		}

		// Only an EnterCashShop that internal answered with an error is known not to have entered; internal rolled that one back.
		// A promise future timeout is not a gRPC status, and its RPC may still enter the cash shop.
		if entered == nil {
			s, ok := status.FromError(err)
			if ok && s.Code() != codes.DeadlineExceeded && s.Code() != codes.Unavailable && s.Code() != codes.Canceled {
				return
			}
		}

		characterID := req.PlayerId
		cashShopID := uint32(h.cs.config.CashShopID)
		logout := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
		logout.ThenRPC(func(c context.Context) (*internal.LogoutSessionReply, error) {
			return h.cs.systems.Get[*rpc.InternalConnection]().Client.LogoutSession(c, &internal.LogoutSessionRequest{
				WorldId:          uint32(h.cs.config.WorldID),
				AccountId:        entered.GetGame().GetCharacter().GetAccountId(),
				CharacterId:      &characterID,
				CashShopId:       &cashShopID,
				DisconnectSource: internal.SessionDisconnectSource_SESSION_DISCONNECT_SOURCE_CASH_SHOP_SERVER,
			})
		}, nil)
		logout.OnError(func(err error) {
			log.Printf("cash shop login: logout session of character %d after login failure: %v", characterID, err)
		})
	})
	return nil
}
