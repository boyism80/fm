package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

// ToDTO converts entity Mob to dto Mob (includes status/debuff for spawn/control packet).
func (m *Mob) ToDTO() *dto.Mob {
	if m == nil {
		return nil
	}

	mobId := uint32(0)
	if m.Wz != nil {
		mobId = m.Wz.ID
	}

	mask, debuffEntries := m.getDebuffMaskAndEntries()
	statuses := make([]dto.MobStatusEntry, 0, len(debuffEntries))
	for _, e := range debuffEntries {
		statuses = append(statuses, dto.MobStatusEntry{
			X:       int16(e.Value),
			SkillID: e.SkillID,
		})
	}

	return &dto.Mob{
		OID:        m.OID,
		MobId:      mobId,
		Position:   m.Position,
		Stance:     m.Stance,
		Foothold:   m.Foothold,
		Hp:         m.Hp,
		MaxHp:      m.Life.GetMaxHp(),
		Mp:         m.Mp,
		MaxMp:      m.Life.GetMaxMp(),
		StatusMask: int32(mask),
		Statuses:   statuses,
	}
}
