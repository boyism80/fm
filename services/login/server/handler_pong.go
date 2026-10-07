package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	loginclient "github.com/boyism80/fm/services/login/client"
)

type Pong struct {
	ls *LoginServer
}

func (Pong) New(ls *LoginServer) *Pong {
	return &Pong{
		ls: ls,
	}
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	c := ctx.Client.(*loginclient.LoginClient)
	c.MarkPongReceived()
	accountId := c.GetAccountId()
	worldId := c.GetWorldId()
	if accountId != 0 {
		async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
			func(cctx context.Context) (*internal.RefreshSessionReply, error) {
				return h.ls.internalClient.RefreshSession(cctx, &internal.RefreshSessionRequest{
					WorldId:   worldId,
					AccountId: accountId,
				})
			}, func(reply *internal.RefreshSessionReply) error {
				if !reply.GetOk() {
					return fmt.Errorf("refresh session failed: %v", reply.GetErrorCode())
				}
				return nil
			}).OnError(func(err error) {
			log.Printf("Login Pong refresh (async): %v", err)
		})
	}
	log.Printf("Pong packet received from %s", ctx.Client.GetConnection().RemoteAddr())
	return nil
}
