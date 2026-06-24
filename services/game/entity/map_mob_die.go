package entity

import (
	"log"

	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

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
