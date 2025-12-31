package entity

import (
	"github.com/boyism80/fm/game/wz"
)

type Npc struct {
	Object
	Wz *wz.NpcSpawn
}
