package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
)

type SkillEffect struct{}

func (SkillEffect) New(*GameServer) *SkillEffect { return &SkillEffect{} }
func (*SkillEffect) Handle(ctx *core.ClientContext, req *request.SkillEffect) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("skill effect: client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		return nil
	}
	ch.ShowChargeSkillEffect(req.SkillID, req.Level, req.Flags, req.Speed)
	return nil
}
