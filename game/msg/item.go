package msg

import "github.com/boyism80/fm/common/types"

type ItemSpawn struct {
	OwnerId      uint32
	SpawnedPoint types.Vector2[int16]
}
