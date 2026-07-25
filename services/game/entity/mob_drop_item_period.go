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

func (m *Mob) armDropItemPeriod() {
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
		m.doDropItemPeriod()
	})
}

func (m *Mob) resolvePeriodDropItemID(nextCount int) (uint32, bool) {
	if m == nil || m.Wz == nil {
		return 0, false
	}
	mapInstance := m.GetMap()
	if mapInstance == nil {
		return 0, false
	}
	root := mapInstance.EnsureLuaRoot(nil)
	if root == nil {
		return 0, false
	}
	mobID := m.Wz.ID
	hook := "on_mob_period_drop"
	thread, err := luax.NewThread(root, fmt.Sprintf("script/mob/%d.lua", mobID))
	if err != nil {
		return 0, false
	}
	luax.SetConfiguration(thread, luax.Configuration{
		ActorPID: mapInstance.LogicActorPID(),
	})
	result, err := luax.Call(thread, hook, m, nextCount)
	if err != nil || result == nil || result == lua.LNil {
		return 0, false
	}
	if result.Type() != lua.LTNumber {
		return 0, false
	}
	itemID := uint32(lua.LVAsNumber(result))
	if itemID == 0 {
		return 0, false
	}
	return itemID, true
}

func (m *Mob) doDropItemPeriod() {
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
	itemID, ok := m.resolvePeriodDropItemID(nextCount)
	if !ok {
		return
	}
	item, err := NewItem(itemID, 1, m.GameWorld)
	if err != nil {
		return
	}
	pos := m.GetPosition()
	fp := &FieldPlacement{
		ObjectCore:   &ObjectCore{Position: pos},
		SpawnedPoint: pos,
		DropType:     constant.DropTypeFFA,
	}
	fp.ObjectCore.self = fp
	item.BindFieldPlacement(fp)
	if err := mapInstance.SpawnItem(item, 0, constant.DropTypeFFA); err != nil {
		return
	}
	if itemID == 4001101 {
		if placed := item.GetFieldPlacement(); placed != nil {
			placed.RegisterExpire(mobAutoDropRiceCakeExpire)
		}
	}
	m.dropItemCount = nextCount
}
