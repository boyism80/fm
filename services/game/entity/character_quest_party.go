package entity

import (
	"fmt"
)

func (ch *Character) RunQuestHook(questID uint32, hook string) {
	if ch == nil || questID == 0 || hook == "" {
		return
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return
	}
	scriptPath := fmt.Sprintf("script/quest/%d.lua", questID)
	mapInstance.runMobLuaHook(root, scriptPath, hook, ch)
}
