package entity

import "github.com/boyism80/fm/common/types"

type Object struct {
	ID       uint32
	Position types.Vec2[int16]
	Type_    types.ObjectType
}
