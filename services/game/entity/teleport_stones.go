package entity

import (
	"sync/atomic"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

const teleportStoneTargetTimeout = 3 * time.Second

type TeleportStones struct {
	owner   *Character
	regular constant.TeleportStoneSlots
	vip     constant.TeleportStoneSlots
}

func (t *TeleportStones) Slots(vip bool) constant.TeleportStoneSlots {
	if vip {
		return t.vip
	}
	return t.regular
}

func (t *TeleportStones) Register(vip bool) {
	m := t.owner.GetMap()
	if m == nil {
		return
	}
	if m.BlocksTeleportStone() || t.query(m, "can_register_teleport_stone", m.TemplateID()) == false {
		t.owner.Listener.OnTeleportStoneFailed(t.owner, vip, pconst.TeleportStoneResultCannotRegister)
		return
	}

	stones := t.Slots(vip)
	if stones.Contains(m.TemplateID()) {
		t.owner.Listener.OnTeleportStones(t.owner, vip)
		return
	}
	for i, mapID := range stones {
		if mapID == constant.TeleportStoneEmpty {
			stones[i] = m.TemplateID()
			break
		}
	}
	t.owner.Listener.OnTeleportStones(t.owner, vip)
}

func (t *TeleportStones) Remove(vip bool, mapID uint32) {
	stones := t.Slots(vip)
	for i, registered := range stones {
		if registered == mapID {
			stones[i] = constant.TeleportStoneEmpty
			t.owner.Listener.OnTeleportStones(t.owner, vip)
			return
		}
	}
}

func (t *TeleportStones) Reset() {
	t.regular = constant.NewTeleportStoneSlots(nil, constant.TeleportStoneCount)
	t.vip = constant.NewTeleportStoneSlots(nil, constant.VipTeleportStoneCount)
	t.owner.Listener.OnTeleportStones(t.owner, false)
	t.owner.Listener.OnTeleportStones(t.owner, true)
}

func (t *TeleportStones) Use(ctx actor.Context, invType constant.InventoryType, slot int16, itemID uint32, mapID uint32, name string) {
	if t.owner.session.cashItemInUse.CompareAndSwap(false, true) == false {
		t.owner.Listener.OnUpdateStats(t.owner, nil, true)
		return
	}
	finish := func() {
		t.owner.session.cashItemInUse.Store(false)
		t.owner.Listener.OnUpdateStats(t.owner, nil, true)
	}
	vip := constant.IsVipTeleportStone(itemID)
	fail := func(result pconst.TeleportStoneResult) {
		t.owner.Listener.OnTeleportStoneFailed(t.owner, vip, result)
		finish()
	}

	switch invType {
	case constant.InventoryTypeCash:
		if constant.IsCashTeleportStone(itemID) == false {
			finish()
			return
		}
	case constant.InventoryTypeConsume:
		if constant.GetConsumeType(itemID) != constant.ConsumeTypeTeleportStone {
			finish()
			return
		}
	default:
		finish()
		return
	}
	item := t.owner.Inventory.Tabs[invType].Get(uint8(slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		finish()
		return
	}
	if t.owner.GetHp() <= 0 || t.owner.Dialog.ShopID != 0 || t.owner.Dialog.Thread() != nil {
		finish()
		return
	}

	consume := func(actor.Context) {
		item := t.owner.Inventory.Tabs[invType].Get(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			t.owner.Inventory.RemoveItem(invType, slot, 1)
		}
		finish()
	}

	if name == "" {
		if mapID == constant.TeleportStoneEmpty || t.Slots(vip).Contains(mapID) == false {
			fail(pconst.TeleportStoneResultCannotGo)
			return
		}
		t.teleport(ctx, itemID, vip, t.owner.GameWorld.GetMapSystem().Get(mapID), 0, false, consume, fail)
		return
	}

	targetID, found := t.owner.GameWorld.GetDispatchSystem().FindCharacterID(name)
	if found == false || targetID == t.owner.GetID() {
		fail(pconst.TeleportStoneResultNotFound)
		return
	}
	admin := t.owner.HasRoleAtLeast(constant.RoleAdmin)
	var settled atomic.Bool
	settle := func(run func(ctx actor.Context)) {
		t.owner.GameWorld.GetDispatchSystem().Call(t.owner.GetID(), func(ctx actor.Context) {
			if settled.CompareAndSwap(false, true) {
				run(ctx)
			}
		})
	}
	time.AfterFunc(teleportStoneTargetTimeout, func() {
		settle(func(actor.Context) {
			fail(pconst.TeleportStoneResultNotFound)
		})
	})
	t.owner.GameWorld.GetDispatchSystem().CallCharacter(targetID, func(_ actor.Context, target *Character) {
		if target == nil || ((target.IsHidden() || target.HasRoleAtLeast(constant.RoleAdmin)) && admin == false) {
			settle(func(actor.Context) {
				fail(pconst.TeleportStoneResultNotFound)
			})
			return
		}
		targetMap := target.GetMap()
		spawnPoint := targetMap.Wz.FindClosestPortalSpawnID(target.Position)
		settle(func(ctx actor.Context) {
			t.teleport(ctx, itemID, vip, targetMap, spawnPoint, true, consume, fail)
		})
	})
}

func (t *TeleportStones) teleport(ctx actor.Context, itemID uint32, vip bool, target *Map, spawnPoint uint8, byName bool, onEnter func(actor.Context), onFail func(pconst.TeleportStoneResult)) {
	m := t.owner.GetMap()
	if m == nil || target == nil {
		onFail(pconst.TeleportStoneResultCannotGo)
		return
	}
	if m.TemplateID() == target.TemplateID() {
		onFail(pconst.TeleportStoneResultCurrentMap)
		return
	}
	if m.BlocksTeleportStone() || target.BlocksTeleportStone() {
		onFail(pconst.TeleportStoneResultCannotGo)
		return
	}
	if vip == false && m.TemplateID()/100000000 != target.TemplateID()/100000000 {
		onFail(pconst.TeleportStoneResultCannotGo)
		return
	}
	if t.query(m, "can_teleport_stone", itemID, m.TemplateID(), target.TemplateID(), byName) == false {
		onFail(pconst.TeleportStoneResultCannotGo)
		return
	}

	err := t.owner.GameWorld.GetMapSystem().Warp(ctx, t.owner, target, spawnPoint, onEnter)
	if err != nil {
		onFail(pconst.TeleportStoneResultCannotGo)
	}
}

func (t *TeleportStones) query(m *Map, hook string, args ...interface{}) bool {
	thread, err := luax.NewThread(m.GetLuaRoot(), constant.CharacterQueryScriptPath)
	if err != nil {
		return false
	}
	allowed, err := luax.Call(thread, hook, append([]interface{}{t.owner}, args...)...)
	return err == nil && allowed == lua.LTrue
}
