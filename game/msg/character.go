package msg

import (
	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/game/entity"
	lua "github.com/yuin/gopher-lua"
)

type CharacterLootFailed struct {
	OID uint32
}

type CharacterItemLooting struct {
	PID  *actor.PID
	OID  uint32
	Item entity.Item
}

type CharacterMesoLooting struct {
	PID  *actor.PID
	OID  uint32
	Meso int32
}

type CharacterMapChanged struct {
	MID        uint32
	Map        *actor.PID
	Init       bool
	SpawnPoint uint8
}

type RunScript struct {
	Script string
}

type ResumeScript struct {
	L    *lua.LState
	Args []any
}

type CharacterBuiltinDialog struct {
	NPC     uint32
	Message string
	Prev    bool
	Next    bool
}

type CharacterBuiltinDialogList struct {
	NPC        uint32
	Message    string
	Selections []string
}

type CharacterBuiltinDialogAccept struct {
	NPC          uint32
	Message      string
	EnableEscape bool
}

type CharacterBuiltinDialogInput struct {
	NPC     uint32
	Message string
}

type CharacterBuiltinChat struct {
	Message           string
	Highlight         bool
	DontRecordHistory bool
}
