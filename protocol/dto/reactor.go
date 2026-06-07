package dto

import "github.com/boyism80/fm/types"

type Reactor struct {
	OID       uint32
	ReactorID uint32
	State     uint8
	Position  types.Vector2[int16]
	Facing    uint8
	Name      string
}
