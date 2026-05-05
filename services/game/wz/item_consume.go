package wz

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

const ConsumeMoveToReturnMap int32 = 999999999

type Consume struct {
	*ItemCore
	ActiveEffect    ActiveEffect
	BuffDuration    time.Duration
	BuffValues      map[constant.BuffFlag]int32
	ScrollSuccess   int32
	ScrollCursed    int32
	ScrollRandStat  int32
	ScrollRecover   int32
	ScrollIncStr    int16
	ScrollIncDex    int16
	ScrollIncInt    int16
	ScrollIncLuk    int16
	ScrollIncMaxHP  int16
	ScrollIncMaxMP  int16
	ScrollIncPAD    int16
	ScrollIncMAD    int16
	ScrollIncPDD    int16
	ScrollIncMDD    int16
	ScrollIncACC    int16
	ScrollIncAvoid  int16
	ScrollIncHands  int16
	ScrollIncSpeed  int16
	ScrollIncJump   int16
	ConsumeOnPickup bool
	Party           bool
	MoveTo          int32
	ExpInc          int32

	CureDebuffs []constant.DebuffFlag
}

func (c *Consume) IsShuriken() bool {
	return constant.GetConsumeType(c.ItemCore.ID) == constant.ConsumeTypeShuriken
}

func (c *Consume) IsBullet() bool {
	return constant.GetConsumeType(c.ItemCore.ID) == constant.ConsumeTypeBullet
}

func (c *Consume) IsArrowForBow() bool {
	return constant.GetConsumeType(c.ItemCore.ID) == constant.ConsumeTypeArrowBow
}

func (c *Consume) IsArrowForCrossBow() bool {
	return constant.GetConsumeType(c.ItemCore.ID) == constant.ConsumeTypeArrowCrossBow
}

func (c *Consume) IsConsumeOnPickup() bool {
	return c != nil && c.ConsumeOnPickup
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
