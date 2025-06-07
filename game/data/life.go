package data

import "github.com/boyism80/fm/common/types"

type FacingDirectionType uint8

const (
	FACING_DIRECTION_RIGHT FacingDirectionType = iota
	FACING_DIRECTION_LEFT
)

type Spawn interface {
}

type SpawnSpec struct {
	ID                 uint32
	Position           types.Vector2[int16]
	RenderX0, RenderX1 int16
	CollisionY         int16
	Hide               bool
	UseDay, UseNight   bool
	Foothold           int16
	FacingDirection    FacingDirectionType
	MobTime            uint64
	Info               uint8
	LimitedName        string
	NoFoothold         bool
}
