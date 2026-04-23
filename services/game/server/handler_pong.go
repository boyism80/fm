package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	gameclient "github.com/boyism80/fm/services/game/client"
)

type Pong struct {
	gs *GameServer
}

func (Pong) New(gs *GameServer) *Pong {
	return &Pong{
		gs: gs,
	}
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	c, ok := ctx.Client.(*gameclient.GameClient)
	if !ok {
		return nil
	}
	c.MarkPongReceived()
	character := c.GetCharacter()
	if character == nil || character.AccountID == 0 || h.gs.internalClient == nil {
		return nil
	}
	if ctx.ActorContext == nil {
		return fmt.Errorf("game pong: actor context required for internal RPC")
	}
	async.ThenRPC(async.NewPromise(ctx.ActorContext, core.InternalRPCPerStepTimeout),
		func(cctx context.Context) (*internal.RefreshSessionReply, error) {
			return h.gs.internalClient.RefreshSession(cctx, &internal.RefreshSessionRequest{
				WorldId:   h.gs.config.WorldId,
				AccountId: character.AccountID,
			})
		}, func(reply *internal.RefreshSessionReply) error {
			if !reply.GetOk() {
				return fmt.Errorf("refresh session failed: %s", reply.GetErrorCode())
			}
			return nil
		}).
		OnError(func(err error) {
			log.Printf("Game Pong refresh (async): %v", err)
		}).
		Run()
	return nil
}
