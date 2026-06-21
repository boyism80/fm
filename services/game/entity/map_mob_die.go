package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

func (m *Map) callMobDieScript(mob *Mob, attacker *Character) {
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
	if thread, err := luax.NewThread(root, scriptPath); err == nil {
		hook := fmt.Sprintf("on_mob_die_%d", mobID)
		if _, err := luax.Call(thread, hook, mob, attackerArg, m); err != nil {
			log.Printf("mob die script %s: %v", hook, err)
		}
	}

	if thread, err := luax.NewThread(root, constant.CharacterHookScriptPath); err == nil {
		if _, err := luax.Call(thread, "on_mob_die", mob, attackerArg, m); err != nil {
			log.Printf("on_mob_die script: %v", err)
		}
	}
}
