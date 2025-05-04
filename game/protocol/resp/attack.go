package resp

import (
	"github.com/boyism80/fm/common/stream"
	"github.com/boyism80/fm/game/entity"
)

type Attack struct {
	entity.AttackInfo
	CharacterId uint32
	SkillLevel  uint8
}

func (a *Attack) Serialize(writer *stream.StreamWriter) error {
	writer.WriteU16(0x83)
	writer.WriteU32(a.CharacterId)
	tbyte := (a.Targets << 4) | (a.Hits & 0xF)
	writer.WriteU8(tbyte)
	if a.Skill == 0 {
		writer.WriteU8(0)
	} else {
		writer.WriteU8(a.SkillLevel)
		writer.WriteU32(a.Skill)
	}

	writer.WriteBoolean(false) // 버프상태
	writer.WriteU8(a.Unk)
	writer.WriteU8(a.Speed)
	writer.WriteU8(a.Display)
	writer.WriteU32(0) // cash bullet?
	for _, oned := range a.AllDamage {
		if oned.Attack != nil {
			writer.WriteU32(oned.ObjectId)
			writer.WriteU8(0x07)
			if a.Skill == 4211006 {
				writer.WriteU8(uint8(len(oned.Attack)))
			}

			for _, v := range oned.Attack {
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
