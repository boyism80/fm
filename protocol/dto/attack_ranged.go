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
	if sr.Remaining() >= 4 {
		a.Position = types.Vector2[int16]{X: sr.Read16(), Y: sr.Read16()}
	}
}

func (a RangedAttackInfo) ToAttackInfo() AttackInfo {
	return AttackInfo{
		AttackHeader: a.AttackHeader,
		Slot:         a.Slot,
		CsStar:       a.CsStar,
		AOE:          a.AOE,
		Damages:      a.Damages,
		Position:     a.Position,
	}
}
