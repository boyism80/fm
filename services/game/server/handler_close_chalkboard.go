package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type CloseChalkboard struct{}

func (CloseChalkboard) New(*GameServer) *CloseChalkboard { return &CloseChalkboard{} }
func (*CloseChalkboard) Handle(ctx *core.ClientContext, _ *request.CloseChalkboard) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("close chalkboard: client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.CloseChalkboard()
	return nil
}
