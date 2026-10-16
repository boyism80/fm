package game

import (
	"github.com/boyism80/fm/maple/wz"
)

const DojoEnergyFull = 300

type Dojo struct {
	owner  *Character
	energy int
}

func (d *Dojo) onField() bool {
	m := d.owner.GetMap()
	return m != nil && m.wz.FieldType == wz.FieldTypeDojo
}

func (d *Dojo) SetEnergy(energy int) {
	energy = min(max(energy, 0), DojoEnergyFull)
	if d.energy == energy {
		return
	}

	d.energy = energy
	d.owner.listener.OnDojoEnergy(d.owner, energy)
}

func (d *Dojo) AddEnergy(n int) {
	if !d.onField() {
		return
	}
	d.SetEnergy(d.energy + n)
}

func (d *Dojo) charge(mob *Mob, damage uint64) {
	if !d.onField() || mob.wz.Boss {
		return
	}
	maxHp := uint64(mob.GetMaxHp())
	if maxHp == 0 {
		return
	}
	d.SetEnergy(d.energy + max(1, int(min(damage, maxHp)*50/maxHp)))
}
