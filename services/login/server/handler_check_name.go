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

type CheckName struct {
	ls *LoginServer
}

func (CheckName) New(ls *LoginServer) *CheckName {
	return &CheckName{
		ls: ls,
	}
}

func (h *CheckName) Handle(ctx *core.ClientContext, req *request.CheckName) error {
	log.Printf("Check name packet received from %s - Name: %s",
		ctx.Client.GetConnection().RemoteAddr(), req.Name)

	reqMsg := &internal.CheckCharacterNameRequest{Name: req.Name, AccountId: ctx.Client.(*client.LoginClient).GetAccountId()}

	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(c context.Context) (*internal.CheckCharacterNameReply, error) {
			return h.ls.internalClient.CheckCharacterName(c, reqMsg)
		}, func(reply *internal.CheckCharacterNameReply) error {
			checkResp := &response.CheckName{
				Name:   req.Name,
				Exists: reply.Exists,
			}
			return ctx.Client.Send(checkResp, types.SEND_POLICY_ENCRYPT)
		}).OnError(func(err error) {
		log.Printf("CheckName (async): %v", err)
		checkResp := &response.CheckName{Name: req.Name, Exists: false}
		_ = ctx.Client.Send(checkResp, types.SEND_POLICY_ENCRYPT)
	})
	return nil
}
