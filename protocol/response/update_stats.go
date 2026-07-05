package response

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
)

type UpdateStats struct {
	Stats        map[constant.Stat]int32
	UnlockAction bool
}

func (p *UpdateStats) Opcode() uint16 {
	return 0x14
}

func (p *UpdateStats) Serialize(writer *stream.StreamWriter) error {
	if p.Stats == nil {
		p.Stats = map[constant.Stat]int32{}
	}

	writer.WriteBoolean(p.UnlockAction)

	mask := uint32(0)
	for k := range p.Stats {
		mask = mask | uint32(k)
	}
	writer.WriteU32(mask)
	order := []constant.Stat{constant.StatSkin, constant.StatFace, constant.StatHair, constant.StatPet, constant.StatLevel, constant.StatClass, constant.StatStr, constant.StatDex, constant.StatInt, constant.StatLuk, constant.StatHP, constant.StatMaxHP, constant.StatMP, constant.StatMaxMP, constant.StatAvailableAP, constant.StatAvailableSP, constant.StatEXP, constant.StatPopulation, constant.StatMeso}
	for _, stat := range order {
		if val, ok := p.Stats[stat]; ok {
			switch stat {
			case constant.StatSkin:
				writer.WriteU16(uint16(val))

			case constant.StatFace, constant.StatHair:
				writer.WriteU32(uint32(val))

			case constant.StatPet:
				writer.WriteU64(uint64(val))

			case constant.StatLevel:
				writer.WriteU8(uint8(val))

			case constant.StatClass, constant.StatStr, constant.StatDex,
				constant.StatInt, constant.StatLuk, constant.StatHP,
				constant.StatMaxHP, constant.StatMP, constant.StatMaxMP,
				constant.StatAvailableAP, constant.StatAvailableSP:
				writer.WriteU16(uint16(val))

			case constant.StatEXP, constant.StatPopulation, constant.StatMeso:
				writer.WriteU32(uint32(val))
			}
		}
	}

	if mask&uint32(constant.StatPet) != 0 {
		writer.WriteU8(1)
	}

	return nil
}

func (p *UpdateStats) Deserialize(reader *stream.StreamReader) {
}
