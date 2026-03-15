package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/game/entity"
	"github.com/boyism80/fm/protocol/dto"
	lua "github.com/yuin/gopher-lua"
)

func CallSkillOnAttack(ctx *core.ClientContext, character *entity.Character, skillID uint32, mapInstance *entity.Map, damages []dto.AttackPair) {
	if skillID == 0 || ctx.LogicActorPID == nil {
		return
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return
	}
	skillEntry := character.Skills[skillID]
	if skillEntry == nil {
		return
	}
	thread, err := luax.NewThread(root, fmt.Sprintf("script/skill/%d.lua", skillID))
	if err != nil {
		return
	}
	defer thread.Close()
	f := thread.GetGlobal("on_attack")
	if f.Type() != lua.LTFunction {
		return
	}
	damagesTable := buildDamagesTable(thread, mapInstance, damages)
	thread.Push(f)
	thread.Push(luax.NewLuable(thread, character))
	thread.Push(luax.NewLuable(thread, skillEntry))
	thread.Push(damagesTable)
	if thread.PCall(3, 1, nil) == nil {
		readDamagesFromLuaTableInto(damagesTable, damages)
	}
}

func CallOnAttackScript(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair, skillID uint32, ranged bool, consumeSlot uint16) {
	if ctx.LogicActorPID == nil {
		return
	}
	root := luax.GetRootLuaState(ctx.LogicActorPID.String())
	if root == nil {
		return
	}
	thread, err := luax.NewThread(root, "script/script.lua")
	if err != nil {
		log.Printf("Failed to load script: %v", err)
		return
	}
	defer thread.Close()
	f := thread.GetGlobal("on_attack")
	if f.Type() != lua.LTFunction {
		return
	}
	var skillLV lua.LValue = lua.LNil
	if skillID != 0 {
		if skillEntry := character.Skills[skillID]; skillEntry != nil {
			skillLV = luax.NewLuable(thread, skillEntry)
		}
	}
	damagesTable := buildDamagesTable(thread, mapInstance, damages)
	attackInfoTable := buildAttackInfoTable(thread, ranged, consumeSlot)
	thread.Push(f)
	thread.Push(luax.NewLuable(thread, character))
	thread.Push(skillLV)
	thread.Push(damagesTable)
	thread.Push(attackInfoTable)
	if err := thread.PCall(4, 1, nil); err != nil {
		log.Printf("Failed to call script on_attack: %v", err)
		return
	}
	thread.Pop(1)
}

func ApplyDamageToMobs(character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair) {
	for _, damage := range damages {
		mob := mapInstance.GetMob(damage.OID)
		if mob == nil {
			log.Printf("Mob not found for OID: %d", damage.OID)
			continue
		}
		for _, damagePair := range damage.DamagePairs {
			mob.ApplyDamage(character, uint32(damagePair.Damage))
		}
	}
}
