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

func (ch *Character) TeleportStones(vip bool) constant.TeleportStones {
	if vip {
		return ch.vipTeleportStones
	}
	return ch.teleportStones
}

func (ch *Character) RegisterTeleportStone(vip bool) {
	m := ch.GetMap()
	if m == nil {
		return
	}
	if m.BlocksTeleportStone() || ch.queryTeleportStone(m, "can_register_teleport_stone", m.TemplateID()) == false {
		ch.Listener.OnTeleportStoneFailed(ch, vip, pconst.TeleportStoneResultCannotRegister)
		return
	}

	stones := ch.TeleportStones(vip)
	if stones.Contains(m.TemplateID()) {
		ch.Listener.OnTeleportStones(ch, vip)
		return
	}
	for i, mapID := range stones {
		if mapID == constant.TeleportStoneEmpty {
			stones[i] = m.TemplateID()
			break
		}
	}
	ch.Listener.OnTeleportStones(ch, vip)
}

func (ch *Character) RemoveTeleportStone(vip bool, mapID uint32) {
	stones := ch.TeleportStones(vip)
	for i, registered := range stones {
		if registered == mapID {
			stones[i] = constant.TeleportStoneEmpty
			ch.Listener.OnTeleportStones(ch, vip)
			return
		}
	}
}

func (ch *Character) ResetTeleportStones() {
	ch.teleportStones = constant.NewTeleportStones(nil, constant.TeleportStoneCount)
	ch.vipTeleportStones = constant.NewTeleportStones(nil, constant.VipTeleportStoneCount)
	ch.Listener.OnTeleportStones(ch, false)
	ch.Listener.OnTeleportStones(ch, true)
}

func (ch *Character) UseTeleportStone(ctx actor.Context, invType constant.InventoryType, slot int16, itemID uint32, mapID uint32, name string) {
	if ch.cashItemInUse.CompareAndSwap(false, true) == false {
		ch.Listener.OnUpdateStats(ch, nil, true)
		return
	}
	finish := func() {
		ch.cashItemInUse.Store(false)
		ch.Listener.OnUpdateStats(ch, nil, true)
	}
	vip := constant.IsVipTeleportStone(itemID)
	fail := func(result pconst.TeleportStoneResult) {
		ch.Listener.OnTeleportStoneFailed(ch, vip, result)
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
	item := ch.Inventory.Tabs[invType].Get(uint8(slot))
	if item == nil || item.GetCount() < 1 || item.GetModel().GetID() != itemID {
		finish()
		return
	}
	if ch.GetHp() <= 0 || ch.CurrentShopID != 0 || ch.Dialog.Thread() != nil {
		finish()
		return
	}

	consume := func(actor.Context) {
		item := ch.Inventory.Tabs[invType].Get(uint8(slot))
		if item != nil && item.GetModel().GetID() == itemID {
			ch.Inventory.RemoveItem(invType, slot, 1)
		}
		finish()
	}

	if name == "" {
		if mapID == constant.TeleportStoneEmpty || ch.TeleportStones(vip).Contains(mapID) == false {
			fail(pconst.TeleportStoneResultCannotGo)
			return
		}
		ch.teleportByStone(ctx, itemID, vip, ch.GameWorld.GetMapSystem().Get(mapID), 0, false, consume, fail)
		return
	}

	targetID, found := ch.GameWorld.GetDispatchSystem().FindCharacterID(name)
	if found == false || targetID == ch.GetID() {
		fail(pconst.TeleportStoneResultNotFound)
		return
	}
	admin := ch.HasRoleAtLeast(constant.RoleAdmin)
	var settled atomic.Bool
	settle := func(run func(ctx actor.Context)) {
		ch.GameWorld.GetDispatchSystem().Call(ch.GetID(), func(ctx actor.Context) {
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
	ch.GameWorld.GetDispatchSystem().CallCharacter(targetID, func(_ actor.Context, target *Character) {
		if target == nil || ((target.IsHidden() || target.HasRoleAtLeast(constant.RoleAdmin)) && admin == false) {
			settle(func(actor.Context) {
				fail(pconst.TeleportStoneResultNotFound)
			})
			return
		}
		targetMap := target.GetMap()
		spawnPoint := targetMap.Wz.FindClosestPortalSpawnID(target.Position)
		settle(func(ctx actor.Context) {
			ch.teleportByStone(ctx, itemID, vip, targetMap, spawnPoint, true, consume, fail)
		})
	})
}

func (ch *Character) teleportByStone(ctx actor.Context, itemID uint32, vip bool, target *Map, spawnPoint uint8, byName bool, onEnter func(actor.Context), onFail func(pconst.TeleportStoneResult)) {
	m := ch.GetMap()
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
	if ch.queryTeleportStone(m, "can_teleport_stone", itemID, m.TemplateID(), target.TemplateID(), byName) == false {
		onFail(pconst.TeleportStoneResultCannotGo)
		return
	}

	err := ch.GameWorld.GetMapSystem().Warp(ctx, ch, target, spawnPoint, onEnter)
	if err != nil {
		onFail(pconst.TeleportStoneResultCannotGo)
	}
}

func (ch *Character) queryTeleportStone(m *Map, hook string, args ...interface{}) bool {
	thread, err := luax.NewThread(m.GetLuaRoot(), constant.CharacterQueryScriptPath)
	if err != nil {
		return false
	}
	allowed, err := luax.Call(thread, hook, append([]interface{}{ch}, args...)...)
	return err == nil && allowed == lua.LTrue
}
