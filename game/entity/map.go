package entity

import (
	"time"

	"github.com/boyism80/fm/game/data"
)

type MobSpawn struct {
	Spec          *data.MobSpawnSpec
	Spawned       bool
	LastSpawnedAt time.Time
}

type Map struct {
	MobSpawns map[uint32]*MobSpawn
}
