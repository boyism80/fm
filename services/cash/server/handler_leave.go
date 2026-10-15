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
	"github.com/boyism80/fm/services/cash/client"
	"github.com/boyism80/fm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Leave struct {
	cs *CashShopServer
}

func (Leave) New(cs *CashShopServer) *Leave {
	return &Leave{cs: cs}
}

func (h *Leave) Handle(ctx *core.ClientContext, req *request.LeaveCashShop) error {
	csClient, ok := ctx.Client.(*client.CashShopClient)
	if !ok {
		return fmt.Errorf("leave cash shop: invalid client type")
	}
	character := csClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("leave cash shop: character not loaded")
	}
	if character.Busy {
		return fmt.Errorf("leave cash shop: character %d is busy", character.ID())
	}

	character.Busy = true
	blocked := func() {
		character.Busy = false
		_ = ctx.Client.Send(&response.ServerBlocked{Reason: constant.ServerBlockedChannelMoveUnavailable}, types.SEND_POLICY_ENCRYPT)
	}
	var route *response.SwitchChannel
	transitionPending := false
	cashShopID := uint32(h.cs.config.CashShopID)
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.GetGameChannelStatusReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.GetGameChannelStatus(c, &internal.GetGameChannelStatusRequest{
			WorldId:   uint32(h.cs.config.WorldID),
			ChannelId: character.ReturnChannel,
		})
	}, func(reply *internal.GetGameChannelStatusReply) error {
		if !reply.GetFound() || !reply.GetAlive() || reply.GetChannelFull() {
			blocked()
			return fmt.Errorf("return channel %d unavailable", character.ReturnChannel)
		}
		route = &response.SwitchChannel{IP: reply.GetHost(), Port: uint16(reply.GetPort())}
		transitionPending = true
		return nil
	})
	promise.ThenRPC(func(c context.Context) (*internal.BeginGameTransitionReply, error) {
		return h.cs.systems.Get[*rpc.InternalConnection]().Client.BeginGameTransition(c, &internal.BeginGameTransitionRequest{
			WorldId:          uint32(h.cs.config.WorldID),
			AccountId:        character.AccountID(),
			CharacterId:      character.ID(),
			ClientIp:         ctx.Client.GetRemoteIP(),
			Debuffs:          character.Game.GetDebuffs(),
			SourceCashShopId: &cashShopID,
		})
	}, func(reply *internal.BeginGameTransitionReply) error {
		transitionPending = false
		if !reply.GetOk() {
			blocked()
			return fmt.Errorf("begin transition failed: %s", reply.GetErrorCode())
		}
		csClient.Leave()
		_ = ctx.Client.Send(route, types.SEND_POLICY_ENCRYPT)
		return nil
	})
	promise.OnError(func(err error) {
		log.Printf("leave cash shop (async): character %d: %v", character.ID(), err)

		// A transition internal never answered may have committed, so the character stays busy and the connection closes;
		// the disconnect logout then ends the session only if it is still this cash shop's.
		// A promise future timeout is not a gRPC status.
		if transitionPending {
			s, ok := status.FromError(err)
			if !ok || s.Code() == codes.DeadlineExceeded || s.Code() == codes.Unavailable || s.Code() == codes.Canceled {
				_ = ctx.Client.GetConnection().Close()
				return
			}
		}
		character.Busy = false
	})
	return nil
}
