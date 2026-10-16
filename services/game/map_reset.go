package game

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
)

func (m *Map) Reset() {
	if m == nil {
		return
	}
	m.KillAllMonsters(constant.MobDieAnimationTypeFadeOut)
	_ = m.ReloadReactors()
	m.RemoveAllFieldDrops()
	m.ClearProperties()
	m.ClearBlockedMobGen()
	m.SetReactorGenEnabled(true)
	m.RemoveRuntimeNpcs()
	m.resetPortalScripts()
	m.summonedMobSpawns = nil
	for _, spawn := range m.mobSpawns {
		if spawn == nil {
			continue
		}
		spawn.Spawned = false
		spawn.LastSpawnedAt = time.Time{}
	}
}

func (m *Map) RemoveAllFieldDrops() {
	if m == nil {
		return
	}
	oids := make([]uint32, 0)
	for oid := range m.GetItems() {
		oids = append(oids, oid)
	}
	for _, oid := range oids {
		_ = m.RemoveItem(oid, constant.RemoveItemTypeAnimated, 0)
	}
}
