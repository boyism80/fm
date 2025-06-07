package entity

import (
	"time"

	"github.com/asynkron/protoactor-go/actor"
	lua "github.com/yuin/gopher-lua"
)

type CharacterListener interface {
	GetContext() actor.Context
	OnDialog(npc uint32, message string, prev bool, next bool)
	OnDialogYesNo(npc uint32, message string, prev bool, next bool)
	OnDialogAccept(npc uint32, message string, enableEscape bool)
	OnDialogList(npc uint32, message string, selections []string)
	OnDialogInput(npc uint32, message string)
	OnChat(message string, highlight bool, dontRecordHistory bool)
	OnMesoChanged(meso int32)
	OnSleep(duration time.Duration, L *lua.LState, args []lua.LValue)
}
