package entity

import (
	"fmt"
	"log"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) OpenNpc(actx actor.Context, npc *Npc) error {
	if actx == nil {
		return fmt.Errorf("actor context is nil")
	}
	npcID := npc.Wz.ID

	if shop := ch.GameWorld.GetResources().GetShop(npcID); shop != nil {
		ch.CurrentShopID = npcID
		ch.Listener.OnOpenNpcShop(ch, npcID, shop)
		return nil
	}

	ch.Dialog.Close()

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
		ch.Listener.OnScriptError(ch, scriptPath, err)
		return err
	}
	luax.SetConfiguration(luaThread, luax.Configuration{
		ActorContext: actx,
		ActorPID:     mapInstance.LogicActorPID(),
	})
	ch.Dialog.Set(luaThread)
	luax.CallAsync(actx, root, luaThread, "on_click", ch, npc).Do(func([]lua.LValue) error {
		ch.Dialog.Reset()
		ch.Listener.OnUnlockAction(ch)
		return nil
	}).OnError(func(err error) {
		log.Printf("npc script on_click npc=%d: %v", npcID, err)
		ch.Listener.OnScriptError(ch, scriptPath, err)
		ch.Dialog.Reset()
		ch.Listener.OnUnlockAction(ch)
	})
	return nil
}
