package wz

import (
	"time"

	"github.com/boyism80/fm/services/game/constant"
)

type Consume struct {
	*ItemCore
	ActiveEffect ActiveEffect
	BuffDuration time.Duration
	BuffValues   map[constant.BuffFlag]int32
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
