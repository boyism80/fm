package actor

import (
	"github.com/boyism80/fm/game/entity"
	lua "github.com/yuin/gopher-lua"
)

// ResumeLua is sent to the map actor after sleep(duration); receiver resumes the thread.
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
