package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
	"github.com/boyism80/fm/types"
)

type SelectCharacter struct {
	ls *LoginServer
}

func (SelectCharacter) New(ls *LoginServer) *SelectCharacter {
	return &SelectCharacter{
		ls: ls,
	}
}

func (h *SelectCharacter) Handle(ctx *core.ClientContext, req *request.SelectCharacter) error {
	log.Printf("Select character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.CharacterId)

	failed := &response.LoginFailed{Reason: response.LoginFailedReasonSystemError6}
	loginClient := ctx.Client.(*client.LoginClient)
	accountId := loginClient.GetAccountId()
	if accountId == 0 {
		_ = ctx.Client.Send(failed, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("select character: account id is missing")
	}
	worldId := loginClient.GetWorldId()
	channelId := loginClient.GetChannelId()
	worldRoutes, ok := h.ls.channelRoutes[worldId]
	if !ok {
		_ = ctx.Client.Send(failed, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("select character: route not found for world=%d channel=%d", worldId, channelId)
	}
	route, ok := worldRoutes[uint32(channelId)]
	if !ok {
		_ = ctx.Client.Send(failed, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("select character: route not found for world=%d channel=%d", worldId, channelId)
	}
	transferResp := &response.Transfer{
		IP:          route.GetHost(),
		Port:        uint16(route.GetPort()),
		CharacterId: req.CharacterId,
	}
	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.BeginGameTransitionReply, error) {
			return h.ls.internalClient.BeginGameTransition(c, &internal.BeginGameTransitionRequest{
				WorldId:     worldId,
				AccountId:   accountId,
				CharacterId: req.CharacterId,
				ClientIp:    ctx.Client.GetRemoteIP(),
			})
		}, func(reply *internal.BeginGameTransitionReply) error {
			if !reply.GetOk() {
				return fmt.Errorf("begin game transition failed: %v", reply.GetErrorCode())
			}
			loginClient.SetTransferDisconnect(true)
			if err := ctx.Client.Send(transferResp, types.SEND_POLICY_ENCRYPT); err != nil {
				log.Printf("Failed to send transfer response: %v", err)
				return err
			}
			return nil
		}).OnError(func(err error) {
		log.Printf("SelectCharacter (async): %v", err)
		_ = ctx.Client.Send(failed, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}
