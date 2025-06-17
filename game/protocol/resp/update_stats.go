package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/constant"
)

type UpdateStats struct {
	Stats        map[constant.Stat]int32
	UnlockAction bool
}

func (p *UpdateStats) Serialize(writer *stream.StreamWriter) error {
	if p.Stats == nil {
		p.Stats = map[constant.Stat]int32{}
	}

	writer.WriteU16(0x14)
	writer.WriteBoolean(p.UnlockAction)

	mask := uint32(0)
	for k := range p.Stats {
		mask = mask | uint32(k)
	}
	writer.WriteU32(mask)
	order := []constant.Stat{constant.STAT_SKIN, constant.STAT_FACE, constant.STAT_HAIR, constant.STAT_PET, constant.STAT_LEVEL, constant.STAT_JOB, constant.STAT_STR, constant.STAT_DEX, constant.STAT_INT, constant.STAT_LUK, constant.STAT_HP, constant.STAT_MAX_HP, constant.STAT_MP, constant.STAT_MAX_MP, constant.STAT_AVAILABLE_AP, constant.STAT_AVAILABLE_SP, constant.STAT_EXP, constant.STAT_FAME, constant.STAT_MESO}
	for _, stat := range order {
		if val, ok := p.Stats[stat]; ok {
			switch stat {
			case constant.STAT_SKIN:
				writer.WriteU16(uint16(val))

			case constant.STAT_FACE, constant.STAT_HAIR:
				writer.WriteU32(uint32(val))

			case constant.STAT_PET:
				writer.WriteU64(uint64(val))

			case constant.STAT_LEVEL:
				writer.WriteU8(uint8(val))

			case constant.STAT_JOB, constant.STAT_STR, constant.STAT_DEX,
				constant.STAT_INT, constant.STAT_LUK, constant.STAT_HP,
				constant.STAT_MAX_HP, constant.STAT_MP, constant.STAT_MAX_MP,
				constant.STAT_AVAILABLE_AP, constant.STAT_AVAILABLE_SP:
				writer.WriteU16(uint16(val))

			case constant.STAT_EXP, constant.STAT_FAME, constant.STAT_MESO:
				writer.WriteU32(uint32(val))
			}
		}
	}

	if mask&uint32(constant.STAT_PET) != 0 {
		writer.WriteU8(1)
	}

	return nil
}

func (p *UpdateStats) Deserialize(reader *stream.StreamReader) error {
	return nil
}
