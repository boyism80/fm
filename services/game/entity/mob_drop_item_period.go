package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

const (
	mobTimerDropItemPeriodKey = "mob:dropItemPeriod"
	mobDropItemHitCooldown    = 11 * time.Second
	mobAutoDropRiceCakeExpire = 6 * time.Second
)

func (m *Mob) MarkHit() {
	if m == nil {
		return
	}
	m.lastHitAt = clock.Now()
}

func (m *Mob) startTimedDrop() {
	if m == nil || m.Wz == nil || m.IsFake() {
		return
	}
	period := m.Wz.DropItemPeriod
	if period <= 0 {
		return
	}
	m.dropItemCount = 0
	m.RemoveTimer(mobTimerDropItemPeriodKey)
	interval := time.Duration(period) * time.Second
	m.AddTimer(mobTimerDropItemPeriodKey, interval, true, func() {
		m.dropTimedItem()
	})
}

func (m *Mob) dropTimedItem() {
	if m == nil || m.Wz == nil || m.GetHp() == 0 {
		return
	}
	mapInstance := m.GetMap()
	if mapInstance == nil || m.GameWorld == nil {
		return
	}
	if !m.lastHitAt.IsZero() && clock.Now().Before(m.lastHitAt.Add(mobDropItemHitCooldown)) {
		return
	}
	nextCount := m.dropItemCount + 1
	root := mapInstance.EnsureLuaRoot(nil)
	if root == nil {
		return
	}
	mobID := m.Wz.ID
	thread, err := luax.NewThread(root, fmt.Sprintf("script/mob/%d.lua", mobID))
	if err != nil {
		return
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	result, err := luax.Call(thread, "on_mob_period_drop", m, nextCount)
	if err != nil || result == nil || result == lua.LNil {
		return
	}
	if result.Type() != lua.LTNumber {
		return
	}
	itemID := uint32(lua.LVAsNumber(result))
	if itemID == 0 {
		return
	}
	item, err := NewItem(itemID, 1, m.GameWorld)
	if err != nil {
		return
	}
	pos := m.GetPosition()
	if err := mapInstance.SpawnItem(item, ItemSpawn{
		Position: pos,
		From:     pos,
		DropType: constant.DropTypeFFA,
	}); err != nil {
		return
	}
	if itemID == 4001101 {
		item.GetFieldPlacement().RegisterExpire(mobAutoDropRiceCakeExpire)
	}
	m.dropItemCount = nextCount
}
