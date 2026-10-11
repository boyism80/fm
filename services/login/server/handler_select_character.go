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
	return h.selectCharacter(ctx, req.CharacterId, nil)
}

func (h *SelectCharacter) selectCharacter(ctx *core.ClientContext, characterID uint32, secondPassword *string) error {
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
		CharacterId: characterID,
	}
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.AuthorizeCharacterSelectionReply, error) {
			return h.ls.internalClient.AuthorizeCharacterSelection(c, &internal.AuthorizeCharacterSelectionRequest{
				WorldId:        worldId,
				AccountId:      accountId,
				CharacterId:    characterID,
				ClientIp:       ctx.Client.GetRemoteIP(),
				SecondPassword: secondPassword,
			})
		}, func(reply *internal.AuthorizeCharacterSelectionReply) error {
			if reply.GetStatus() != internal.SecondPasswordStatus_SECOND_PASSWORD_SUCCESS {
				return ctx.Client.Send(&response.SecondPasswordResult{
					Mode:   request.SecondPasswordModeSet,
					Result: response.SecondPasswordResultInvalid,
				}, types.SEND_POLICY_ENCRYPT)
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
