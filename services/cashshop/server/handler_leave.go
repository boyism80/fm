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
	"github.com/boyism80/fm/types"
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
	cashShopID := h.cs.cashShopID()
	promise := async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout)
	promise.ThenRPC(func(c context.Context) (*internal.GetGameChannelStatusReply, error) {
		return h.cs.internalClient.GetGameChannelStatus(c, &internal.GetGameChannelStatusRequest{
			WorldId:   h.cs.worldID(),
			ChannelId: character.ReturnChannel,
		})
	}, func(reply *internal.GetGameChannelStatusReply) error {
		if reply.GetFound() == false || reply.GetAlive() == false || reply.GetChannelFull() {
			blocked()
			return fmt.Errorf("return channel %d unavailable", character.ReturnChannel)
		}
		route = &response.SwitchChannel{IP: reply.GetHost(), Port: uint16(reply.GetPort())}
		return nil
	})
	promise.ThenRPC(func(c context.Context) (*internal.BeginGameTransitionReply, error) {
		return h.cs.internalClient.BeginGameTransition(c, &internal.BeginGameTransitionRequest{
			WorldId:          h.cs.worldID(),
			AccountId:        character.AccountID(),
			CharacterId:      character.ID(),
			ClientIp:         ctx.Client.GetRemoteIP(),
			Debuffs:          character.Game.GetDebuffs(),
			SourceCashShopId: &cashShopID,
		})
	}, func(reply *internal.BeginGameTransitionReply) error {
		if reply.GetOk() == false {
			blocked()
			return fmt.Errorf("begin transition failed: %s", reply.GetErrorCode())
		}
		csClient.Leave()
		_ = ctx.Client.Send(route, types.SEND_POLICY_ENCRYPT)
		return nil
	})
	promise.OnError(func(err error) {
		character.Busy = false
		log.Printf("leave cash shop (async): character %d: %v", character.ID(), err)
	})
	return nil
}
