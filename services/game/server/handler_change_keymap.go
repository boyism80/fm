package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/entity"
)

type ChangeKeymap struct {
	gs *GameServer
}

func (ChangeKeymap) New(gs *GameServer) *ChangeKeymap {
	return &ChangeKeymap{
		gs: gs,
	}
}

func (h *ChangeKeymap) Handle(ctx *core.ClientContext, req *request.ChangeKeymap) error {
	gc, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return fmt.Errorf("change keymap: client is not a GameClient")
	}
	ch := gc.GetCharacter()
	if ch == nil {
		return fmt.Errorf("change keymap: no character")
	}
	kl := ch.KeyLayout()
	if kl == nil {
		return fmt.Errorf("change keymap: nil key layout")
	}

	if len(req.Changes) == 0 {
		_, _ = req.Type, req.Data
		return nil
	}

	for _, c := range req.Changes {
		if !h.allowBinding(ch, c.Type, c.Action) {
			continue
		}
		kl.ApplyChange(int(c.Key), c.Type, c.Action)
	}
	return nil
}

func (h *ChangeKeymap) allowBinding(ch *entity.Character, typ byte, action int32) bool {
	if typ != 1 || action < 1000 {
		return true
	}
	if ch == nil {
		return true
	}
	res := ch.GameWorld.GetResources()
	if res == nil {
		return true
	}
	skillID := uint32(action)
	skill, ok := res.Skills[skillID]
	if !ok || skill == nil {
		return true
	}
	if !skill.Invisible {
		return true
	}
	ent := ch.Skills.Get(skillID)
	return ent != nil && ent.Level() > 0
}
