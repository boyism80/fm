package game

import (
	"time"

	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/types"
)

type Mist struct {
	objectCore
	causer               uint32
	skillWz              *wz.Skill
	mobLevelData         *wz.MobSkillLevelData
	skillLevel           uint8
	mistType             constant.MistType
	mobMist              bool
	mobSkill             bool
	skillDelay           uint16
	bounds               types.Rect[int32]
	expiresAt            time.Time
	nextPoisonTickAt     time.Time
	poisonTickMultiplier float64
}

func (mist *Mist) GetObjectType() constant.ObjectType {
	return constant.ObjectTypeMist
}

func (mist *Mist) Is(typ constant.ObjectType) bool {
	return mist.GetObjectType().Has(typ)
}

func (mist *Mist) SchedulePoison(at time.Time) {
	mist.nextPoisonTickAt = at
}

func (mist *Mist) IsExpired(now time.Time) bool {
	return !mist.expiresAt.IsZero() && !now.Before(mist.expiresAt)
}

func (mist *Mist) SendSpawnSyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	skillID := uint32(0)
	if mist.skillWz != nil {
		skillID = mist.skillWz.ID
	}
	viewer.Send(&response.SpawnMist{
		OID:        mist.oid,
		Type:       mist.mistType,
		MobMist:    mist.mobMist,
		CauserID:   mist.causer,
		SkillID:    skillID,
		SkillLevel: mist.skillLevel,
		SkillDelay: mist.skillDelay,
		Bounds:     mist.bounds,
		MobSkill:   mist.mobSkill,
	}, types.SEND_POLICY_ENCRYPT)
}

func (mist *Mist) SendDestroySyncToViewer(viewer *Character) {
	if mist == nil || viewer == nil {
		return
	}
	viewer.Send(&response.RemoveMist{
		OID:      mist.oid,
		Eruption: false,
	}, types.SEND_POLICY_ENCRYPT)
}
