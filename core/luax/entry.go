package luax

import (
	"strings"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

const entrySource = `
local function pack(...)
	return select("#", ...), { ... }
end
local fn = ...
local n, results = pack(fn(select(2, ...)))
return unpack(results, 1, n)
`

var entryProto = compileEntry()

func compileEntry() *lua.FunctionProto {
	chunk, err := parse.Parse(strings.NewReader(entrySource), "entry")
	if err != nil {
		panic(err)
	}
	proto, err := lua.Compile(chunk, "entry")
	if err != nil {
		panic(err)
	}
	return proto
}

func NewEntry(L *lua.LState, fn *lua.LFunction, args []lua.LValue) (*lua.LFunction, []lua.LValue) {
	return L.NewFunctionFromProto(entryProto), append([]lua.LValue{fn}, args...)
}
