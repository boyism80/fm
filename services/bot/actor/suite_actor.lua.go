package actor

import (
	"fmt"
	"log"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/bot"
	"github.com/boyism80/fm/services/bot/conn"
	"github.com/boyism80/fm/types"
	lua "github.com/yuin/gopher-lua"
)

var requests = []outbound{
	&request.NormalChat{},
	&request.NpcClick{},
	&request.Dialog{},
	&request.Warp{},
	&request.PartyOperation{},
	&request.SwitchChannel{},
	&request.QuestAction{},
}

func (a *SuiteActor) register() {
	L := a.L
	L.SetGlobal("test_suite", L.NewFunction(func(L *lua.LState) int {
		a.def = L.CheckTable(1)
		return 0
	}))
	L.SetGlobal("log", L.NewFunction(func(L *lua.LState) int {
		log.Printf("[%s] %s: %s", a.name, L.CheckString(1), L.CheckString(2))
		return 0
	}))

	resp := L.NewTable()
	for _, pkt := range conn.Responses {
		name := a.packetName(pkt)
		resp.RawSetString(name, lua.LString(name))
	}
	L.SetGlobal("resp", resp)

	req := L.NewTable()
	for _, pkt := range requests {
		t := reflect.TypeOf(pkt).Elem()
		req.RawSetString(a.marshal.Name(t.Name()), L.NewFunction(func(L *lua.LState) int {
			out := reflect.New(t).Interface()
			if err := a.marshal.FromLua(L.OptTable(1, L.NewTable()), out); err != nil {
				L.RaiseError("req.%s: %v", a.marshal.Name(t.Name()), err)
				return 0
			}
			L.Push(a.newUserData(out, "bot_request"))
			return 1
		}))
	}
	L.SetGlobal("req", req)
	L.NewTypeMetatable("bot_request")

	party := L.NewTable()
	party.RawSetString("Create", lua.LNumber(pconst.PartyC2SCreate))
	party.RawSetString("Leave", lua.LNumber(pconst.PartyC2SLeave))
	party.RawSetString("AcceptInvite", lua.LNumber(pconst.PartyC2SAcceptInvite))
	party.RawSetString("Invite", lua.LNumber(pconst.PartyC2SInvite))
	party.RawSetString("Expel", lua.LNumber(pconst.PartyC2SExpel))
	party.RawSetString("ChangeLeader", lua.LNumber(pconst.PartyC2SChangeLeader))
	L.SetGlobal("PARTY", party)
	L.SetGlobal("wz", L.SetFuncs(L.NewTable(), a.wzFuncs()))

	ctxIndex := L.SetFuncs(L.NewTable(), a.ctxFuncs())
	botIndex := L.SetFuncs(L.NewTable(), a.botFuncs())
	a.wrapWaits(ctxIndex, "sleep")
	a.wrapWaits(botIndex, "request", "request_on", "instance_move", "map_move", "warp", "transfer", "npc", "npc_click", "dialog")
	L.SetField(L.NewTypeMetatable("bot_ctx"), "__index", ctxIndex)
	L.SetField(L.NewTypeMetatable("bot"), "__index", botIndex)
	a.ctxUD = a.newUserData(a, "bot_ctx")
}

func (a *SuiteActor) wrapWaits(index *lua.LTable, names ...string) {
	if err := a.L.DoString("return function(raw) return function(...) local r, name = raw(...) return r, name end end"); err != nil {
		panic(err)
	}
	wrap := a.L.Get(-1).(*lua.LFunction)
	a.L.Pop(1)
	for _, name := range names {
		wrapped, err := luax.CallFunction(a.L, wrap, index.RawGetString(name))
		if err != nil {
			panic(err)
		}
		index.RawSetString(name, wrapped)
	}
}

func (a *SuiteActor) newUserData(v any, typeName string) *lua.LUserData {
	ud := a.L.NewUserData()
	ud.Value = v
	a.L.SetMetatable(ud, a.L.GetTypeMetatable(typeName))
	return ud
}

