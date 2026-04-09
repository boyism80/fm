package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/stream"
)

type RangedAttack struct {
	dto.AttackInfo
	CharacterId uint32
	SkillLevel  uint8
	CashBullet  uint32
}

func (a *RangedAttack) Opcode() uint16 {
	return 0x84
}

func (a *RangedAttack) Serialize(writer *stream.StreamWriter) error {
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
	writer.WriteU32(a.CashBullet)
	for _, oned := range a.Damages {
		if oned.DamagePairs != nil {
			writer.WriteU32(oned.OID)
			writer.WriteU8(0x07)
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

func (a *RangedAttack) Deserialize(reader *stream.StreamReader) {
}
