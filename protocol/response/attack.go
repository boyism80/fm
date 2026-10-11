package response

import (
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
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

	writer.WriteBoolean(false)
	writer.WriteU8(a.Unk)
	writer.WriteU8(a.Speed)
	writer.WriteU8(a.Display)
	writer.WriteU32(0)
	for _, oned := range a.Damages {
		if oned.DamagePairs != nil {
			writer.WriteU32(oned.OID)
			writer.WriteU8(0x07)
			if a.Skill == uint32(constant.SkillMesoExplosion) {
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

	switch constant.SkillID(a.Skill) {
	case constant.SkillBigBang, constant.SkillBigBang2221001, constant.SkillBigBang2321001:
		writer.WriteU32(a.Charge)
	}
	return nil
}

func (a *Attack) Deserialize(reader *stream.StreamReader) {
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
	reader.ReadU32()
	// Serialize skips targets without damage pairs, so Targets can count more than were written.
	for len(a.Damages) < int(a.Targets) && reader.Remaining() >= 5 {
		target := dto.AttackPair{OID: reader.ReadU32()}
		reader.ReadU8()
		hits := int(a.Hits)
		if a.Skill == uint32(constant.SkillMesoExplosion) {
			hits = int(reader.ReadU8())
		}
		for range hits {
			damage := reader.ReadU32()
			target.DamagePairs = append(target.DamagePairs, dto.DamagePair{Damage: damage &^ 0x80000000, Unknown: damage&0x80000000 != 0})
		}
		a.Damages = append(a.Damages, target)
	}

	switch constant.SkillID(a.Skill) {
	case constant.SkillBigBang, constant.SkillBigBang2221001, constant.SkillBigBang2321001:
		a.Charge = reader.ReadU32()
	}
}
