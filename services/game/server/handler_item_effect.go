package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type ItemEffect struct{}

func (ItemEffect) New(*GameServer) *ItemEffect { return &ItemEffect{} }
func (*ItemEffect) Handle(ctx *core.ClientContext, req *request.ItemEffect) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("item effect: client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.SetItemEffect(req.ItemID)
	return nil
}
