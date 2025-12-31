package response

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/protocol/dto"
)

type Attack struct {
	dto.AttackInfo
	CharacterId uint32
	SkillLevel  uint8
}

func (a *Attack) Opcode() uint16 {
	return 0x83
}

func (a *Attack) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(a.CharacterId)
	tbyte := (a.Targets << 4) | (a.Hits & 0xF)
	writer.WriteU8(tbyte)
	if a.Skill == 0 {
		writer.WriteU8(0)
	} else {
		writer.WriteU8(a.SkillLevel)
		writer.WriteU32(a.Skill)
	}

	writer.WriteBoolean(false) // 버프?�태
	writer.WriteU8(a.Unk)
	writer.WriteU8(a.Speed)
	writer.WriteU8(a.Display)
	writer.WriteU32(0) // cash bullet?
	for _, oned := range a.Damages {
		if oned.DamagePairs != nil {
			writer.WriteU32(oned.OID)
			writer.WriteU8(0x07)
			if a.Skill == 4211006 {
				writer.WriteU8(uint8(len(oned.DamagePairs)))
			}

			for _, v := range oned.DamagePairs {
				if v.Unknown {
					writer.WriteU32(v.Damage | 0x80000000)
				} else {
					writer.WriteU32(v.Damage)
				}
			}
		}
	}

	if a.Charge > 0 {
		writer.WriteU32(a.Charge)
	} else {
		writer.WriteU32(0)
	}
	return nil
}

func (a *Attack) Deserialize(reader *stream.StreamReader) error {
	return nil
}
