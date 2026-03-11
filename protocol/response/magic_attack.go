package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type MagicAttack struct {
	dto.AttackInfo
	CharacterId uint32
	SkillLevel  uint8
}

func (a *MagicAttack) Opcode() uint16 {
	return 0x85
}

func (a *MagicAttack) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU32(a.CharacterId)
	tbyte := (a.Targets << 4) | (a.Hits & 0xF)
	writer.WriteU8(tbyte)
	if a.Skill == 0 {
		writer.WriteU8(0)
	} else {
		writer.WriteU8(a.SkillLevel)
		writer.WriteU32(a.Skill)
	}

	writer.WriteBoolean(false)
	writer.WriteU8(a.Unk)
	writer.WriteU8(a.Speed)
	writer.WriteU8(a.Display)
	writer.WriteU32(0)

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
	}

	return nil
}

func (a *MagicAttack) Deserialize(reader *stream.StreamReader) error {
	return nil
}
