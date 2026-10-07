package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/cashshop/client"
	"github.com/boyism80/fm/types"
)

type Refresh struct {
	cs *CashShopServer
}

func (Refresh) New(cs *CashShopServer) *Refresh {
	return &Refresh{cs: cs}
}

func (h *Refresh) Handle(ctx *core.ClientContext, req *request.CashShopRefresh) error {
	csClient, ok := ctx.Client.(*client.CashShopClient)
	if !ok {
		return nil
	}
	character := csClient.GetCharacter()
	if character == nil {
		return nil
	}

	return ctx.Client.Send(&response.CashShopBalance{NXCash: character.NXCash, MaplePoint: character.MaplePoint}, types.SEND_POLICY_ENCRYPT)
}
