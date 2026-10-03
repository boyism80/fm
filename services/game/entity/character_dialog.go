package entity

import (
	"errors"
	"fmt"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	lua "github.com/yuin/gopher-lua"
)

var (
	ErrDialogNotExpected = errors.New("dialog answer not expected")
	ErrDialogChoiceRange = errors.New("dialog choice out of range")
	ErrDialogNoLuaRoot   = errors.New("dialog lua root not found")
)

type npcDialog struct {
	dialogType constant.DialogType
	choices    int
}

func (ch *Character) GetDialog() *lua.LState {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	return ch.luaDialog
}

func (ch *Character) SetDialog(thread *lua.LState) {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.luaDialog = thread
	ch.npcDialog = nil
}

func (ch *Character) AskDialog(thread *lua.LState, dialogType constant.DialogType, choices int) {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.luaDialog = thread
	ch.npcDialog = &npcDialog{dialogType: dialogType, choices: choices}
}

func (ch *Character) ResetDialog() {
	ch.dialogMutex.Lock()
	defer ch.dialogMutex.Unlock()
	ch.luaDialog = nil
	ch.npcDialog = nil
}

func (ch *Character) CloseDialog() {
	thread := ch.GetDialog()
	if thread == nil {
		return
	}
	ch.ResetDialog()
	if cfg, ok := luax.GetConfiguration(thread); ok {
		cfg.CallPromise = nil
		luax.SetConfiguration(thread, cfg)
	}
	luax.Close(thread)
}

func (ch *Character) ResumeDialog(actx actor.Context, dialogType constant.DialogType, next bool, selected uint32, text string) error {
	ch.dialogMutex.Lock()
	thread := ch.luaDialog
	asked := ch.npcDialog
	if thread == nil || asked == nil || asked.dialogType != dialogType {
		ch.dialogMutex.Unlock()
		return ErrDialogNotExpected
	}
	ch.npcDialog = nil
	ch.dialogMutex.Unlock()

	var args []interface{}
	switch dialogType {
	case constant.DialogTypeDefault, constant.DialogTypeYesNo, constant.DialogTypeAccept:
		args = append(args, lua.LBool(next))
	case constant.DialogTypeList, constant.DialogTypeStyle:
		if next == false {
			args = append(args, lua.LNil)
			break
		}
		if selected >= uint32(asked.choices) {
			ch.CloseDialog()
			ch.Listener.OnUnlockAction(ch)
			return ErrDialogChoiceRange
		}
		args = append(args, lua.LNumber(selected+1))
	case constant.DialogTypeInput:
		if next {
			args = append(args, lua.LString(text))
		} else {
			args = append(args, lua.LNil)
		}
	}

	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return ErrDialogNoLuaRoot
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return ErrDialogNoLuaRoot
	}
	cfg, _ := luax.GetConfiguration(thread)
	cfg.ActorContext = actx
	luax.SetConfiguration(thread, cfg)

	state, _, err := luax.Resume(root, thread, args...)
	if err != nil {
		ch.ResetDialog()
		return fmt.Errorf("resume dialog: %w", err)
	}
	if state != lua.ResumeYield {
		ch.ResetDialog()
	}
	return nil
}
