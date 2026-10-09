package entity

import (
	"github.com/boyism80/fm/services/game/wz"
)

const DojoEnergyFull = 300

func (ch *Character) OnDojoField() bool {
	m := ch.GetMap()
	return m != nil && m.Wz.FieldType == wz.FieldTypeDojo
}

func (ch *Character) DojoEnergy() int {
	return ch.dojoEnergy
}

func (ch *Character) SetDojoEnergy(energy int) {
	energy = min(max(energy, 0), DojoEnergyFull)
	if ch.dojoEnergy == energy {
		return
	}

	ch.dojoEnergy = energy
	ch.Listener.OnDojoEnergy(ch, energy)
}

func (ch *Character) AddDojoEnergy(n int) {
	if !ch.OnDojoField() {
		return
	}
	ch.SetDojoEnergy(ch.dojoEnergy + n)
}

func (ch *Character) chargeDojoEnergy(mob *Mob, damage uint64) {
	if !ch.OnDojoField() || mob.Wz.Boss {
		return
	}
	maxHp := uint64(mob.GetMaxHp())
	if maxHp == 0 {
		return
	}
	ch.SetDojoEnergy(ch.dojoEnergy + max(1, int(min(damage, maxHp)*50/maxHp)))
}
