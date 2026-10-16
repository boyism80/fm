package game

import (
	"slices"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/protocol/dto"
)

func (m *Mob) BuildS2C() *dto.Mob {
	if m == nil {
		return nil
	}

	mobId := uint32(0)
	if m.wz != nil {
		mobId = m.wz.ID
	}

	var mask int32
	var statuses []dto.MobBuffEntry
	flags := make([]constant.MobBuffFlag, 0, len(m.buffs.byFlag))
	for f := range m.buffs.byFlag {
		flags = append(flags, f)
	}
	slices.Sort(flags)

	for _, f := range flags {
		ent := m.buffs.byFlag[f]
		if ent == nil {
			continue
		}
		v, ok := ent.values[f]
		if !ok {
			continue
		}
		mask |= int32(f)
		skillID := uint16(0)
		if ent.wz != nil {
			skillID = uint16(ent.wz.ID)
		}
		statuses = append(statuses, dto.MobBuffEntry{
			X:          int16(v),
			SkillID:    skillID,
			SkillLevel: uint16(ent.skillLevel),
		})
	}

	return &dto.Mob{
		OID:          m.oid,
		MobId:        mobId,
		Position:     m.position,
		Stance:       m.stance,
		Foothold:     m.foothold,
		Hp:           m.GetHp(),
		MaxHp:        m.GetMaxHp(),
		Mp:           m.GetMp(),
		MaxMp:        m.GetMaxMp(),
		StatusMask:   mask,
		Statuses:     statuses,
		Reflections:  m.buffs.Reflections(),
		CarnivalTeam: m.carnivalTeam,
	}
}
