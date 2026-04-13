package entity

import (
	"github.com/boyism80/fm/services/game/constant"
)

type DoorSpawn struct {
	OwnerID        uint32
	SkillID        constant.SkillID
	FieldMapID     uint32
	ReturnPortalID uint8
	FieldPortalID  uint8
}

type DoorRemove struct {
	OwnerID            uint32
	SkillID            uint32
	CounterpartMapWZID uint32
}
