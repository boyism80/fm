package entity

import (
	"time"

	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	"github.com/boyism80/fm/types"
)

type Mist struct {
	ObjectCore
	Causer               uint32
	SkillWz              *wz.Skill
	MobLevelData         *wz.MobSkillLevelData
	SkillLevel           uint8
	MistType             constant.MistType
	MobMist              bool
	MobSkill             bool
	SkillDelay           uint16
	Bounds               types.Rect[int32]
	ExpiresAt            time.Time
	NextPoisonTickAt     time.Time
	PoisonTickMultiplier float64
}

func MistWorldBounds(pos types.Point[int16], ld *wz.SkillLevelData) types.Rect[int32] {
	px := int32(pos.X)
	py := int32(pos.Y)
	if ld == nil {
		return types.Rect[int32]{Left: px, Top: py, Right: px, Bottom: py}
	}
	x1 := px + ld.LT.X
	x2 := px + ld.RB.X
	y1 := py + ld.LT.Y
	y2 := py + ld.RB.Y
	return types.Rect[int32]{
		Left:   min(x1, x2),
		Top:    min(y1, y2),
		Right:  max(x1, x2),
		Bottom: max(y1, y2),
	}
}

func (mist *Mist) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMist
}

func (mist *Mist) Is(typ constant.ObjectType) bool {
	return mist.GetObjectType().Has(typ)
}

func (mist *Mist) IsExpired(now time.Time) bool {
	return !mist.ExpiresAt.IsZero() && !now.Before(mist.ExpiresAt)
}

func (mist *Mist) SendSpawnSyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	skillID := uint32(0)
	if mist.SkillWz != nil {
		skillID = mist.SkillWz.ID
	}
	viewer.Send(&response.SpawnMist{
		OID:        mist.OID,
		Type:       mist.MistType,
		MobMist:    mist.MobMist,
		CauserID:   mist.Causer,
		SkillID:    skillID,
		SkillLevel: mist.SkillLevel,
		SkillDelay: mist.SkillDelay,
		Bounds:     mist.Bounds,
		MobSkill:   mist.MobSkill,
	}, types.SEND_POLICY_ENCRYPT)
}
