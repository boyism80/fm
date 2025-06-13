package builtin

import (
	"time"

	"github.com/boyism80/fm/common/luax"
	"github.com/boyism80/fm/game/msg"
	lua "github.com/yuin/gopher-lua"
)

func Sleep(L *lua.LState) int {
	ms := L.CheckInt(1)
	luax.SendAfter(time.Duration(ms)*time.Millisecond, &msg.LuaResume{
		Lua: L,
	})
	return L.Yield(lua.LNumber(0))
}