func (a *SuiteActor) ctxFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"bot": func(L *lua.LState) int {
			i := L.CheckInt(2)
			if i < 0 || i >= len(a.bots) {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(a.botUDs[a.bots[i]])
			return 1
		},
		"bot_count": func(L *lua.LState) int {
			L.Push(lua.LNumber(len(a.bots)))
			return 1
		},
		"seat": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.seat))
			return 1
		},
		"fail": func(L *lua.LState) int {
			a.Fail(L.CheckString(2))
			L.Push(lua.LFalse)
			return 1
		},
		"sleep": func(L *lua.LState) int {
			a.park(L)
			a.Sleep(L, time.Duration(L.CheckInt(2))*time.Millisecond)
			return L.Yield()
		},
		"hook": func(L *lua.LState) int {
			name := L.CheckString(2)
			fn := L.CheckFunction(3)
			a.hooks[name] = func(b *bot.Bot, pkt any) {
				if _, err := luax.CallFunction(a.L, fn, a.ctxUD, a.botUDs[b], a.marshal.ToLua(a.L, pkt)); err != nil {
					log.Printf("[%s] hook %s: %v", a.name, name, err)
				}
			}
			return 0
		},
		"unhook": func(L *lua.LState) int {
			delete(a.hooks, L.CheckString(2))
			return 0
		},
	}
}

