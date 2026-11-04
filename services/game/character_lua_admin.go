package game

import (
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaAdminFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"disconnect": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.listener.Disconnect(ch)
			return 0
		},
		"invincible": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.GetInvincible()))
				return 1
			case 2:
				ch.SetInvincible(L.CheckBool(2))
				return 0
			default:
				L.ArgError(2, "invincible() requires 0 or 1 arguments")
				return 0
			}
		},
		"dojo_energy": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			switch L.GetTop() {
			case 1:
				L.Push(lua.LNumber(ch.dojo.energy))
				return 1
			case 2:
				ch.dojo.SetEnergy(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "dojo_energy() requires 0 or 1 arguments")
				return 0
			}
		},
		"instant_kill": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.gm.instantKill))
				return 1
			case 2:
				ch.gm.instantKill = L.CheckBool(2)
				return 0
			default:
				L.ArgError(2, "instant_kill() requires 0 or 1 arguments")
				return 0
			}
		},
		"player_mode": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.gm.playerMode))
				return 1
			case 2:
				ch.gm.playerMode = L.CheckBool(2)
				return 0
			default:
				L.ArgError(2, "player_mode() requires 0 or 1 arguments")
				return 0
			}
		},
		"timer_limit": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.gm.timerLimit))
				return 1
			case 2:
				ch.gm.timerLimit = uint32(L.CheckNumber(2))
				return 0
			default:
				L.ArgError(2, "timer_limit() requires 0 or 1 arguments")
				return 0
			}
		},
	}
}
