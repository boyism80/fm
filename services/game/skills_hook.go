package game

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	pconst "github.com/boyism80/fm/protocol/constant"
	lua "github.com/yuin/gopher-lua"
)

func (sc *Skills) UseActive(actx actor.Context, skillID uint32, level uint8, direction *uint8, params func(L *lua.LState) lua.LValue) bool {
	if sc.owner.Spectating() {
		return false
	}
	m := sc.owner.GetMap()
	if m == nil {
		return false
	}
	root := m.GetLuaRoot()
	if root == nil {
		return false
	}
	entry := sc.Get(skillID)
	if entry == nil || entry.Level() <= 0 || entry.Level() != int(level) {
		return false
	}
	levelData := entry.Wz.GetLevelData(entry.Level())
	if levelData == nil {
		return false
	}
	if levelData.Cooldown > 0 && entry.IsCooling() {
		return false
	}

	args := params(root)
	if !sc.CallHook(actx, entry, "can_activate", args) {
		return false
	}
	if !sc.CallHook(actx, entry, "on_activating", args) {
		return false
	}
	if levelData.ItemCon != 0 && levelData.ItemConNo > 0 {
		if !sc.owner.Inventory.RemoveByItemIDCount(uint32(levelData.ItemCon), uint16(levelData.ItemConNo)) {
			return false
		}
	}
	if !sc.CallHook(actx, entry, "on_activated", args) {
		return false
	}

	if levelData.Cooldown > 0 && !entry.Wz.ManualCooldown {
		entry.StartCooldown(levelData.Cooldown)
	}
	sc.owner.Listener.OnShowSkillEffect(sc.owner, pconst.SkillEffectTypeCast, skillID, level, direction)
	sc.owner.Listener.OnUpdateStats(sc.owner, nil, true)
	return true
}

func (sc *Skills) CallHook(ctx actor.Context, skill *SkillEntry, hook string, args ...interface{}) bool {
	mapInstance := sc.owner.GetMap()
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
	callArgs := append([]interface{}{sc.owner, skillArg}, args...)
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
		if ret != nil && ret.Type() == lua.LTBool && !lua.LVAsBool(ret) {
			return false
		}
	}
	return true
}
