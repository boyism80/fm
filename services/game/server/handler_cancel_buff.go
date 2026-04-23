package server

import (
	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	"github.com/boyism80/fm/services/game/constant"
)

type CancelBuff struct {
	gs *GameServer
}

func (CancelBuff) New(gs *GameServer) *CancelBuff {
	return &CancelBuff{
		gs: gs,
	}
}

func (h *CancelBuff) Handle(ctx *core.ClientContext, req *request.CancelBuff) error {
	client, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := client.GetCharacter()
	if character == nil {
		return nil
	}
	sourceID := req.SourceID
	if sourceID == 0 {
		return nil
	}
	if sourceID < 0 {
		sourceID = -sourceID
	}
	flagSet := make(map[constant.BuffFlag]struct{})
	for _, buff := range character.Buffs.Entities() {
		if buff == nil {
			continue
		}
		buffID := buff.GetBuffID()
		if buffID == 0 {
			continue
		}
		if buffID < 0 {
			buffID = -buffID
		}
		if buffID != sourceID {
			continue
		}
		for _, flag := range buff.GetFlags() {
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
