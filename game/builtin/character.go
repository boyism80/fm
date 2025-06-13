package builtin

import (
	"github.com/boyism80/fm/game/msg"
	lua "github.com/yuin/gopher-lua"
)

var characterBuiltinFuncs = map[string]lua.LGFunction{

	"name": func(L *lua.LState) int {
		// ud := L.CheckUserData(1)
		// ch, ok := ud.Value.(*Character)
		// if !ok {
		// 	L.ArgError(1, "Character expected")
		// 	return 0
		// }
		// L.Push(lua.LString(ch.Name))
		// return 1
		return 0
	},
	"meso": func(L *lua.LState) int {
		// argc := L.GetTop()
		// ud := L.CheckUserData(1)
		// ch, ok := ud.Value.(*Character)
		// if !ok {
		// 	L.ArgError(1, "Character expected")
		// 	return 0
		// }

		// if argc > 1 {
		// 	ch.Meso = int32(L.CheckInt(2))
		// 	ch.Listener.OnMesoChanged(ch.Meso)
		// 	return 0
		// } else {
		// 	L.Push(lua.LNumber(ch.Meso))
		// 	return 1
		// }
		return 0
	},
	"chat": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*LuaCharacter)
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

		ch.Context.Send(ch.PID, &msg.CharacterBuiltinChat{
			Message:           message,
			Highlight:         highlight,
			DontRecordHistory: dontRecordHistory,
		})

		return 0
	},
	"dialog": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*LuaCharacter)
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

		ch.Context.Send(ch.PID, &msg.CharacterBuiltinDialog{
			NPC:     uint32(npc),
			Message: message,
			Prev:    prev,
			Next:    next,
		})

		return L.Yield(lua.LNumber(1))
	},
	"dialog_list": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*LuaCharacter)
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

		ch.Context.Send(ch.PID, &msg.CharacterBuiltinDialogList{
			NPC:        uint32(npc),
			Message:    message,
			Selections: selections,
		})

		return L.Yield(lua.LNumber(1))
	},
	"dialog_accept": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*LuaCharacter)
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

		ch.Context.Send(ch.PID, &msg.CharacterBuiltinDialogAccept{
			NPC:          uint32(npc),
			Message:      message,
			EnableEscape: enableEscape,
		})

		return L.Yield(lua.LNumber(1))
	},
	"dialog_yes_no": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*LuaCharacter)
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

		ch.Context.Send(ch.PID, &msg.CharacterBuiltinDialog{
			NPC:     uint32(npc),
			Message: message,
			Prev:    prev,
			Next:    next,
		})

		return L.Yield(lua.LNumber(1))
	},
	"dialog_input": func(L *lua.LState) int {
		argc := L.GetTop()
		ud := L.CheckUserData(1)
		ch, ok := ud.Value.(*LuaCharacter)
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

		ch.Context.Send(ch.PID, &msg.CharacterBuiltinDialogInput{
			NPC:     uint32(npc),
			Message: message,
		})

		return L.Yield(lua.LNumber(1))
	},
}

type LuaCharacter struct {
	LuaLife
}

func (ch *LuaCharacter) LuaTypeName() string {
	return "LuaCharacter"
}

func (ch *LuaCharacter) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return characterBuiltinFuncs
}
