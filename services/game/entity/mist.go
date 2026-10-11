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

func (mist *Mist) SendDestroySyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	viewer.Send(&response.RemoveMist{
		OID:      mist.OID,
		Eruption: false,
	}, types.SEND_POLICY_ENCRYPT)
}
