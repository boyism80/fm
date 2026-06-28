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

	ic := h.ls.internalClient
	if ic == nil {
		checkResp := &response.CheckName{Name: req.Name, Exists: false}
		return ctx.Client.Send(checkResp, types.SEND_POLICY_ENCRYPT)
	}

	reqMsg := &internal.CheckCharacterNameRequest{Name: req.Name}

	if ctx.ActorContext == nil {
		log.Printf("CheckName: no actor context, cannot run internal RPC")
		checkResp := &response.CheckName{Name: req.Name, Exists: false}
		_ = ctx.Client.Send(checkResp, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("check name: actor context required for internal RPC")
	}

	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.CheckCharacterNameReply, error) {
			return ic.CheckCharacterName(c, reqMsg)
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
