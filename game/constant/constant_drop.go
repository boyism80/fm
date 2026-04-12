package constant

import "time"

type DropItemAnimationType uint8

const (
	DROP_ITEM_ANIMATION_TYPE_LOOTING        DropItemAnimationType = 1
	DROP_ITEM_ANIMATION_TYPE_NONE           DropItemAnimationType = 2
	DROP_ITEM_ANIMATION_TYPE_DISAPPEAR_FADE DropItemAnimationType = 3
	DROP_ITEM_ANIMATION_TYPE_DISAPPEAR      DropItemAnimationType = 4
)

type DropType uint8

const (
	DROP_TYPE_FFA   DropType = 0
	DROP_TYPE_PARTY DropType = 1
	DROP_TYPE_OWNED DropType = 2
)

const (
	ITEM_EXPIRE_TIME = 120 * time.Second
	ITEM_FFA_TIME    = 30 * time.Second
)
