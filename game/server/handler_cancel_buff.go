package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/game/client"
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/protocol/request"
)

type CancelBuff struct {
	gs     *GameServer
	opcode byte
}

func (CancelBuff) New(gs *GameServer) *CancelBuff {
	return &CancelBuff{
		gs:     gs,
		opcode: 0x4B,
	}
}

func (h *CancelBuff) GetOpcode() byte {
	return h.opcode
}

func (h *CancelBuff) Handle(ctx *core.ClientContext, req *request.CancelBuff) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := client.GetCharacter()
	if character == nil || character.Buffs == nil {
		return nil
	}
	if req.SourceID <= 0 {
		return nil
	}

	skillID := uint32(req.SourceID)
	flagSet := make(map[constant.BuffFlag]struct{})
	for _, buff := range character.Buffs.Entities() {
		if buff == nil || buff.Wz == nil || buff.Wz.ID != skillID {
			continue
		}
		for _, flag := range buff.Flags {
			flagSet[flag] = struct{}{}
		}
	}
	if len(flagSet) == 0 {
		return nil
	}

	flags := make([]constant.BuffFlag, 0, len(flagSet))
	for flag := range flagSet {
		flags = append(flags, flag)
	}
	character.Buffs.RemoveBuff(flags)
	return nil
}
