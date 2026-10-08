package entity

import (
	lua "github.com/yuin/gopher-lua"
)

func (r *Records) LuaTypeName() string {
	return "LuaRecords"
}

func (r *Records) String() string {
	return r.LuaTypeName()
}

func (r *Records) Type() lua.LValueType {
	return lua.LTUserData
}

func LuaCheckRecords(L *lua.LState, idx int) (*Records, bool) {
	ud := L.CheckUserData(idx)
	r, ok := ud.Value.(*Records)
	if !ok || r == nil {
		L.ArgError(idx, "LuaRecords expected")
		return nil, false
	}
	return r, true
}

func (r *Records) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"get": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			L.Push(lua.LNumber(records.Get(L.CheckString(2))))
			return 1
		},
		"text": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			L.Push(lua.LString(records.Text(L.CheckString(2))))
			return 1
		},
		"set": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			records.Set(L.CheckString(2), int64(L.CheckNumber(3)), RecordPeriod(L.OptInt(4, int(RecordPeriodNone))))
			return 0
		},
		"set_text": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			records.SetText(L.CheckString(2), L.CheckString(3), RecordPeriod(L.OptInt(4, int(RecordPeriodNone))))
			return 0
		},
		"add": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			value := records.Add(L.CheckString(2), int64(L.OptNumber(3, 1)), RecordPeriod(L.OptInt(4, int(RecordPeriodNone))))
			L.Push(lua.LNumber(value))
			return 1
		},
		"remove": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			records.Remove(L.CheckString(2))
			return 0
		},
		"clear": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			records.Clear()
			return 0
		},
	}
}
