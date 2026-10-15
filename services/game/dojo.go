package game

import (
	"github.com/boyism80/fm/maple/wz"
)

const DojoEnergyFull = 300

type Dojo struct {
	owner  *Character
	Energy int
}

func (d *Dojo) onField() bool {
	m := d.owner.GetMap()
	return m != nil && m.Wz.FieldType == wz.FieldTypeDojo
}

func (d *Dojo) SetEnergy(energy int) {
	energy = min(max(energy, 0), DojoEnergyFull)
	if d.Energy == energy {
		return
	}

	d.Energy = energy
	d.owner.Listener.OnDojoEnergy(d.owner, energy)
}

func (d *Dojo) AddEnergy(n int) {
	if !d.onField() {
		return
	}
	d.SetEnergy(d.Energy + n)
}

func (d *Dojo) charge(mob *Mob, damage uint64) {
	if !d.onField() || mob.Wz.Boss {
		return
	}
	maxHp := uint64(mob.GetMaxHp())
	if maxHp == 0 {
		return
	}
	d.SetEnergy(d.Energy + max(1, int(min(damage, maxHp)*50/maxHp)))
}
