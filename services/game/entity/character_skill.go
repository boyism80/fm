package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) canLearn(skillID uint32) bool {
	class := uint32(ch.Class)
	skillClass := skillID / 10000
	switch skillClass {
	case uint32(constant.ClassBeginner):
		return class < uint32(constant.ClassNoblesse)
	case uint32(constant.ClassNoblesse):
		return class >= uint32(constant.ClassNoblesse) && class < uint32(constant.ClassLegend)
	case uint32(constant.ClassLegend):
		return class >= uint32(constant.ClassLegend) && class <= uint32(constant.ClassAran5)
	}

	if class/100 != skillClass/100 {
		return false
	}
	skillBranch := (skillClass / 10) % 10
	if skillBranch != 0 && skillBranch != (class/10)%10 {
		return false
	}
	return skillClass%10 <= class%10
}

func (ch *Character) CallSkillHook(ctx actor.Context, skill *SkillEntry, hook string, args ...interface{}) bool {
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return true
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return true
	}

	var skillArg interface{}
	paths := []string{constant.SkillHookScriptPath}
	if skill != nil {
		skillArg = skill
		paths = append(paths, fmt.Sprintf("script/skill/%d.lua", skill.Wz.ID))
	}
	callArgs := append([]interface{}{ch, skillArg}, args...)
	for _, path := range paths {
		thread, err := luax.NewThread(root, path)
		if err != nil {
			continue
		}
		if ctx != nil {
			luax.SetConfiguration(thread, luax.Configuration{ActorContext: ctx})
		}
		ret, err := luax.Call(thread, hook, callArgs...)
		if err != nil {
			log.Printf("skill hook %s %s: %v", path, hook, err)
			return false
		}
		if ret != nil && ret.Type() == lua.LTBool && lua.LVAsBool(ret) == false {
			return false
		}
	}
	return true
}

func (ch *Character) GetTotalSkillLevel(skillID uint32) int {
	entry := ch.Skills.Get(skillID)
	if entry == nil {
		return 0
	}
	return entry.Level()
}

func (ch *Character) UseAttackSkill(skillID uint32, activate func() bool) bool {
	if ch.GameWorld == nil {
		return false
	}
	resources := ch.GameWorld.GetResources()
	if resources == nil {
		return false
	}
	wzSkill := resources.GetSkill(skillID)
	if wzSkill == nil {
		log.Printf("Skill not found: %d", skillID)
		return false
	}
	skillLevel := ch.GetTotalSkillLevel(skillID)
	if skillLevel <= 0 {
		log.Printf("Character does not have skill %d or skill level is 0", skillID)
		return false
	}
	levelData := wzSkill.GetLevelData(skillLevel)
	if levelData == nil {
		return false
	}
	skillEntry := ch.Skills.Get(skillID)
	if levelData.Cooldown > 0 && (skillEntry == nil || skillEntry.IsCooling()) {
		return false
	}
	if activate() == false {
		return false
	}
	if levelData.Cooldown > 0 {
		skillEntry.StartCooldown(levelData.Cooldown)
	}
	return true
}
