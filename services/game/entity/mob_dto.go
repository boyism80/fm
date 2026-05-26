package entity

import (
	"github.com/boyism80/fm/protocol/dto"
)

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
			X:          int16(e.Value),
			SkillID:    uint16(e.SkillID),
			SkillLevel: 0,
		})
	}

	return &dto.Mob{
		OID:        m.OID,
		MobId:      mobId,
		Position:   m.Position,
		Stance:     m.Stance,
		Foothold:   m.Foothold,
		Hp:         m.GetHp(),
		MaxHp:      m.GetMaxHp(),
		Mp:         m.GetMp(),
		MaxMp:      m.GetMaxMp(),
		StatusMask: int32(mask),
		Statuses:   statuses,
	}
}
