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

type AttackHeader struct {
	Targets  uint8
	Hits     uint8
	Skill    uint32
	Charge   uint32
	Unk      uint8
	Speed    uint8
	Display  uint8
	LastTick uint32
}

type AttackPayload interface {
	ToAttackInfo() AttackInfo
}

type AttackInfo struct {
	AttackHeader
	Slot     uint16
	CsStar   uint16
	AOE      uint8
	Damages  []AttackPair
	MesoOIDs []uint32
	Position types.Vector2[int16]
}

func (header *AttackHeader) deserialize(sr *stream.StreamReader) {
	sr.Skip(1)
	tbyte := sr.ReadU8()
	header.Targets = (tbyte >> 4) & 0xF
	header.Hits = tbyte & 0xF
	header.Skill = sr.ReadU32()

	switch constant.SkillID(header.Skill) {
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
		header.Charge = sr.ReadU32()
	default:
		header.Charge = 0
	}

	sr.Skip(1)
	header.Unk = sr.ReadU8()
	header.Speed = sr.ReadU8()
	header.Display = sr.ReadU8()
	header.LastTick = sr.ReadU32()
}

func (header *AttackHeader) parseNormalDamages(sr *stream.StreamReader, skipBytes int) []AttackPair {
	damages := make([]AttackPair, 0, header.Targets)
	for range int(header.Targets) {
		oid := sr.ReadU32()
		sr.Skip(skipBytes)
		damagePairs := make([]DamagePair, 0, header.Hits)
		for range int(header.Hits) {
			damagePairs = append(damagePairs, DamagePair{Damage: sr.ReadU32(), Unknown: false})
		}
		damages = append(damages, AttackPair{OID: oid, DamagePairs: damagePairs})
	}
	return damages
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
