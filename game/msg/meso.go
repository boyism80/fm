package msg

import (
	"github.com/boyism80/fm/common/types"
)

type MesoSpawn struct {
	OwnerId      uint32
	Position     types.Vector2[int16]
	SpawnedPoint types.Vector2[int16]
}
