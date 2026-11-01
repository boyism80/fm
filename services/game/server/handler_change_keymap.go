package server

import (
	"fmt"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game"
	"github.com/boyism80/fm/services/game/client"
)

type ChangeKeymap struct {
	gs *GameServer
}

func (ChangeKeymap) New(gs *GameServer) *ChangeKeymap {
	return &ChangeKeymap{
		gs: gs,
	}
}

// petPotion is the item the pets drink, or 0 to clear it when the client sends anything but a consume item.
func (h *ChangeKeymap) petPotion(data int32) uint32 {
	if data <= 0 {
		return 0
	}
	model, ok := h.gs.GetResources().Items[uint32(data)]
	if !ok || model.GetInventoryType() != constant.InventoryTypeConsume {
		return 0
	}
	return uint32(data)
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
	switch req.Mode {
	case 1:
		ch.Pets().SetHPItem(h.petPotion(req.Data))
		return nil
	case 2:
		ch.Pets().SetMPItem(h.petPotion(req.Data))
		return nil
	}

	for _, c := range req.Changes {
		if !h.allowBinding(ch, c.Type, c.Action) {
			continue
		}
		ch.KeyLayout().SetKey(int(c.Key), c.Type, c.Action)
	}
	return nil
}

func (h *ChangeKeymap) allowBinding(ch *game.Character, typ byte, action int32) bool {
	if typ != 1 || action < 1000 {
		return true
	}
	if ch == nil {
		return true
	}
	res := ch.GetGameWorld().GetResources()
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
	ent := ch.Skills().Get(skillID)
	return ent != nil && ent.Level() > 0
}
