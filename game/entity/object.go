package entity

import "github.com/boyism80/fm/common/types"

type Object struct {
	ID       int64
	Position types.Vec2
	Name     string
	Type_    types.ObjectType
}
