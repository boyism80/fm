package entity

import (
	"fmt"
	"log"

	"github.com/boyism80/fm/core/luax"
)

func (ch *Character) RunQuestHook(questID uint32, hook string, args ...interface{}) bool {
	if ch == nil || questID == 0 || hook == "" {
		return false
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return false
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return false
	}
	scriptPath := fmt.Sprintf("script/quest/%d.lua", questID)
	thread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return false
	}
	if !luax.HasFunc(thread, hook) {
		luax.Close(thread)
		return false
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	luax.CallAsync(nil, root, thread, hook, append([]interface{}{ch}, args...)...).OnError(func(err error) {
		log.Printf("quest script %s %s: %v", scriptPath, hook, err)
	})
	return true
}
