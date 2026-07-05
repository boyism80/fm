package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
)

func (ch *Character) RunQuestScript(actx actor.Context, questID uint32, npcID uint32, entry string) error {
	if ch == nil {
		return fmt.Errorf("character is nil")
	}
	if ch.GetDialog() != nil {
		return fmt.Errorf("dialog already active")
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return fmt.Errorf("lua state not available")
	}
	scriptPath := fmt.Sprintf("script/quest/%d.lua", questID)
	luaThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return err
	}
	luax.SetConfiguration(luaThread, luax.Configuration{
		ActorContext: actx,
		MapActorPID:  mapInstance.GetActorPID(),
		KeepAlive:    true,
	})
	luax.CallAsync(root, luaThread, entry, ch, npcID).Then(func(_ interface{}) (interface{}, error) {
		ch.ResetDialog()
		if ch.Listener != nil {
			ch.Listener.OnUnlockAction(ch)
		}
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("quest script %s quest=%d npc=%d: %v", entry, questID, npcID, err)
		ch.ResetDialog()
		if ch.Listener != nil {
			ch.Listener.OnUnlockAction(ch)
		}
	})
	return nil
}
