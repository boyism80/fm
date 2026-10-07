package server

import (
	"context"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/login/client"
	"github.com/boyism80/fm/types"
)

type DeleteCharacter struct {
	ls *LoginServer
}

func (DeleteCharacter) New(ls *LoginServer) *DeleteCharacter {
	return &DeleteCharacter{
		ls: ls,
	}
}

func (h *DeleteCharacter) Handle(ctx *core.ClientContext, req *request.DeleteCharacter) error {
	log.Printf("Delete character packet received from %s - Character ID: %d",
		ctx.Client.GetConnection().RemoteAddr(), req.ID)

	accountId := ctx.Client.(*client.LoginClient).GetAccountId()
	if accountId == 0 {
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	}

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.DeleteCharacterReply, error) {
			return h.ls.internalClient.DeleteCharacter(c, &internal.DeleteCharacterRequest{AccountId: accountId, CharacterId: req.ID})
		}, func(reply *internal.DeleteCharacterReply) error {
			deleteResp := &response.DeleteCharacter{
				ID:      req.ID,
				Success: reply.Success,
			}
			return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
		}).OnError(func(err error) {
		log.Printf("DeleteCharacter (async): %v", err)
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		_ = ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}
