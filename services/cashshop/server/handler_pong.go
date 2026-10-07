package server

import (
	"context"
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/async"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/cashshop/client"
)

type Pong struct {
	cs *CashShopServer
}

func (Pong) New(cs *CashShopServer) *Pong {
	return &Pong{cs: cs}
}

func (h *Pong) Handle(ctx *core.ClientContext, req *request.Pong) error {
	csClient, ok := ctx.Client.(*client.CashShopClient)
	if !ok {
		return nil
	}
	csClient.MarkPongReceived()
	character := csClient.GetCharacter()
	if character == nil {
		return nil
	}

	characterID := character.ID()
	cashShopID := h.cs.cashShopID()
	async.NewTask(ctx.ActorContext, core.InternalRPCPerStepTimeout).ThenRPC(func(c context.Context) (*internal.RefreshSessionReply, error) {
		return h.cs.internalClient.RefreshSession(c, &internal.RefreshSessionRequest{
			WorldId:     h.cs.worldID(),
			AccountId:   character.AccountID(),
			CharacterId: &characterID,
			CashShopId:  &cashShopID,
		})
	}, func(reply *internal.RefreshSessionReply) error {
		if reply.GetOk() {
			return nil
		}
		if csClient.Leaving() == false {
			_ = csClient.GetConnection().Close()
		}
		return fmt.Errorf("refresh session failed: %s", reply.GetErrorCode())
	}).OnError(func(err error) {
		log.Printf("cash shop pong refresh (async): %v", err)
	})
	return nil
}
