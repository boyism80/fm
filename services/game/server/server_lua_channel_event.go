package server

import (
	"log"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/services/game"
	lua "github.com/yuin/gopher-lua"
)

func registerChannelEventLuaFuncs(gs *GameServer, luaState *lua.LState) {
	luax.RegisterFunc(luaState, "ox_questions", func(L *lua.LState) int {
		tbl := L.NewTable()
		for _, q := range gs.GetResources().OXQuestions {
			row := L.NewTable()
			row.RawSetString("set", lua.LNumber(q.Set))
			row.RawSetString("id", lua.LNumber(q.ID))
			row.RawSetString("o", lua.LBool(q.O))
			tbl.Append(row)
		}
		L.Push(tbl)
		return 1
	})

	// The event's timers have no character to speak through, so this notice is not tied to one and does not wait.
	luax.RegisterFunc(luaState, "world_notice", func(L *lua.LState) int {
		text := L.CheckString(1)
		msgType := constant.MsgNotice
		if L.GetTop() >= 2 {
			msgType = constant.ServerMessageType(L.CheckInt(2))
		}
		cfg, ok := luax.GetConfiguration(L)
		if !ok || cfg.ActorContext == nil {
			return 0
		}
		gs.systems.Get[*game.InternalSystem]().BroadcastNoticeAsync(cfg.ActorContext, nil, msgType, text, false).OnError(func(err error) {
			log.Printf("world notice: %v", err)
		})
		return 0
	})
}
