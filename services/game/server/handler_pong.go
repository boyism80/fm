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
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(
		func(cctx context.Context) (*internal.RefreshSessionReply, error) {
			req := &internal.RefreshSessionRequest{
				WorldId:   h.gs.config.WorldId,
				AccountId: character.AccountID,
			}
			if cid := character.GetID(); cid != 0 {
				v := cid
				req.CharacterId = &v
				ch := h.gs.config.ChannelId
				req.ChannelId = &ch
			}
			return h.gs.internalClient.RefreshSession(cctx, req)
		}, func(reply *internal.RefreshSessionReply) error {
			if reply.GetOk() {
				return nil
			}
			switch reply.GetErrorCode() {
			case internal.SessionErrorCode_SESSION_NOT_FOUND:
				if c.GetCharacter() == character {
					_ = c.GetConnection().Close()
				}
			case internal.SessionErrorCode_SESSION_NOT_OWNER:
				if c.GetCharacter() == character {
					c.LoseSession()
					_ = c.GetConnection().Close()
				}
			}
			return fmt.Errorf("refresh session failed: %s", reply.GetErrorCode())
		}).OnError(func(err error) {
		log.Printf("Game Pong refresh (async): %v", err)
	})
	return nil
}