func (a *SuiteActor) botFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).CharID))
			return 1
		},
		"name": func(L *lua.LState) int {
			L.Push(lua.LString(a.checkBot(L).Name))
			return 1
		},
		"map": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).Map))
			return 1
		},
		"hp": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).HP))
			return 1
		},
		"position": func(L *lua.LState) int {
			pos, ok := a.checkBot(L).Position(a.wz)
			if ok == false {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LNumber(pos.X))
			L.Push(lua.LNumber(pos.Y))
			return 2
		},
		"move": func(L *lua.LState) int {
			b := a.checkBot(L)
			pos := types.Point[int16]{X: int16(L.CheckInt(2)), Y: int16(L.CheckInt(3))}
			L.Push(lua.LBool(b.Move(pos, int16(L.OptInt(4, 0))) == nil))
			return 1
		},
		"send": func(L *lua.LState) int {
			b := a.checkBot(L)
			L.Push(lua.LBool(b.Send(a.checkRequest(L, 2)) == nil))
			return 1
		},
		"chat": func(L *lua.LState) int {
			b := a.checkBot(L)
			L.Push(lua.LBool(a.Command(b, L.CheckString(2)) == nil))
			return 1
		},
		"command": func(L *lua.LState) int {
			b := a.checkBot(L)
			text := L.CheckString(2)
			if strings.HasPrefix(text, "/") == false {
				text = "/" + text
			}
			L.Push(lua.LBool(a.Command(b, text) == nil))
			return 1
		},
		"request": func(L *lua.LState) int {
			b := a.checkBot(L)
			names := a.checkNames(L, 2)
			var pkt outbound
			if L.Get(3) != lua.LNil {
				pkt = a.checkRequest(L, 3)
			}
			pred := L.OptFunction(4, nil)
			timeout := time.Duration(L.OptInt(5, a.cfg.TimeoutMs)) * time.Millisecond
			a.park(L)
			a.Request(L, b, b, pkt, timeout, func(p any) bool {
				return names[a.packetName(p)] && (pred == nil || a.callPredicate(pred, p))
			}, a.wakeWithPacket(L))
			return L.Yield()
		},
		"request_on": func(L *lua.LState) int {
			sender := a.checkBot(L)
			listener, ok := L.CheckUserData(2).Value.(*bot.Bot)
			if ok == false {
				L.ArgError(2, "bot expected")
			}
			names := a.checkNames(L, 3)
			var pkt outbound
			if L.Get(4) != lua.LNil {
				pkt = a.checkRequest(L, 4)
			}
			pred := L.OptFunction(5, nil)
			timeout := time.Duration(L.OptInt(6, a.cfg.TimeoutMs)) * time.Millisecond
			a.park(L)
			a.Request(L, sender, listener, pkt, timeout, func(p any) bool {
				return names[a.packetName(p)] && (pred == nil || a.callPredicate(pred, p))
			}, a.wakeWithPacket(L))
			return L.Yield()
		},
		"warp": func(L *lua.LState) int {
			b := a.checkBot(L)
			pkt := &request.Warp{Target: 0xFFFFFFFF, PortalName: L.CheckString(2)}
			a.park(L)
			a.Request(L, b, b, pkt, a.timeout(), func(p any) bool {
				_, ok := p.(*response.Warp)
				return ok
			}, a.wakeWithPacket(L))
			return L.Yield()
		},
		"transfer": func(L *lua.LState) int {
			b := a.checkBot(L)
			channel := L.CheckInt(2)
			if channel < 0 || channel > 255 {
				L.ArgError(2, "channel must fit in one byte")
			}
			i := slices.Index(a.bots, b)
			a.park(L)
			a.switchChannel(L, i, uint8(channel), func(ok bool) {
				a.wake(L, lua.LBool(ok))
			})
			return L.Yield()
		},
		"instance_move": func(L *lua.LState) int {
			b := a.checkBot(L)
			mapID := uint32(L.CheckInt(2))
			text := fmt.Sprintf("/인스턴스이동 %d %d", mapID, a.seat)
			if L.GetTop() >= 3 {
				text = fmt.Sprintf("%s %d", text, L.CheckInt(3))
			}
			a.park(L)
			a.warpByCommand(L, b, text, mapID, true, func(ok bool) {
				a.wake(L, lua.LBool(ok))
			})
			return L.Yield()
		},
		"map_move": func(L *lua.LState) int {
			b := a.checkBot(L)
			mapID := uint32(L.CheckInt(2))
			if b.Map == mapID {
				L.Push(lua.LTrue)
				return 1
			}
			text := fmt.Sprintf("/맵이동 %d", mapID)
			if L.GetTop() >= 3 {
				text = fmt.Sprintf("%s %d", text, L.CheckInt(3))
			}
			a.park(L)
			a.warpByCommand(L, b, text, mapID, true, func(ok bool) {
				a.wake(L, lua.LBool(ok))
			})
			return L.Yield()
		},
		"npc": func(L *lua.LState) int {
			b := a.checkBot(L)
			templateID := uint32(L.CheckInt(2))
			if oid, ok := b.FindNPC(templateID); ok {
				L.Push(lua.LNumber(oid))
				return 1
			}
			timeout := time.Duration(L.OptInt(3, a.cfg.TimeoutMs)) * time.Millisecond
			a.park(L)
			a.Request(L, b, b, nil, timeout, func(p any) bool {
				spawn, ok := p.(*response.SpawnNpc)
				return ok && spawn.NPC.NpcId == templateID
			}, func(p any, ok bool) {
				if ok == false {
					a.wake(L, lua.LNil)
					return
				}
				a.wake(L, lua.LNumber(p.(*response.SpawnNpc).NPC.OID))
			})
			return L.Yield()
		},
		"npc_click": func(L *lua.LState) int {
			b := a.checkBot(L)
			oid := uint32(L.CheckInt(2))
			a.park(L)
			a.Request(L, b, b, &request.NpcClick{OID: oid}, a.timeout(), a.dialogPacket, a.wakeWithPacket(L))
			return L.Yield()
		},
		"dialog": func(L *lua.LState) int {
			b := a.checkBot(L)
			pkt := &request.Dialog{
				DialogType: b.Dialog,
				Next:       L.ToBool(2),
				Selected:   uint32(L.OptInt(3, 0)),
				Text:       L.OptString(4, ""),
			}
			if pkt.Next == false {
				if err := b.Send(pkt); err != nil {
					a.Fail(fmt.Sprintf("%s send: %v", b.Name, err))
				}
				L.Push(lua.LNil)
				return 1
			}
			a.park(L)
			a.Request(L, b, b, pkt, a.timeout(), a.dialogPacket, func(p any, ok bool) {
				if ok == false {
					a.wake(L, lua.LNil)
					return
				}
				a.wake(L, a.marshal.ToLua(a.L, p))
			})
			return L.Yield()
		},
	}
}

