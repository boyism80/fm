package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

const commonSkillScriptPath = "script/skill/common.lua"

func CallSkillHook(ctx *core.ClientContext, character *entity.Character, skillID uint32, hook string) bool {
	if skillID == 0 || character == nil {
		return true
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return true
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return false
	}

	skillEntry := character.Skills.Get(skillID)
	if skillEntry == nil {
		return false
	}

	commonThread, commonErr := luax.NewThread(root, commonSkillScriptPath)
	if commonErr != nil {
		log.Printf("Skill common %s: %v", hook, commonErr)
		return false
	}
	commonResult, commonErr := luax.Call(commonThread, hook, character, skillEntry)
	if commonErr != nil {
		log.Printf("Skill common %s: %v", hook, commonErr)
		return false
	}
	if commonResult != nil && commonResult.Type() == lua.LTBool && !lua.LVAsBool(commonResult) {
		return false
	}

	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Skill hook %s failed for %s: %v", hook, scriptPath, err)
		return true
	}
	result, err := luax.Call(thread, luax.SkillScriptHookName(hook, skillID), character, skillEntry)
	if err != nil {
		log.Printf("Skill hook %s failed for %s: %v", hook, scriptPath, err)
		return true
	}
	if result != nil && result.Type() == lua.LTBool && !lua.LVAsBool(result) {
		return false
	}
	return true
}

func CallPassiveSkillHook(ctx *core.ClientContext, character *entity.Character, skillID uint32, hook string) {
	if skillID == 0 || character == nil {
		return
	}
	var root *lua.LState
	if m := character.GetMap(); m != nil {
		root = m.GetLuaRoot()
	}
	if root == nil {
		return
	}
	skillEntry := character.Skills.Get(skillID)
	if skillEntry == nil {
		return
	}
	commonThread, commonErr := luax.NewThread(root, commonSkillScriptPath)
	if commonErr != nil {
		log.Printf("Skill common %s: %v", hook, commonErr)
		return
	}
	commonResult, commonErr := luax.Call(commonThread, hook, character, skillEntry)
	if commonErr != nil {
		log.Printf("Skill common %s: %v", hook, commonErr)
		return
	}
	if commonResult != nil && commonResult.Type() == lua.LTBool && !lua.LVAsBool(commonResult) {
		return
	}
	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Skill passive hook %s failed for %s: %v", hook, scriptPath, err)
		return
	}
	result, err := luax.Call(thread, luax.SkillScriptHookName(hook, skillID), character, skillEntry)
	if err != nil {
		log.Printf("Skill passive hook %s failed for %s: %v", hook, scriptPath, err)
		return
	}
	if result != nil && result.Type() == lua.LTBool && !lua.LVAsBool(result) {
		return
	}
}

func CallOnAttackHooks(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair, skillID uint32, ranged bool, consumeSlot uint16) {
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}

	commonThread, err := luax.NewThread(root, commonSkillScriptPath)
	if err != nil {
		log.Printf("common.lua: %v", err)
		return
	}
	defer commonThread.Close()

	var skillLV lua.LValue = lua.LNil
	if skillID != 0 {
		if skillEntry := character.Skills.Get(skillID); skillEntry != nil {
			skillLV = luax.NewLuable(commonThread, skillEntry)
		}
	}
	damagesTable := buildDamagesTable(commonThread, mapInstance, damages)
	attackInfoTable := buildAttackInfoTable(commonThread, ranged, consumeSlot)
	f := commonThread.GetGlobal("on_attack")
	if f.Type() == lua.LTFunction {
		commonThread.Push(f)
		commonThread.Push(luax.NewLuable(commonThread, character))
		commonThread.Push(skillLV)
		commonThread.Push(damagesTable)
		commonThread.Push(attackInfoTable)
		if err := commonThread.PCall(4, 0, nil); err != nil {
			log.Printf("common on_attack: %v", err)
			return
		}
		readDamagesFromLuaTableInto(damagesTable, damages)
	}

	if skillID == 0 {
		return
	}
	skillEntry := character.Skills.Get(skillID)
	if skillEntry == nil {
		return
	}

	skillThread, err := luax.NewThread(root, fmt.Sprintf("script/skill/%d.lua", skillID))
	if err != nil {
		return
	}
	defer skillThread.Close()
	f2 := skillThread.GetGlobal(luax.SkillScriptHookName("on_attack", skillID))
	if f2.Type() != lua.LTFunction {
		return
	}
	damagesTable2 := buildDamagesTable(skillThread, mapInstance, damages)
	skillThread.Push(f2)
	skillThread.Push(luax.NewLuable(skillThread, character))
	skillThread.Push(luax.NewLuable(skillThread, skillEntry))
	skillThread.Push(damagesTable2)
	if skillThread.PCall(3, 1, nil) == nil {
		readDamagesFromLuaTableInto(damagesTable2, damages)
	}
}

func CallSummonOnAttackHooks(ctx *core.ClientContext, character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair, skillID uint32) {
	if ctx == nil || character == nil || mapInstance == nil {
		return
	}
	if skillID == 0 {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	skillEntry := character.Skills.Get(skillID)
	if skillEntry == nil {
		return
	}
	skillThread, err := luax.NewThread(root, fmt.Sprintf("script/skill/%d.lua", skillID))
	if err != nil {
		log.Printf("summon on_attack thread %d: %v", skillID, err)
		return
	}
	defer skillThread.Close()
	f := skillThread.GetGlobal(luax.SkillScriptHookName("on_attack", skillID))
	if f.Type() != lua.LTFunction {
		return
	}
	damagesTable := buildDamagesTable(skillThread, mapInstance, damages)
	skillThread.Push(f)
	skillThread.Push(luax.NewLuable(skillThread, character))
	skillThread.Push(luax.NewLuable(skillThread, skillEntry))
	skillThread.Push(damagesTable)
	if err := skillThread.PCall(3, 0, nil); err != nil {
		log.Printf("summon on_attack %d: %v", skillID, err)
	}
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
