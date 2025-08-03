package entity

import (
	"github.com/boyism80/fm/common/types"
)

type Object struct {
	OID      uint32
	Position types.Vector2[int16]
	Context  GameContext // GameContext for accessing game resources
}
