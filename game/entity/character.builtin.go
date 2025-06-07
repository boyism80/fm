package entity

import (
	"time"

	lua "github.com/yuin/gopher-lua"
)

var characterBuiltinFuncs = map[string]lua.LGFunction{
	"sleep": func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}
		ms := L.CheckInt(2)
		if ms == 0 {
			return 0
		}

		ch.Listener.OnSleep(time.Duration(ms)*time.Millisecond, L, []lua.LValue{})
		return L.Yield(lua.LNumber(0))
	},
	"hp": func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}
		top := L.GetTop()
		if top == 1 {
			L.Push(lua.LNumber(ch.Hp))
			return 1
		}

		newHp := L.CheckInt(2)
		ch.Hp = uint16(newHp)
		return 0
	},
	"name": func(L *lua.LState) int {
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}
		L.Push(lua.LString(ch.Name))
		return 1
	},
	"meso": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		if argc > 1 {
			ch.Meso = int32(L.CheckInt(2))
			ch.Listener.OnMesoChanged(ch.Meso)
			return 0
		} else {
			L.Push(lua.LNumber(ch.Meso))
			return 1
		}
	},
	"chat": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		message := L.CheckString(2)
		highlight := false
		if argc > 2 {
			highlight = L.CheckBool(3)
		}
		dontRecordHistory := false
		if argc > 3 {
			dontRecordHistory = L.CheckBool(4)
		}
		if ch.Listener != nil {
			ch.Listener.OnChat(message, highlight, dontRecordHistory)
		}

		ch.Dialog = L
		return 0
	},
	"dialog": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		npc := 0
		if argc > 1 {
			npc = L.CheckInt(2)
		}

		message := ""
		if argc > 2 {
			message = L.CheckString(3)
		}

		prev := false
		if argc > 3 {
			prev = L.CheckBool(4)
		}
		next := false
		if argc > 4 {
			next = L.CheckBool(5)
		}
		if ch.Listener != nil {
			ch.Listener.OnDialog(uint32(npc), message, prev, next)
		}

		ch.Dialog = L
		return L.Yield(lua.LNumber(1))
	},
	"dialog_list": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		npc := 0
		if argc > 1 {
			npc = L.CheckInt(2)
		}

		message := ""
		if argc > 2 {
			message = L.CheckString(3)
		}

		selections := []string{}
		if argc > 3 {
			tbl := L.CheckTable(4)
			tbl.ForEach(func(_, value lua.LValue) {
				if str, ok := value.(lua.LString); ok {
					selections = append(selections, string(str))
				}
			})
		}
		if ch.Listener != nil {
			ch.Listener.OnDialogList(uint32(npc), message, selections)
		}

		ch.Dialog = L
		return L.Yield(lua.LNumber(1))
	},
	"dialog_accept": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		npc := 0
		if argc > 1 {
			npc = L.CheckInt(2)
		}

		message := ""
		if argc > 2 {
			message = L.CheckString(3)
		}

		enableEscape := false
		if argc > 3 {
			enableEscape = L.CheckBool(4)
		}
		if ch.Listener != nil {
			ch.Listener.OnDialogAccept(uint32(npc), message, enableEscape)
		}

		ch.Dialog = L
		return L.Yield(lua.LNumber(1))
	},
	"dialog_yes_no": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		npc := 0
		if argc > 1 {
			npc = L.CheckInt(2)
		}

		message := ""
		if argc > 2 {
			message = L.CheckString(3)
		}

		prev := false
		if argc > 3 {
			prev = L.CheckBool(4)
		}
		next := false
		if argc > 4 {
			next = L.CheckBool(5)
		}
		if ch.Listener != nil {
			ch.Listener.OnDialogYesNo(uint32(npc), message, prev, next)
		}

		ch.Dialog = L
		return L.Yield(lua.LNumber(1))
	},
	"dialog_input": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*Character)
		if !ok {
			L.ArgError(1, "Character expected")
			return 0
		}

		npc := 0
		if argc > 1 {
			npc = L.CheckInt(2)
		}

		message := ""
		if argc > 2 {
			message = L.CheckString(3)
		}
		if ch.Listener != nil {
			ch.Listener.OnDialogInput(uint32(npc), message)
		}

		ch.Dialog = L
		return L.Yield(lua.LNumber(1))
	},
}
