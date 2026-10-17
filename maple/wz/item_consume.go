package wz

import (
	"math"
	"time"

	"github.com/boyism80/fm/maple/constant"
)

const ConsumeMoveToReturnMap int32 = 999999999

type Consume struct {
	*ItemCore
	ConsumeType    constant.ConsumeType
	UnitPrice      float64
	ActiveEffect   ActiveEffect
	BuffDuration   time.Duration
	BuffValues     map[constant.BuffFlag]int32
	ScrollSuccess  int32
	ScrollCursed   int32
	ScrollRandStat int32
	ScrollRecover  int32
	ScrollFlag     constant.ItemFlag
	ScrollReqRUC   int32
	ScrollReqs     []uint32
	ScrollIncStr   int16
	ScrollIncDex   int16
	ScrollIncInt   int16
	ScrollIncLuk   int16
	ScrollIncMaxHP int16
	ScrollIncMaxMP int16
	ScrollIncPAD   int16
	ScrollIncMAD   int16
	ScrollIncPDD   int16
	ScrollIncMDD   int16
	ScrollIncACC   int16
	ScrollIncAvoid int16
	ScrollIncHands int16
	ScrollIncSpeed int16
	ScrollIncJump  int16
	PetFood        *PetFood
	Party          bool
	MoveTo         int32
	ExpInc         int32
	CP             int
	NuffSkillID    uint32
	NuffSkillLevel uint8
	MobID          uint32
	MobHP          int32
	CreateID       uint32
	MonsterBook    bool
	MountFatigue   int
	SkillBook      *SkillBook
	SummonBag      []SummonBagEntry
	Rewards        []RewardEntry
	RewardProb     int

	CureDebuffs   []constant.DebuffFlag
	AvailableMaps []MapRange
}

func (model *Consume) setConsumeType() {
	switch model.ID / 10000 {
	case 203:
		model.ConsumeType = constant.ConsumeTypeReturnScroll
	case 204:
		model.ConsumeType = constant.ConsumeTypeEnhanceScroll
	case 206:
		if model.ID >= 2061000 && model.ID < 2062000 {
			model.ConsumeType = constant.ConsumeTypeArrowCrossBow
		} else {
			model.ConsumeType = constant.ConsumeTypeArrowBow
		}
	case 207:
		model.ConsumeType = constant.ConsumeTypeShuriken
	case 210:
		model.ConsumeType = constant.ConsumeTypeMonsterSummon
	case 212:
		model.ConsumeType = constant.ConsumeTypePetFood
	case 224:
		model.ConsumeType = constant.ConsumeTypeWeddingRing
	case 226:
		model.ConsumeType = constant.ConsumeTypeRidingFood
	case 228:
		model.ConsumeType = constant.ConsumeTypeSkillBookGet
	case 229:
		model.ConsumeType = constant.ConsumeTypeSkillBookUnlock
	case 231:
		model.ConsumeType = constant.ConsumeTypeMinervaOwl
	case 232:
		model.ConsumeType = constant.ConsumeTypeTeleportStone
	case 233:
		model.ConsumeType = constant.ConsumeTypeBullet
	case 238:
		model.ConsumeType = constant.ConsumeTypeCard
	default:
		model.ConsumeType = constant.ConsumeTypeNone
	}
}

type SkillBook struct {
	Success       int
	RequiredLevel int
	MasterLevel   int
	Skills        []uint32
}

type SummonBagEntry struct {
	MobID       uint32
	Probability int
}

type RewardEntry struct {
	ItemID uint32
	Count  uint16
	Prob   int
	Period time.Duration
	Effect string
}

type MapRange struct {
	Start uint32
	End   uint32
}

func (c *Consume) RechargePrice(missing uint16) int {
	return int(math.Round(c.UnitPrice * float64(missing)))
}

func (c *Consume) IsAvailableAt(mapID uint32) bool {
	if len(c.AvailableMaps) == 0 {
		return true
	}
	for _, r := range c.AvailableMaps {
		if mapID >= r.Start && mapID <= r.End {
			return true
		}
	}
	return false
}

func (c *Consume) BuffSpecValues() map[constant.BuffFlag]int32 {
	if c == nil || c.BuffDuration <= 0 || len(c.BuffValues) == 0 {
		return nil
	}
	out := make(map[constant.BuffFlag]int32, len(c.BuffValues))
	for k, v := range c.BuffValues {
		out[k] = v
	}
	return out
}
