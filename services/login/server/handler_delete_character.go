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

	loginClient, ok := ctx.Client.(*client.LoginClient)
	if !ok {
		return nil
	}

	ic := h.ls.internalClient
	if ic == nil {
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
	}

	reqProto := &internal.DeleteCharacterRequest{
		AccountId:   loginClient.GetAccountId(),
		CharacterId: req.ID,
	}

	if ctx.ActorContext == nil {
		log.Printf("DeleteCharacter: no actor context, cannot run internal RPC")
		deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
		_ = ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
		return fmt.Errorf("delete character: actor context required for internal RPC")
	}

	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(c context.Context) (*internal.DeleteCharacterReply, error) {
			return ic.DeleteCharacter(c, reqProto)
		}, func(reply *internal.DeleteCharacterReply) error {
			deleteResp := &response.DeleteCharacter{
				ID:      req.ID,
				Success: reply.Success,
			}
			return ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
		}).
		OnError(func(err error) {
			log.Printf("DeleteCharacter (async): %v", err)
			deleteResp := &response.DeleteCharacter{ID: req.ID, Success: false}
			_ = ctx.Client.Send(deleteResp, types.SEND_POLICY_ENCRYPT)
		}).
		Run()
	return nil
}
