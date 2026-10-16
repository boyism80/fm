package server

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func CallOnAttackHooks(ctx actor.Context, character *game.Character, damages []dto.AttackPair, skill *game.SkillEntry, magicAttack bool, ranged bool) {
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	damagesTable := buildDamagesTable(root, character, damages)
	attackInfoTable := root.NewTable()
	attackInfoTable.RawSetString("magic", lua.LBool(magicAttack))
	attackInfoTable.RawSetString("ranged", lua.LBool(ranged))
	character.Skills().CallHook(ctx, skill, "on_attack", damagesTable, attackInfoTable)
	readDamagesFromLuaTableInto(damagesTable, damages)
}

func CallSummonOnAttackHooks(character *game.Character, mapInstance *game.Map, damages []dto.AttackPair, skillID uint32) {
	if character == nil || mapInstance == nil || skillID == 0 {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	skillEntry := character.Skills().Get(skillID)
	if skillEntry == nil {
		return
	}
	skillThread, err := luax.NewThread(root, fmt.Sprintf("script/skill/%d.lua", skillID))
	if err != nil {
		log.Printf("summon on_attack thread %d: %v", skillID, err)
		return
	}
	damagesTable := buildDamagesTable(skillThread, character, damages)
	skillHook := "on_attack"
	if _, err := luax.Call(skillThread, skillHook, character, skillEntry, damagesTable); err != nil {
		log.Printf("summon on_attack %d: %v", skillID, err)
	}
}
