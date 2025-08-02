package entity

import (
	"github.com/boyism80/fm/game/data"
)

type Npc struct {
	Object
	Spec *data.NpcSpawnSpec
}
