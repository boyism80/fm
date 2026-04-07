package dto

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type DamagePair struct {
	Damage  uint32
	Unknown bool
}

type AttackPair struct {
	OID         uint32
	DamagePairs []DamagePair
}

type AttackInfo struct {
	Targets  uint8
	Hits     uint8
	Skill    uint32
	Charge   uint32
	Unk      uint8
	Speed    uint8
	Display  uint8
	LastTick uint32
	Slot     uint16
	CsStar   uint16
	AOE      uint8
	Damages  []AttackPair
	Position types.Vector2[int16]
}

func (a *AttackInfo) Deserialize(sr *stream.StreamReader, opcode uint16) {
	sr.Skip(1)
	tbyte := sr.ReadU8()
	a.Targets = (tbyte >> 4) & 0xF
	a.Hits = tbyte & 0xF

	a.Skill = sr.ReadU32()

	switch constant.SkillID(a.Skill) {
	case constant.SkillBigBang,
		constant.SkillBigBang2221001,
		constant.SkillBigBang2321001,
		constant.SkillPiercing,
		constant.SkillStormArrow,
		constant.SkillStormArrowCygnus,
		constant.SkillCorkscrewBlow,
		constant.SkillCorkscrewBlowCygnus,
		constant.SkillRapidFire,
		constant.SkillGrenade:
		a.Charge = sr.ReadU32()
	default:
		a.Charge = 0
	}

	sr.Skip(1)

	a.Unk = sr.ReadU8()
	a.Speed = sr.ReadU8()
	a.Display = sr.ReadU8()
	a.LastTick = sr.ReadU32()

	if opcode == 0x1C {
		a.Slot = sr.ReadU16()
		a.CsStar = sr.ReadU16()
		a.AOE = sr.ReadU8()
	}

	damages := make([]AttackPair, 0, a.Targets)
	for range int(a.Targets) {
		oid := sr.ReadU32()

		sr.Skip(14)

		damagePairs := make([]DamagePair, 0, a.Hits)
		for range int(a.Hits) {
			damage := sr.ReadU32()

			damagePairs = append(damagePairs, DamagePair{
				Damage:  damage,
				Unknown: false,
			})
		}

		damages = append(damages, AttackPair{
			OID:         oid,
			DamagePairs: damagePairs,
		})
	}

	if sr.Remaining() >= 4 {
		a.Position = types.Vector2[int16]{X: sr.Read16(), Y: sr.Read16()}
	}

	a.Damages = damages
}

func (a *AttackInfo) Serialize(sw *stream.StreamWriter) error {
	sw.WriteU8(0)
	sw.WriteU8((a.Targets << 4) | a.Hits)
	sw.WriteU32(a.Skill)

	switch constant.SkillID(a.Skill) {
	case constant.SkillBigBang,
		constant.SkillBigBang2221001,
		constant.SkillBigBang2321001,
		constant.SkillPiercing,
		constant.SkillStormArrow,
		constant.SkillStormArrowCygnus,
		constant.SkillCorkscrewBlow,
		constant.SkillCorkscrewBlowCygnus,
		constant.SkillRapidFire,
		constant.SkillGrenade:
		sw.WriteU32(a.Charge)
	}

	sw.WriteU8(0)
	sw.WriteU8(a.Unk)
	sw.WriteU8(a.Speed)
	sw.WriteU8(a.Display)
	sw.WriteU32(a.LastTick)

	if a.Slot != 0 || a.CsStar != 0 || a.AOE != 0 {
		sw.WriteU16(a.Slot)
		sw.WriteU16(a.CsStar)
		sw.WriteU8(a.AOE)
	}

	for _, damage := range a.Damages {
		sw.WriteU32(damage.OID)
		sw.Write(make([]byte, 14))

		for _, pair := range damage.DamagePairs {
			sw.WriteU32(pair.Damage)
		}
	}

	if a.Position.X != 0 || a.Position.Y != 0 {
		sw.Write16(a.Position.X)
		sw.Write16(a.Position.Y)
	}

	return nil
}
