package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

// ToDTO converts entity Mob to dto Mob
func (m *Mob) ToDTO() *dto.Mob {
	if m == nil {
		return nil
	}

	mobId := uint32(0)
	if m.Wz != nil {
		mobId = m.Wz.ID
	}

	return &dto.Mob{
		OID:      m.OID,
		MobId:    mobId,
		Position: m.Position,
		Stance:   m.Stance,
		Foothold: m.Foothold,
		Hp:       m.Hp,
		MaxHp:    m.Life.GetMaxHp(),
		Mp:       m.Mp,
		MaxMp:    m.Life.GetMaxMp(),
	}
}
