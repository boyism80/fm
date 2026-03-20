package wz

import "github.com/boyism80/fm/game/constant"

type Consume struct {
	*ItemCore
	ActiveEffect ActiveEffect
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
