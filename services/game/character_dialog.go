package game

import (
	"errors"
	"fmt"
	"sync"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	lua "github.com/yuin/gopher-lua"
)

var (
	ErrDialogNotExpected = errors.New("dialog answer not expected")
	ErrDialogChoiceRange = errors.New("dialog choice out of range")
	ErrDialogNoLuaRoot   = errors.New("dialog lua root not found")
	ErrDialogBusy        = errors.New("dialog already open")
)

type npcDialog struct {
	dialogType constant.DialogType
	choices    int
	min, max   int32
}

type Dialog struct {
	owner   *Character
	mutex   sync.Mutex
	thread  *lua.LState
	asked   *npcDialog
	waiting bool
	shop    *wz.Shop
}

func NewDialog(owner *Character) *Dialog {
	return &Dialog{owner: owner}
}

func (d *Dialog) Thread() *lua.LState {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	return d.thread
}

func (d *Dialog) Busy() bool {
	return d.Thread() != nil || d.shop != nil
}

func (d *Dialog) OpenShop(shop *wz.Shop) {
	d.shop = shop
	d.owner.listener.OnOpenNpcShop(d.owner, shop.NpcID, shop)
}

func (d *Dialog) Shop() *wz.Shop {
	return d.shop
}

func (d *Dialog) CloseShop() {
	d.shop = nil
}

func (d *Dialog) Set(thread *lua.LState) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = thread
	d.asked = nil
	d.waiting = false
}

func (d *Dialog) Ask(thread *lua.LState, dialogType constant.DialogType, choices int) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = thread
	d.asked = &npcDialog{dialogType: dialogType, choices: choices}
	d.waiting = false
}

func (d *Dialog) AskNumber(thread *lua.LState, min, max int32) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = thread
	d.asked = &npcDialog{dialogType: constant.DialogTypeNumber, min: min, max: max}
	d.waiting = false
}

func (d *Dialog) Wait(thread *lua.LState) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = thread
	d.asked = nil
	d.waiting = true
}

func (d *Dialog) Reset() {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.thread = nil
	d.asked = nil
	d.waiting = false
}

func (d *Dialog) Close() {
	d.owner.account.storage.Close()
	d.shop = nil
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

// The reply can arrive after the dialog was closed or replaced, so only the thread still waiting for it resumes.
func (d *Dialog) Wake(thread *lua.LState, value lua.LValue) error {
	d.mutex.Lock()
	if d.thread != thread || !d.waiting {
		d.mutex.Unlock()
		return ErrDialogNotExpected
	}
	d.waiting = false
	d.mutex.Unlock()

	mapInstance := d.owner.GetMap()
	if mapInstance == nil {
		return ErrDialogNoLuaRoot
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return ErrDialogNoLuaRoot
	}
	if _, _, err := luax.Resume(root, thread, value); err != nil {
		return fmt.Errorf("wake dialog: %w", err)
	}
	return nil
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
		if !next {
			args = append(args, lua.LNil)
			break
		}
		if selected >= uint32(asked.choices) {
			d.Close()
			d.owner.listener.OnUnlockAction(d.owner)
			return ErrDialogChoiceRange
		}
		args = append(args, lua.LNumber(selected+1))
	case constant.DialogTypeInput:
		if next {
			args = append(args, lua.LString(text))
		} else {
			args = append(args, lua.LNil)
		}
	case constant.DialogTypeNumber:
		if !next {
			args = append(args, lua.LNil)
			break
		}
		if int32(selected) < asked.min || int32(selected) > asked.max {
			d.Close()
			d.owner.listener.OnUnlockAction(d.owner)
			return ErrDialogChoiceRange
		}
		args = append(args, lua.LNumber(int32(selected)))
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
