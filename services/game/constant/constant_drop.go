package constant

import "time"

type DropItemAnimationType uint8

const (
	DropItemAnimationTypeLooting       DropItemAnimationType = 1
	DropItemAnimationTypeNone          DropItemAnimationType = 2
	DropItemAnimationTypeDisappearFade DropItemAnimationType = 3
	DropItemAnimationTypeDisappear     DropItemAnimationType = 4
)

type DropType uint8

const (
	DropTypeFFA   DropType = 0
	DropTypeParty DropType = 1
	DropTypeOwned DropType = 2
)

const (
	ItemExpireTime = 120 * time.Second
	ItemFFATime    = 30 * time.Second
)
