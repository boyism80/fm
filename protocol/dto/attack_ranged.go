package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type RangedAttackInfo struct {
	AttackHeader
	Slot     uint16
	CsStar   uint16
	AOE      uint8
	Damages  []AttackPair
	Position types.Vector2[int16]
}

func (a *RangedAttackInfo) Deserialize(sr *stream.StreamReader) {
	a.AttackHeader.deserialize(sr)
	a.Slot = sr.ReadU16()
	a.CsStar = sr.ReadU16()
	a.AOE = sr.ReadU8()
	a.Damages = a.AttackHeader.parseNormalDamages(sr, 14)
	a.Position = types.Vector2[int16]{X: sr.Read16(), Y: sr.Read16()}
}

func (a *RangedAttackInfo) Serialize(sw *stream.StreamWriter) error {
	a.AttackHeader.serialize(sw)
	sw.WriteU16(a.Slot)
	sw.WriteU16(a.CsStar)
	sw.WriteU8(a.AOE)
	a.AttackHeader.serializeNormalDamages(sw, a.Damages)
	sw.Write16(a.Position.X)
	sw.Write16(a.Position.Y)
	return nil
}

func (a RangedAttackInfo) ToAttackInfo() AttackInfo {
	return AttackInfo{
		AttackHeader: a.AttackHeader,
		Ranged:       true,
		Slot:         a.Slot,
		CsStar:       a.CsStar,
		AOE:          a.AOE,
		Damages:      a.Damages,
		Position:     a.Position,
	}
}
