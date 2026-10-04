package actor

import (
	"fmt"
	"log"
	"maps"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/boyism80/fm/core/luax"
	pconst "github.com/boyism80/fm/protocol/constant"
	"github.com/boyism80/fm/protocol/dto"
	"github.com/boyism80/fm/protocol/request"
	"github.com/boyism80/fm/protocol/response"
	"github.com/boyism80/fm/services/bot/bot"
	"github.com/boyism80/fm/services/bot/conn"
	"github.com/boyism80/fm/services/game/constant"
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
	&request.Attack{},
	&request.ItemLoot{},
	&request.MoveItem{},
	&request.DamageReactor{},
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
	a.wrapWaits(botIndex, "request", "request_on", "instance_move", "map_move", "warp", "transfer", "npc", "npc_click", "dialog", "kill", "loot", "drop", "hit_reactor")
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
		"report": func(L *lua.LState) int {
			text := L.CheckString(2)
			if a.report == nil {
				path := filepath.Join(a.cfg.ReportDir, a.suite.Name+".txt")
				if err := os.MkdirAll(a.cfg.ReportDir, 0o755); err != nil {
					L.RaiseError("report: %v", err)
					return 0
				}
				f, err := os.Create(path)
				if err != nil {
					L.RaiseError("report: %v", err)
					return 0
				}
				a.report = f
				log.Printf("[%s] report: %s", a.name, path)
			}
			if _, err := a.report.WriteString(text + "\n"); err != nil {
				L.RaiseError("report: %v", err)
			}
			return 0
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
		"level": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).Level))
			return 1
		},
		"job": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).Job))
			return 1
		},
		"exp": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).EXP))
			return 1
		},
		"meso": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).Meso))
			return 1
		},
		"fame": func(L *lua.LState) int {
			L.Push(lua.LNumber(a.checkBot(L).Fame))
			return 1
		},
		"items": func(L *lua.LState) int {
			counts := make(map[uint32]int)
			for _, tab := range a.checkBot(L).Items {
				for slot, item := range tab {
					if slot > 0 {
						counts[item.ItemID] += int(item.Count)
					}
				}
			}
			t := L.NewTable()
			for id, count := range counts {
				t.RawSetInt(int(id), lua.LNumber(count))
			}
			L.Push(t)
			return 1
		},
		"quests": func(L *lua.LState) int {
			t := L.NewTable()
			for id, status := range a.checkBot(L).Quests {
				t.RawSetInt(int(id), lua.LNumber(status))
			}
			L.Push(t)
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
		"areas": func(L *lua.LState) int {
			t := L.NewTable()
			if m, ok := a.wz.Maps[a.checkBot(L).Map]; ok {
				for _, r := range m.Areas {
					area := L.NewTable()
					area.RawSetString("x", lua.LNumber((int(r.Left)+int(r.Right))/2))
					area.RawSetString("y", lua.LNumber((int(r.Top)+int(r.Bottom))/2))
					t.Append(area)
				}
			}
			L.Push(t)
			return 1
		},
		"npc_position": func(L *lua.LState) int {
			templateID := uint32(L.CheckInt(2))
			if m, ok := a.wz.Maps[a.checkBot(L).Map]; ok {
				for _, spawn := range m.NpcSpawns {
					if spawn.BaseSpawn != nil && spawn.ID == templateID {
						L.Push(lua.LNumber(spawn.Position.X))
						L.Push(lua.LNumber(spawn.Position.Y))
						return 2
					}
				}
			}
			L.Push(lua.LNil)
			return 1
		},
		"mob_spots": func(L *lua.LState) int {
			t := L.NewTable()
			if m, ok := a.wz.Maps[a.checkBot(L).Map]; ok {
				for _, key := range slices.Sorted(maps.Keys(m.MobSpawns)) {
					spawn := m.MobSpawns[key]
					if spawn.BaseSpawn == nil {
						continue
					}
					spot := L.NewTable()
					spot.RawSetString("id", lua.LNumber(spawn.ID))
					spot.RawSetString("x", lua.LNumber(spawn.Position.X))
					spot.RawSetString("y", lua.LNumber(spawn.Position.Y))
					t.Append(spot)
				}
			}
			L.Push(t)
			return 1
		},
		"reactor_spots": func(L *lua.LState) int {
			t := L.NewTable()
			if m, ok := a.wz.Maps[a.checkBot(L).Map]; ok {
				for _, key := range slices.Sorted(maps.Keys(m.ReactorSpawns)) {
					spawn := m.ReactorSpawns[key]
					spot := L.NewTable()
					spot.RawSetString("id", lua.LNumber(spawn.ReactorID))
					spot.RawSetString("name", lua.LString(spawn.Name))
					spot.RawSetString("x", lua.LNumber(spawn.Position.X))
					spot.RawSetString("y", lua.LNumber(spawn.Position.Y))
					t.Append(spot)
				}
			}
			L.Push(t)
			return 1
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
		"mobs": func(L *lua.LState) int {
			templateID := uint32(L.OptInt(2, 0))
			t := L.NewTable()
			for oid, id := range a.checkBot(L).Mobs {
				if templateID != 0 && id != templateID {
					continue
				}
				mob := L.NewTable()
				mob.RawSetString("oid", lua.LNumber(oid))
				mob.RawSetString("id", lua.LNumber(id))
				t.Append(mob)
			}
			L.Push(t)
			return 1
		},
		"attack": func(L *lua.LState) int {
			b := a.checkBot(L)
			pkt, ok := a.attackPacket(b, uint32(L.CheckInt(2)), uint32(L.OptInt(3, 0)), uint8(L.OptInt(4, 1)))
			if ok == false {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(b.Send(pkt) == nil))
			return 1
		},
		"kill": func(L *lua.LState) int {
			b := a.checkBot(L)
			oid := uint32(L.CheckInt(2))
			pkt, ok := a.attackPacket(b, oid, 0, 1)
			if ok == false {
				L.Push(lua.LFalse)
				return 1
			}
			a.park(L)
			a.Request(L, b, b, pkt, a.timeout(), func(p any) bool {
				die, ok := p.(*response.DieMob)
				return ok && die.OID == oid
			}, func(_ any, ok bool) {
				a.wake(L, lua.LBool(ok))
			})
			return L.Yield()
		},
		"drops": func(L *lua.LState) int {
			itemID := uint32(L.OptInt(2, 0))
			t := L.NewTable()
			for oid, d := range a.checkBot(L).Drops {
				if itemID != 0 && d.ItemID != itemID {
					continue
				}
				drop := L.NewTable()
				drop.RawSetString("oid", lua.LNumber(oid))
				drop.RawSetString("item", lua.LNumber(d.ItemID))
				drop.RawSetString("meso", lua.LNumber(d.Meso))
				drop.RawSetString("owner", lua.LNumber(d.Owner))
				t.Append(drop)
			}
			L.Push(t)
			return 1
		},
		"loot": func(L *lua.LState) int {
			b := a.checkBot(L)
			oid := uint32(L.CheckInt(2))
			pos, _ := b.Position(a.wz)
			a.park(L)
			a.Request(L, b, b, &request.ItemLoot{Position: pos, OID: oid}, a.timeout(), func(p any) bool {
				remove, ok := p.(*response.RemoveItem)
				return ok && remove.OID == oid
			}, func(_ any, ok bool) {
				a.wake(L, lua.LBool(ok))
			})
			return L.Yield()
		},
		"drop": func(L *lua.LState) int {
			b := a.checkBot(L)
			itemID := uint32(L.CheckInt(2))
			count := uint16(L.OptInt(3, 1))
			var pkt *request.MoveItem
			for typ, tab := range b.Items {
				for slot, item := range tab {
					if slot > 0 && item.ItemID == itemID {
						pkt = &request.MoveItem{InventoryType: typ, Source: slot, Count: count}
					}
				}
			}
			if pkt == nil {
				L.Push(lua.LNil)
				return 1
			}
			a.park(L)
			a.Request(L, b, b, pkt, a.timeout(), func(p any) bool {
				spawn, ok := p.(*response.SpawnItem)
				return ok && spawn.OwnerID == b.CharID && spawn.ItemModel.GetID() == itemID
			}, func(p any, ok bool) {
				if ok == false {
					a.wake(L, lua.LNil)
					return
				}
				a.wake(L, lua.LNumber(p.(*response.SpawnItem).ID))
			})
			return L.Yield()
		},
		"reactors": func(L *lua.LState) int {
			templateID := uint32(L.OptInt(2, 0))
			t := L.NewTable()
			for oid, r := range a.checkBot(L).Reactors {
				if templateID != 0 && r.ReactorID != templateID {
					continue
				}
				reactor := L.NewTable()
				reactor.RawSetString("oid", lua.LNumber(oid))
				reactor.RawSetString("id", lua.LNumber(r.ReactorID))
				reactor.RawSetString("state", lua.LNumber(r.State))
				reactor.RawSetString("name", lua.LString(r.Name))
				reactor.RawSetString("x", lua.LNumber(r.Position.X))
				reactor.RawSetString("y", lua.LNumber(r.Position.Y))
				t.Append(reactor)
			}
			L.Push(t)
			return 1
		},
		"hit_reactor": func(L *lua.LState) int {
			b := a.checkBot(L)
			oid := uint32(L.CheckInt(2))
			pkt := &request.DamageReactor{OID: oid, HitSide: constant.ReactorHitGroundRight}
			a.park(L)
			a.Request(L, b, b, pkt, a.timeout(), func(p any) bool {
				switch p := p.(type) {
				case *response.TriggerReactor:
					return p.Reactor.OID == oid
				case *response.DestroyReactor:
					return p.Reactor.OID == oid
				default:
					return false
				}
			}, a.wakeWithPacket(L))
			return L.Yield()
		},
	}
}

func (a *SuiteActor) attackPacket(b *bot.Bot, oid uint32, damage uint32, hits uint8) (*request.Attack, bool) {
	templateID, ok := b.Mobs[oid]
	if ok == false {
		return nil, false
	}
	if damage == 0 {
		damage = 1
		if mob, ok := a.wz.Monsters[templateID]; ok && mob.MaxHP > 0 {
			damage = uint32(mob.MaxHP)
		}
	}
	pairs := make([]dto.DamagePair, hits)
	for i := range pairs {
		pairs[i] = dto.DamagePair{Damage: damage}
	}
	pos, _ := b.Position(a.wz)
	return &request.Attack{CloseAttackInfo: dto.CloseAttackInfo{
		AttackHeader: dto.AttackHeader{Targets: 1, Hits: hits},
		Damages:      []dto.AttackPair{{OID: oid, DamagePairs: pairs}},
		Position:     pos,
	}}, true
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
