package entity

import (
	"slices"

	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
)

func (m *Mob) ToDTO() *dto.Mob {
	if m == nil {
		return nil
	}

	mobId := uint32(0)
	if m.Wz != nil {
		mobId = m.Wz.ID
	}

	var mask int32
	var statuses []dto.MobBuffEntry
	flags := make([]constant.MobBuffFlag, 0, len(m.Buffs.byFlag))
	for f := range m.Buffs.byFlag {
		flags = append(flags, f)
	}
	slices.Sort(flags)

	for _, f := range flags {
		ent := m.Buffs.byFlag[f]
		if ent == nil {
			continue
		}
		v, ok := ent.Values[f]
		if !ok {
			continue
		}
		mask |= int32(f)
		skillID := uint16(0)
		if ent.Wz != nil {
			skillID = uint16(ent.Wz.ID)
		}
		statuses = append(statuses, dto.MobBuffEntry{
			X:          int16(v),
			SkillID:    skillID,
			SkillLevel: uint16(ent.SkillLevel),
		})
	}

	return &dto.Mob{
		OID:         m.OID,
		MobId:       mobId,
		Position:    m.Position,
		Stance:      m.Stance,
		Foothold:    m.Foothold,
		Hp:          m.GetHp(),
		MaxHp:       m.GetMaxHp(),
		Mp:          m.GetMp(),
		MaxMp:       m.GetMaxMp(),
		StatusMask:  mask,
		Statuses:    statuses,
		Reflections: m.Buffs.Reflections(),
	}
}
