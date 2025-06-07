package entity

import (
	"github.com/boyism80/fm/game/data"
)

type Npc struct {
	ID   uint32
	Spec *data.NpcSpawnSpec
}
