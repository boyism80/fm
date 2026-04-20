package dto

import (
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/stream"
	"github.com/boyism80/fm/types"
)

type CloseAttackInfo struct {
	AttackHeader
	Damages  []AttackPair
	MesoOIDs []uint32
	Position types.Vector2[int16]
}

func (a *CloseAttackInfo) deserializeMesoExplosion(sr *stream.StreamReader) {
	damages := make([]AttackPair, 0, a.Targets)
	if a.Hits == 0 {
		sr.Skip(4)
		bullets := sr.ReadU8()
		a.MesoOIDs = make([]uint32, 0, bullets)
		for range int(bullets) {
			a.MesoOIDs = append(a.MesoOIDs, sr.ReadU32())
			sr.Skip(1)
		}
		sr.Skip(2)
		a.Damages = damages
		return
	}

	for range int(a.Targets) {
		oid := sr.ReadU32()
		sr.Skip(12)
		bullets := sr.ReadU8()
		damagePairs := make([]DamagePair, 0, bullets)
		for range int(bullets) {
			damagePairs = append(damagePairs, DamagePair{Damage: sr.ReadU32(), Unknown: false})
		}
		damages = append(damages, AttackPair{OID: oid, DamagePairs: damagePairs})
	}

	sr.Skip(4)
	bullets := sr.ReadU8()
	a.MesoOIDs = make([]uint32, 0, bullets)
	for range int(bullets) {
		a.MesoOIDs = append(a.MesoOIDs, sr.ReadU32())
		sr.Skip(1)
	}
	a.Damages = damages
}

func (a *CloseAttackInfo) Deserialize(sr *stream.StreamReader) {
	a.AttackHeader.deserialize(sr)
	a.MesoOIDs = nil
	if constant.SkillID(a.Skill) == constant.SkillMesoExplosion {
		a.deserializeMesoExplosion(sr)
	} else {
		a.Damages = a.AttackHeader.parseNormalDamages(sr, 14)
		if sr.Remaining() >= 4 {
			a.Position = types.Vector2[int16]{X: sr.Read16(), Y: sr.Read16()}
		}
	}
}

func (a CloseAttackInfo) ToAttackInfo() AttackInfo {
	return AttackInfo{
		AttackHeader: a.AttackHeader,
		Damages:      a.Damages,
		MesoOIDs:     a.MesoOIDs,
		Position:     a.Position,
	}
}
