package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) OpenNpc(actx actor.Context, npc *Npc, caller *lua.LState) error {
	if actx == nil {
		return fmt.Errorf("actor context is nil")
	}
	npcID := npc.Wz.ID

	if shop := ch.GameWorld.GetResources().GetShop(npcID); shop != nil {
		ch.CurrentShopID = npcID
		ch.Listener.OnOpenNpcShop(ch, npcID, shop)
		return nil
	}

	if old := ch.GetDialog(); old != nil {
		ch.ResetDialog()
		if cfg, ok := luax.GetConfiguration(old); ok {
			cfg.CallPromise = nil
			luax.SetConfiguration(old, cfg)
		}
		// The caller is still running and closes itself when it finishes.
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
	})
	luax.CallAsync(root, luaThread, "on_click", ch, npc).Then(func(_ interface{}) (interface{}, error) {
		ch.ResetDialog()
		ch.Listener.OnUnlockAction(ch)
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("npc script on_click npc=%d: %v", npcID, err)
		ch.ResetDialog()
		ch.Listener.OnUnlockAction(ch)
	})
	return nil
}
