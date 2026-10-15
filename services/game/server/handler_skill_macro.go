package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type SkillMacro struct{}

func (SkillMacro) New(*GameServer) *SkillMacro { return &SkillMacro{} }
func (*SkillMacro) Handle(ctx *core.ClientContext, req *request.SkillMacro) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("skill macro: client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		return nil
	}
	if req.Count != uint8(len(req.Entries)) {
		return nil
	}
	entries := make([]game.SkillMacro, 0, len(req.Entries))
	for _, entry := range req.Entries {
		if entry.Shout > 1 {
			return nil
		}
		entries = append(entries, game.SkillMacro{Name: entry.Name, Shout: entry.Shout == 1, Skills: entry.Skills})
	}
	ch.SkillMacros.Replace(entries)
	return nil
}
