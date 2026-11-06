package game

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
)

// HurtByEnvironment takes the map's WZ decHP from every living character once per decInterval, the first time one
// interval after someone is there.
func (m *Map) HurtByEnvironment(now time.Time) {
	if m.wz == nil || m.wz.DecHP <= 0 {
		return
	}
	if m.GetPlayerCount() == 0 {
		m.nextEnvironmentHurt = time.Time{}
		return
	}
	if m.nextEnvironmentHurt.IsZero() {
		m.nextEnvironmentHurt = now.Add(m.wz.DecInterval)
		return
	}
	if now.Before(m.nextEnvironmentHurt) {
		return
	}
	m.nextEnvironmentHurt = now.Add(m.wz.DecInterval)
	for _, obj := range m.GetAllPlayers() {
		if ch, ok := obj.(*Character); ok && ch.IsAlive() {
			ch.hurtByEnvironment(m.wz)
		}
	}
}

// hurtByEnvironment spares a character wearing the map's protect item. An HP loss guard potion takes its thaw off the
// damage, a negative thaw in water and a positive one out of it.
func (ch *Character) hurtByEnvironment(field *wz.Map) {
	if field.ProtectItem != 0 && ch.inventory.wears(field.ProtectItem) {
		return
	}
	damage := field.DecHP
	if _, thaw, ok := ch.buffs.GetBuffValue(constant.BuffFlagHpLossGuard); ok && (thaw < 0) == field.Swim {
		damage -= int(max(thaw, -thaw))
	}
	if damage > 0 {
		ch.AddHp(-damage)
	}
}
