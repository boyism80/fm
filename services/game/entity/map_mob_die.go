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
	luax.CallAsync(root, thread, hook, args...).OnError(func(err error) {
		log.Printf("mob script %s: %v", hook, err)
	})
}
