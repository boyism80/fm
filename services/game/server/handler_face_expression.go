package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type FaceExpression struct{}

func (FaceExpression) New(*GameServer) *FaceExpression { return &FaceExpression{} }
func (*FaceExpression) Handle(ctx *core.ClientContext, req *request.FaceExpression) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("face expression: client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.Express(req.Expression)
	return nil
}
