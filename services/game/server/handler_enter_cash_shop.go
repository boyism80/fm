package server

import (
	"context"
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/client"
)

type EnterCashShop struct {
	gs *GameServer
}

func (EnterCashShop) New(gs *GameServer) *EnterCashShop {
	return &EnterCashShop{
		gs: gs,
	}
}

func (h *EnterCashShop) Handle(ctx *core.ClientContext, req *request.EnterCashShop) error {
	if ctx.ActorContext == nil {
		return fmt.Errorf("enter cash shop: actor context required for internal RPC")
	}

	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("enter cash shop: invalid client type")
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return fmt.Errorf("enter cash shop: character not loaded")
	}

	worldID := h.gs.config.WorldId
	h.gs.handOver(ctx, gameClient, character, constant.ServerBlockedCashShopUnavailable, func(c context.Context) (*response.SwitchChannel, error) {
		reply, err := h.gs.internalClient.FindCashShop(c, &internal.FindCashShopRequest{WorldId: worldID})
		if err != nil {
			return nil, err
		}
		if !reply.GetFound() {
			return nil, nil
		}
		return &response.SwitchChannel{IP: reply.GetHost(), Port: uint16(reply.GetPort())}, nil
	})
	return nil
}
