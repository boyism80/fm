package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) OpenNpc(actx actor.Context, npcID uint32, caller *lua.LState) error {
	if ch == nil {
		return fmt.Errorf("character is nil")
	}
	if actx == nil {
		return fmt.Errorf("actor context is nil")
	}
	if npcID == 0 {
		return fmt.Errorf("npc id required")
	}

	if old := ch.GetDialog(); old != nil {
		ch.ResetDialog()
		if cfg, ok := luax.GetConfiguration(old); ok {
			cfg.CallPromise = nil
			cfg.KeepAlive = false
			luax.SetConfiguration(old, cfg)
		}
		if old != caller {
			luax.Close(old)
		}
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return fmt.Errorf("lua state not available")
	}
	scriptPath := fmt.Sprintf("script/npc/%d.lua", npcID)
	luaThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return err
	}
	luax.SetConfiguration(luaThread, luax.Configuration{
		ActorContext: actx,
		ActorPID:     mapInstance.LogicActorPID(),
		KeepAlive:    true,
	})
	luax.CallAsync(root, luaThread, "on_click", ch, npcID).Then(func(_ interface{}) (interface{}, error) {
		ch.ResetDialog()
		if ch.Listener != nil {
			ch.Listener.OnUnlockAction(ch)
		}
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("npc script on_click npc=%d: %v", npcID, err)
		ch.ResetDialog()
		if ch.Listener != nil {
			ch.Listener.OnUnlockAction(ch)
		}
	})
	return nil
}