func (a *SuiteActor) checkBot(L *lua.LState) *bot.Bot {
	b, ok := L.CheckUserData(1).Value.(*bot.Bot)
	if ok == false {
		L.ArgError(1, "bot expected")
	}
	return b
}

func (a *SuiteActor) checkRequest(L *lua.LState, n int) outbound {
	pkt, ok := L.CheckUserData(n).Value.(outbound)
	if ok == false {
		L.ArgError(n, "req.* packet expected")
	}
	return pkt
}

func (a *SuiteActor) checkNames(L *lua.LState, n int) map[string]bool {
	names := make(map[string]bool)
	switch v := L.Get(n).(type) {
	case lua.LString:
		names[string(v)] = true
	case *lua.LTable:
		v.ForEach(func(_, name lua.LValue) {
			names[name.String()] = true
		})
	default:
		L.ArgError(n, "resp.* name or a list of names expected")
	}
	return names
}

func (a *SuiteActor) packetName(pkt any) string {
	return a.marshal.Name(reflect.TypeOf(pkt).Elem().Name())
}

func (a *SuiteActor) dialogPacket(pkt any) bool {
	switch pkt.(type) {
	case *response.Dialog, *response.DialogYesNo, *response.DialogInput, *response.DialogList, *response.DialogStyle, *response.DialogAccept:
		return true
	default:
		return false
	}
}

func (a *SuiteActor) callPredicate(pred *lua.LFunction, pkt any) bool {
	ret, err := luax.CallFunction(a.L, pred, a.marshal.ToLua(a.L, pkt), lua.LString(a.packetName(pkt)))
	if err != nil {
		log.Printf("[%s] predicate: %v", a.name, err)
		return false
	}
	return lua.LVAsBool(ret)
}

func (a *SuiteActor) park(L *lua.LState) {
	if _, ok := a.threads[L]; ok == false {
		L.RaiseError("waiting builtins only run in a scenario coroutine, not in hooks or predicates")
	}
}

func (a *SuiteActor) spawn(fn *lua.LFunction, done func(values []lua.LValue, err error), args ...lua.LValue) {
	co, _ := a.L.NewThread()
	a.threads[co] = done
	a.resume(co, fn, args...)
}

func (a *SuiteActor) wake(co *lua.LState, values ...lua.LValue) {
	if a.finished {
		return
	}
	if _, ok := a.threads[co]; ok == false {
		return
	}
	a.resume(co, nil, values...)
}

func (a *SuiteActor) wakeWithPacket(co *lua.LState) func(pkt any, ok bool) {
	return func(pkt any, ok bool) {
		if ok == false {
			a.wake(co, lua.LFalse)
			return
		}
		a.wake(co, a.marshal.ToLua(a.L, pkt), lua.LString(a.packetName(pkt)))
	}
}

func (a *SuiteActor) resume(co *lua.LState, fn *lua.LFunction, args ...lua.LValue) {
	state, err, values := a.L.Resume(co, fn, args...)
	if err == nil && state == lua.ResumeYield {
		return
	}
	done := a.threads[co]
	delete(a.threads, co)

	pending := a.sleeping[co]
	delete(a.sleeping, co)
	var left []*waiter
	for _, w := range a.waiters {
		if w.thread == co {
			left = append(left, w)
		}
	}
	for _, w := range left {
		a.removeWaiter(w)
	}
	if err == nil && (pending || len(left) > 0) {
		err = fmt.Errorf("a waiting builtin ran inside pcall; call it outside pcall")
	}
	if a.finished {
		return
	}
	done(values, err)
}
