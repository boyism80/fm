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
	ls     *LoginServer
	opcode byte
}

func (Pong) New(ls *LoginServer) *Pong {
	return &Pong{
		ls:     ls,
		opcode: 0x0A,
	}
}

func (h *Pong) GetOpcode() byte {
	return h.opcode
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	c, ok := ctx.Client.(*loginclient.LoginClient)
	if !ok {
		return nil
	}
	c.MarkPongReceived()
	accountId := c.GetAccountId()
	worldId := c.GetWorldId()
	if accountId != 0 && h.ls.context.InternalClient != nil {
		if ctx.ActorContext == nil {
			return fmt.Errorf("login pong: actor context required for internal RPC")
		}
		async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
			func(cctx context.Context) (*internal.RefreshSessionReply, error) {
				return h.ls.context.InternalClient.RefreshSession(cctx, &internal.RefreshSessionRequest{
					WorldId:   worldId,
					AccountId: accountId,
				})
			}, func(reply *internal.RefreshSessionReply) error {
				if !reply.GetOk() {
					return fmt.Errorf("refresh session failed: %v", reply.GetErrorCode())
				}
				return nil
			}).
			OnError(func(err error) {
				log.Printf("Login Pong refresh (async): %v", err)
			}).
			Run()
	}
	log.Printf("Pong packet received from %s", ctx.Client.GetConnection().RemoteAddr())
	return nil
}
