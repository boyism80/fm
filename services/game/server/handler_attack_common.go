package server

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/entity"
	lua "github.com/yuin/gopher-lua"
)

func skillHookAllowed(ret lua.LValue, err error, hook string, defaultAllow bool) bool {
	if err != nil {
		log.Printf("Skill hook %s: %v", hook, err)
		return false
	}
	if ret == nil {
		return defaultAllow
	}
	if ret.Type() == lua.LTBool && !lua.LVAsBool(ret) {
		return false
	}
	return defaultAllow
}

func CallSkillHook(character *entity.Character, skillID uint32, hook string) bool {
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
	commonThread, err := luax.NewThread(root, constant.SkillHookScriptPath)
	if err != nil {
		log.Printf("Skill common %s: %v", hook, err)
		return false
	}
	commonRet, err := luax.Call(commonThread, hook, character, skillEntry)
	if !skillHookAllowed(commonRet, err, hook, true) {
		return false
	}
	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Skill hook %s failed for %s: %v", hook, scriptPath, err)
		return true
	}
	skillHook := fmt.Sprintf("%s_%d", hook, skillID)
	skillRet, err := luax.Call(thread, skillHook, character, skillEntry)
	return skillHookAllowed(skillRet, err, skillHook, true)
}

func CallPassiveSkillHook(character *entity.Character, skillID uint32, hook string) {
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
	commonThread, err := luax.NewThread(root, constant.SkillHookScriptPath)
	if err != nil {
		log.Printf("Skill common %s: %v", hook, err)
		return
	}
	commonRet, err := luax.Call(commonThread, hook, character, skillEntry)
	if !skillHookAllowed(commonRet, err, hook, true) {
		return
	}
	scriptPath := fmt.Sprintf("script/skill/%d.lua", skillID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		log.Printf("Skill passive hook %s failed for %s: %v", hook, scriptPath, err)
		return
	}
	skillHook := fmt.Sprintf("%s_%d", hook, skillID)
	if _, err := luax.Call(thread, skillHook, character, skillEntry); err != nil {
		log.Printf("Skill passive hook %s failed for %s: %v", skillHook, scriptPath, err)
	}
}

func CallOnAttackHooks(character *entity.Character, damages []dto.AttackPair, skillID uint32, magicAttack bool, ranged bool, consumeSlot uint16) {
	if character == nil {
		return
	}
	mapInstance := character.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	commonThread, err := luax.NewThread(root, constant.SkillHookScriptPath)
	if err != nil {
		log.Printf("skill hook: %v", err)
		return
	}
	var skillLV lua.LValue = lua.LNil
	if skillID != 0 {
		if skillEntry := character.Skills.Get(skillID); skillEntry != nil {
			skillLV = luax.NewLuable(commonThread, skillEntry)
		}
	}
	damagesTable := buildDamagesTable(commonThread, character, damages)
	attackInfoTable := buildAttackInfoTable(commonThread, magicAttack, ranged, consumeSlot)
	if _, err := luax.Call(commonThread, "on_attack", character, skillLV, damagesTable, attackInfoTable); err != nil {
		log.Printf("skill hook on_attack: %v", err)
		return
	}
	readDamagesFromLuaTableInto(damagesTable, damages)
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
	damagesTable2 := buildDamagesTable(skillThread, character, damages)
	skillHook := "on_attack"
	if _, err := luax.Call(skillThread, skillHook, character, skillEntry, damagesTable2); err == nil {
		readDamagesFromLuaTableInto(damagesTable2, damages)
	}
}

func CallSummonOnAttackHooks(character *entity.Character, mapInstance *entity.Map, damages []dto.AttackPair, skillID uint32) {
	if character == nil || mapInstance == nil || skillID == 0 {
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
	damagesTable := buildDamagesTable(skillThread, character, damages)
	skillHook := "on_attack"
	if _, err := luax.Call(skillThread, skillHook, character, skillEntry, damagesTable); err != nil {
		log.Printf("summon on_attack %d: %v", skillID, err)
	}
}
