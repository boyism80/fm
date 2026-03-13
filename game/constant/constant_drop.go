package constant

import "time"

// DropItemAnimationType represents the animation when items are dropped.
type DropItemAnimationType uint8

const (
	DROP_ITEM_ANIMATION_TYPE_LOOTING        DropItemAnimationType = 1 // Item being looted
	DROP_ITEM_ANIMATION_TYPE_NONE           DropItemAnimationType = 2 // No animation
	DROP_ITEM_ANIMATION_TYPE_DISAPPEAR_FADE DropItemAnimationType = 3 // Fade out disappear
	DROP_ITEM_ANIMATION_TYPE_DISAPPEAR      DropItemAnimationType = 4 // Instant disappear
)

// DropType represents who can pick up a dropped item.
type DropType uint8

const (
	DROP_TYPE_FFA   DropType = 0 // Free for all
	DROP_TYPE_PARTY DropType = 1 // Party members only
	DROP_TYPE_OWNED DropType = 2 // Owner only
)

// Item cleanup times
const (
	ITEM_EXPIRE_TIME = 120 * time.Second // Time before item/meso expires
	ITEM_FFA_TIME    = 30 * time.Second  // Time before owned/party drop becomes FFA
)
