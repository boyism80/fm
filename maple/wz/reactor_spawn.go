package wz

import (
	"time"

	"github.com/boyism80/fm/types"
)

type ReactorSpawn struct {
	ReactorID       uint32
	Position        types.Vector2[int16]
	FacingDirection FacingDirectionType
	RespawnDelay    time.Duration
	Name            string
}
