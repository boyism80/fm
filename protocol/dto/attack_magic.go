package dto

import (
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type MagicAttackInfo struct {
	AttackHeader
	Damages  []AttackPair
	Position types.Vector2[int16]
}

func (a *MagicAttackInfo) Deserialize(sr *stream.StreamReader) {
	a.AttackHeader.deserialize(sr)
	a.Damages = a.AttackHeader.parseNormalDamages(sr, 14)
	if sr.Remaining() >= 4 {
		a.Position = types.Vector2[int16]{X: sr.Read16(), Y: sr.Read16()}
	}
}

func (a MagicAttackInfo) ToAttackInfo() AttackInfo {
	return AttackInfo{
		AttackHeader: a.AttackHeader,
		Damages:      a.Damages,
		Position:     a.Position,
	}
}
