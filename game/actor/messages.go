package actor

import (
	"github.com/boyism80/fm/game/constant"
	"github.com/boyism80/fm/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type SpawnDoor struct {
	OwnerID        uint32
	SkillID        constant.SkillID
	FieldMapID     uint32
	ReturnPortalID uint8
	FieldPortalID  uint8
}

type RemoveDoor struct {
	OwnerID uint32
	SkillID uint32
}

type ResumeLua struct {
	Root   *lua.LState
	Thread *lua.LState
}

type AddCharacter struct {
	Character  *entity.Character
	SpawnPoint uint8
	Init       bool
}

type RemoveCharacter struct {
	CharacterID uint32
}

type WarpCharacter struct {
	Character *entity.Character
	Portal    uint8
}

type TimerTick struct {
	HandlerName string
}
