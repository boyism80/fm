package entity

import (
	"strings"
	"time"

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
			var reset RecordReset
			reset.ParseLua(L, 4)
			records.Set(L.CheckString(2), int64(L.CheckNumber(3)), reset)
			return 0
		},
		"set_text": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			var reset RecordReset
			reset.ParseLua(L, 4)
			records.SetText(L.CheckString(2), L.CheckString(3), reset)
			return 0
		},
		"add": func(L *lua.LState) int {
			records, ok := LuaCheckRecords(L, 1)
			if !ok {
				return 0
			}
			var reset RecordReset
			reset.ParseLua(L, 4)
			value := records.Add(L.CheckString(2), int64(L.OptNumber(3, 1)), reset)
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

func (reset *RecordReset) ParseLua(L *lua.LState, idx int) {
	opts := L.OptTable(idx, nil)
	if opts == nil {
		return
	}
	reset.Restart = lua.LVAsBool(opts.RawGetString("restart"))

	if daily := opts.RawGetString("daily"); daily != lua.LNil {
		reset.Kind = RecordResetDaily
		if text, ok := daily.(lua.LString); ok {
			reset.At = reset.parseClock(L, idx, string(text))
		}
		return
	}

	if weekly := opts.RawGetString("weekly"); weekly != lua.LNil {
		reset.Kind = RecordResetWeekly
		text, ok := weekly.(lua.LString)
		if !ok {
			return
		}
		fields := strings.Fields(string(text))
		day := -1
		if len(fields) != 0 {
			for d := time.Sunday; d <= time.Saturday; d++ {
				if strings.EqualFold(d.String()[:3], fields[0]) {
					day = (int(d) + 6) % 7
				}
			}
		}
		if day == -1 {
			L.ArgError(idx, "weekly must start with a weekday such as Thu")
			return
		}
		reset.At = time.Duration(day) * 24 * time.Hour
		if len(fields) > 1 {
			reset.At += reset.parseClock(L, idx, fields[1])
		}
		return
	}

	if every := opts.RawGetString("every"); every != lua.LNil {
		seconds, ok := every.(lua.LNumber)
		if !ok || seconds <= 0 {
			L.ArgError(idx, "every must be a positive number of seconds")
			return
		}
		reset.Kind = RecordResetEvery
		reset.Every = time.Duration(float64(seconds) * float64(time.Second))
	}
}

func (reset *RecordReset) parseClock(L *lua.LState, idx int, text string) time.Duration {
	t, err := time.Parse("15:04", text)
	if err != nil {
		L.ArgError(idx, "reset time must be HH:MM")
		return 0
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute
}
