package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/services/game/client"
	lua "github.com/yuin/gopher-lua"
)

type SummonSkill struct {
	gs *GameServer
}

func (SummonSkill) New(gs *GameServer) *SummonSkill {
	return &SummonSkill{
		gs: gs,
	}
}

func (h *SummonSkill) Handle(ctx *core.ClientContext, req *request.SummonSkill) error {
	gameClient, ok := ctx.Client.(*client.GameClient)
	if !ok {
		return nil
	}
	character := gameClient.GetCharacter()
	if character == nil {
		return nil
	}
	if req.SubSkillID == 0 {
		return nil
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return nil
	}
	summon := mapInstance.GetSummon(req.SummonOID)
	if summon == nil || summon.OwnerID != character.GetID() {
		return nil
	}
	skillEntry := character.Skills.Get(req.SubSkillID)
	if skillEntry == nil {
		return nil
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return nil
	}
	params := buildSummonSkillParams(root, req)
	scriptPath := fmt.Sprintf("script/skill/%d.lua", req.SubSkillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("summon skill script %s: %v", scriptPath, err)
		return nil
	}
	_, err = luax.Call(thread, luax.SkillScriptHookName("on_summon_skill", req.SubSkillID), character, summon, skillEntry, params)
	if err != nil {
		log.Printf("summon skill script %s: %v", scriptPath, err)
	}
	return nil
}

func buildSummonSkillParams(L *lua.LState, req *request.SummonSkill) lua.LValue {
	t := L.NewTable()
	t.RawSetString("buff_effect_index", lua.LNumber(req.BuffEffectIndex))
	return t
}
