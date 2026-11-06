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
	a.CharacterId = reader.ReadU32()
	tbyte := reader.ReadU8()
	a.Targets = tbyte >> 4
	a.Hits = tbyte & 0xF
	a.SkillLevel = reader.ReadU8()
	if a.SkillLevel != 0 {
		a.Skill = reader.ReadU32()
	}

	reader.ReadBool()
	a.Unk = reader.ReadU8()
	a.Speed = reader.ReadU8()
	a.Display = reader.ReadU8()
	a.CashBullet = reader.ReadU32()
	for len(a.Damages) < int(a.Targets) && reader.Remaining() >= 5 {
		target := dto.AttackPair{OID: reader.ReadU32()}
		reader.ReadU8()
		for range a.Hits {
			damage := reader.ReadU32()
			target.DamagePairs = append(target.DamagePairs, dto.DamagePair{Damage: damage &^ 0x80000000, Unknown: damage&0x80000000 != 0})
		}
		a.Damages = append(a.Damages, target)
	}

	if reader.Remaining() >= 4 {
		a.Charge = reader.ReadU32()
	}
}
