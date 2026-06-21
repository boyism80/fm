package entity

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

const mobTimerRemoveAfterKey = "mob:removeAfter"

func (m *Mob) armRemoveAfter() {
	if m == nil || m.Wz == nil {
		return
	}
	if m.IsFake() || m.SpawnLink != 0 {
		return
	}
	if m.SpawnType == constant.MobSpawnTypeRevive {
		return
	}
	removeAfter := m.Wz.RemoveAfter
	if removeAfter <= 0 {
		return
	}

	delay := time.Duration(removeAfter) * time.Second
	m.RemoveTimer(mobTimerRemoveAfterKey)
	m.AddTimer(mobTimerRemoveAfterKey, delay, false, func() {
		m.expireRemoveAfter()
	})
}

func (m *Mob) expireRemoveAfter() {
	if m == nil || m.GetHp() == 0 {
		return
	}

	m.RemoveTimer(mobTimerRemoveAfterKey)

	mapInstance := m.GetMap()
	if mapInstance == nil {
		return
	}

	m.SetHp(0, false)

	pos := m.Position
	linkOID := m.OID
	if !m.IsFake() && m.Wz != nil && len(m.Wz.Revives) > 0 && m.SpawnLink == 0 {
		m.handleRevives(mapInstance, pos, linkOID, m.Wz.Revives)
	}

	dieAnim := constant.MobDieAnimationTypeFadeOut
	if m.Wz != nil && m.Wz.SelfDestructionAction >= 0 {
		dieAnim = constant.MobDieAnimationType(m.Wz.SelfDestructionAction)
	}
	_ = mapInstance.RemoveMob(m.OID, dieAnim)
}
