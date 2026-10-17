package game

import (
	"fmt"
	"log"
	"os"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (d *Dialog) Open(actx actor.Context, npc *Npc) error {
	if actx == nil {
		return fmt.Errorf("actor context is nil")
	}
	d.owner.account.storage.Close()
	if d.Busy() {
		return ErrDialogBusy
	}
	npcID := npc.wz.ID

	if shop := d.owner.gameWorld.GetResources().GetShop(npcID); shop != nil {
		d.OpenShop(npcID)
		d.owner.listener.OnOpenNpcShop(d.owner, npcID, shop)
		return nil
	}

	mapInstance := d.owner.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return fmt.Errorf("lua state not available")
	}
	scriptPath := fmt.Sprintf("script/npc/%d.lua", npcID)
	if _, err := os.Stat(scriptPath); err != nil {
		scriptPath = "script/npc/default.lua"
	}
	luaThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		d.owner.listener.OnScriptError(d.owner, scriptPath, err)
		return err
	}
	luax.SetConfiguration(luaThread, luax.Configuration{
		ActorContext: actx,
		ActorPID:     mapInstance.LogicActorPID(),
	})
	d.Set(luaThread)
	luax.CallAsync(actx, root, luaThread, "on_click", d.owner, npc).Do(func([]lua.LValue) error {
		d.Reset()
		d.owner.listener.OnUnlockAction(d.owner)
		return nil
	}).OnError(func(err error) {
		log.Printf("npc script on_click npc=%d: %v", npcID, err)
		d.owner.listener.OnScriptError(d.owner, scriptPath, err)
		d.Reset()
		d.owner.listener.OnUnlockAction(d.owner)
	})
	return nil
}
