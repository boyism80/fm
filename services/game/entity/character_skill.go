package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) DistributeSP(skillID uint32) bool {
	if ch.SkillPoint == 0 || ch.canLearn(skillID) == false {
		return false
	}

	skillEntry := ch.Skills.Get(skillID)
	if skillEntry == nil {
		if ch.GameWorld == nil {
			return false
		}
		resources := ch.GameWorld.GetResources()
		if resources == nil {
			return false
		}
		wzSkill := resources.GetSkill(skillID)
		if wzSkill == nil {
			return false
		}
		skillEntry = NewSkillEntry(ch, wzSkill, 0, wzSkill.DefaultMasterLevel())
	}

	if skillEntry.Level() >= min(skillEntry.MasterLevel, skillEntry.Wz.MaxLevel) {
		return false
	}

	ch.SkillPoint--
	if skillEntry.Level() == 0 {
		skillEntry.SetLevel(1)
		ch.Skills.Bind(skillID, skillEntry)
	} else {
		skillEntry.SetLevel(skillEntry.Level() + 1)
	}
	ch.Listener.OnUpdateStats(ch, map[constant.Stat]int32{
		constant.StatAvailableSP: int32(ch.SkillPoint),
	}, false)
	return true
}

func (ch *Character) canLearn(skillID uint32) bool {
	job := uint32(ch.Class)
	skillJob := skillID / 10000
	switch skillJob {
	case uint32(constant.ClassBeginner):
		return job < uint32(constant.ClassNoblesse)
	case uint32(constant.ClassNoblesse):
		return job >= uint32(constant.ClassNoblesse) && job < uint32(constant.ClassLegend)
	case uint32(constant.ClassLegend):
		return job >= uint32(constant.ClassLegend) && job <= uint32(constant.ClassAran5)
	}

	if job/100 != skillJob/100 {
		return false
	}
	skillBranch := (skillJob / 10) % 10
	if skillBranch != 0 && skillBranch != (job/10)%10 {
		return false
	}
	return skillJob%10 <= job%10
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
