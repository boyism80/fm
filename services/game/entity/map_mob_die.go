package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (m *Map) runDieScript(mob *Mob, attacker *Character) {
	if m == nil || mob == nil || mob.Wz == nil {
		return
	}
	root := m.GetLuaRoot()
	if root == nil {
		return
	}

	var attackerArg interface{} = lua.LNil
	if attacker != nil {
		attackerArg = attacker
	}

	mobID := mob.Wz.ID
	scriptPath := fmt.Sprintf("script/mob/%d.lua", mobID)
	m.runMobLuaHook(root, scriptPath, fmt.Sprintf("on_mob_die_%d", mobID), mob, attackerArg, m)
	m.runMobLuaHook(root, constant.CharacterHookScriptPath, "on_mob_die", mob, attackerArg, m)
}

func (m *Map) runMobLuaHook(root *lua.LState, scriptPath, hook string, args ...interface{}) {
	if m == nil || root == nil || hook == "" {
		return
	}
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		MapActorPID: m.GetActorPID(),
	})
	state, _, err := luax.Execute(root, thread, hook, args...)
	if err != nil {
		log.Printf("mob script %s: %v", hook, err)
	}
	if state == lua.ResumeYield {
		return
	}
}
