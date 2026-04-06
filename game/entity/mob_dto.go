package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

// ToDTO converts entity Mob to dto Mob (includes mob buff mask/entries for spawn/control packet).
func (m *Mob) ToDTO() *dto.Mob {
	if m == nil {
		return nil
	}

	mobId := uint32(0)
	if m.Wz != nil {
		mobId = m.Wz.ID
	}

	mask, buffEntries := m.getMobBuffMaskAndEntries()
	statuses := make([]dto.MobBuffEntry, 0, len(buffEntries))
	for _, e := range buffEntries {
		statuses = append(statuses, dto.MobBuffEntry{
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
		MaxHp:      m.GetMaxHp(),
		Mp:         m.Mp,
		MaxMp:      m.GetMaxMp(),
		StatusMask: int32(mask),
		Statuses:   statuses,
	}
}
