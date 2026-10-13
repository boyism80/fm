package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

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
		if ret != nil && ret.Type() == lua.LTBool && !lua.LVAsBool(ret) {
			return false
		}
	}
	return true
}
