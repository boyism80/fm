package entity

import (
	"errors"
	"fmt"
	"sync"

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

type Dialog struct {
	owner  *Character
	mutex  sync.Mutex
	thread *lua.LState
	asked  *npcDialog
	ShopID uint32
}

func NewDialog(owner *Character) *Dialog {
	return &Dialog{owner: owner}
}

func (d *Dialog) Thread() *lua.LState {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.thread
}

func (d *Dialog) Set(thread *lua.LState) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = thread
	d.asked = nil
}

func (d *Dialog) Ask(thread *lua.LState, dialogType constant.DialogType, choices int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = thread
	d.asked = &npcDialog{dialogType: dialogType, choices: choices}
}

func (d *Dialog) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = nil
	d.asked = nil
}

func (d *Dialog) Close() {
	thread := d.Thread()
	if thread == nil {
		return
	}
	d.Reset()
	if cfg, ok := luax.GetConfiguration(thread); ok {
		cfg.CallPromise = nil
		luax.SetConfiguration(thread, cfg)
	}
	luax.Close(thread)
}

func (d *Dialog) Resume(actx actor.Context, dialogType constant.DialogType, next bool, selected uint32, text string) error {
	d.mutex.Lock()
	thread := d.thread
	asked := d.asked
	if thread == nil || asked == nil || asked.dialogType != dialogType {
		d.mutex.Unlock()
		return ErrDialogNotExpected
	}
	d.asked = nil
	d.mutex.Unlock()

	var args []interface{}
	switch dialogType {
	case constant.DialogTypeDefault, constant.DialogTypeYesNo, constant.DialogTypeAccept, constant.DialogTypeAcceptEscape:
		args = append(args, lua.LBool(next))
	case constant.DialogTypeList, constant.DialogTypeStyle:
		if next == false {
			args = append(args, lua.LNil)
			break
		}
		if selected >= uint32(asked.choices) {
			d.Close()
			d.owner.Listener.OnUnlockAction(d.owner)
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

	mapInstance := d.owner.GetMap()
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
		if d.Thread() == thread {
			d.Reset()
		}
		return fmt.Errorf("resume dialog: %w", err)
	}
	if state != lua.ResumeYield && d.Thread() == thread {
		d.Reset()
	}
	return nil
}
