package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

func (m *Map) callMobReviveScript(mob *Mob, pos types.Point[int16], linkOID uint32, revives []uint32) bool {
	if m == nil || mob == nil || mob.Wz == nil {
		return false
	}
	root := m.GetLuaRoot()
	if root == nil {
		return false
	}

	mobID := mob.Wz.ID
	scriptPath := fmt.Sprintf("script/mob/%d.lua", mobID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return false
	}

	hook := fmt.Sprintf("on_revive_%d", mobID)
	if thread.GetGlobal(hook).Type() != lua.LTFunction {
		thread.Close()
		return false
	}

	revivesTbl := thread.NewTable()
	for i, id := range revives {
		revivesTbl.RawSetInt(i+1, lua.LNumber(id))
	}

	if _, err := luax.Call(thread, hook, mob, m, pos.X, pos.Y, linkOID, revivesTbl); err != nil {
		log.Printf("mob revive script %s: %v", hook, err)
	}
	return true
}
