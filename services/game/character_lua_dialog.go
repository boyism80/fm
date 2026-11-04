package game

import (
	"log"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	"github.com/boyism80/fm/maple/wz"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaDialogFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"open_shop": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			npcID := uint32(L.CheckInt(2))
			stock := L.CheckTable(3)
			if ch.dialog.Busy() {
				L.Push(lua.LFalse)
				return 1
			}
			shop := &wz.Shop{NpcID: npcID}
			stock.ForEach(func(_, entry lua.LValue) {
				row, ok := entry.(*lua.LTable)
				if !ok {
					return
				}
				shop.Items = append(shop.Items, wz.ShopItem{
					ItemID: uint32(lua.LVAsNumber(row.RawGetString("item"))),
					Price:  int(lua.LVAsNumber(row.RawGetString("price"))),
				})
			})
			ch.dialog.OpenShop(shop)
			L.Push(lua.LTrue)
			return 1
		},
		"dialog": func(L *lua.LState) int {
			argc := L.GetTop()
			if argc == 1 {
				switch L.Get(1).Type() {
				case lua.LTBool:
					b := lua.LVAsBool(L.Get(1))
					L.SetTop(0)
					L.Push(lua.LBool(b))
					return 1
				case lua.LTNil:
					L.SetTop(0)
					L.Push(lua.LBool(false))
					return 1
				}
			}
			ch := luax.Check[Character](L, 1)
			npc := 0
			message := ""
			prev := false
			next := false
			switch argc {
			case 5:
				next = L.CheckBool(5)
				fallthrough
			case 4:
				prev = L.CheckBool(4)
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}
			ch.listener.OnDialog(ch, uint32(npc), message, prev, next)
			ch.dialog.Ask(L, constant.DialogTypeDefault, 0)
			return L.Yield()
		},
		"dialog_yes_no": func(L *lua.LState) int {
			argc := L.GetTop()
			if argc == 1 {
				switch L.Get(1).Type() {
				case lua.LTBool:
					b := lua.LVAsBool(L.Get(1))
					L.SetTop(0)
					L.Push(lua.LBool(b))
					return 1
				case lua.LTNil:
					L.SetTop(0)
					L.Push(lua.LBool(false))
					return 1
				}
			}
			ch := luax.Check[Character](L, 1)
			npc := 0
			message := ""
			prev := false
			next := false
			switch argc {
			case 5:
				next = L.CheckBool(5)
				fallthrough
			case 4:
				prev = L.CheckBool(4)
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}
			ch.listener.OnDialogYesNo(ch, uint32(npc), message, prev, next)
			ch.dialog.Ask(L, constant.DialogTypeYesNo, 0)
			return L.Yield()
		},
		"dialog_list": func(L *lua.LState) int {
			argc := L.GetTop()
			ch := luax.Check[Character](L, 1)

			npc := 0
			message := ""
			selections := []string{}
			switch argc {
			case 4:
				tbl := L.CheckTable(4)
				tbl.ForEach(func(_, value lua.LValue) {
					if str, ok := value.(lua.LString); ok {
						selections = append(selections, string(str))
					}
				})
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}

			ch.listener.OnDialogList(ch, uint32(npc), message, selections)
			ch.dialog.Ask(L, constant.DialogTypeList, len(selections))
			return L.Yield(lua.LNumber(0))
		},
		"dialog_style": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 4 {
				L.ArgError(2, "dialog_style(npc, text, styles) requires npc, text, and style ids")
				return 0
			}
			npc := 0
			switch npcArg := L.Get(2).(type) {
			case lua.LNumber:
				npc = int(npcArg)
			case *lua.LUserData:
				if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
					npc = int(n.wz.BaseSpawn.ID)
				}
			}
			styles := []uint32{}
			L.CheckTable(4).ForEach(func(_, value lua.LValue) {
				if n, ok := value.(lua.LNumber); ok {
					styles = append(styles, uint32(n))
				}
			})
			ch.listener.OnDialogStyle(ch, uint32(npc), L.CheckString(3), styles)
			ch.dialog.Ask(L, constant.DialogTypeStyle, len(styles))
			return L.Yield(lua.LNumber(0))
		},
		"dialog_accept": func(L *lua.LState) int {
			argc := L.GetTop()
			if argc == 1 {
				switch L.Get(1).Type() {
				case lua.LTBool:
					b := lua.LVAsBool(L.Get(1))
					L.SetTop(0)
					L.Push(lua.LBool(b))
					return 1
				case lua.LTNil:
					L.SetTop(0)
					L.Push(lua.LBool(false))
					return 1
				}
			}
			ch := luax.Check[Character](L, 1)
			npc := 0
			message := ""
			enableEscape := false
			switch argc {
			case 4:
				enableEscape = L.CheckBool(4)
				fallthrough
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}
			ch.listener.OnDialogAccept(ch, uint32(npc), message, enableEscape)
			if enableEscape {
				ch.dialog.Ask(L, constant.DialogTypeAcceptEscape, 0)
			} else {
				ch.dialog.Ask(L, constant.DialogTypeAccept, 0)
			}
			return L.Yield()
		},
		"dialog_input": func(L *lua.LState) int {
			argc := L.GetTop()
			ch := luax.Check[Character](L, 1)

			npc := 0
			message := ""
			switch argc {
			case 3:
				message = L.CheckString(3)
				fallthrough
			case 2:
				switch npcArg := L.Get(2).(type) {
				case lua.LNumber:
					npc = int(npcArg)
				case *lua.LUserData:
					if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
						npc = int(n.wz.BaseSpawn.ID)
					}
				}
			}

			ch.listener.OnDialogInput(ch, uint32(npc), message)
			ch.dialog.Ask(L, constant.DialogTypeInput, 0)
			return L.Yield(lua.LNumber(0))
		},
		"dialog_number": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			npc := 0
			switch npcArg := L.Get(2).(type) {
			case lua.LNumber:
				npc = int(npcArg)
			case *lua.LUserData:
				if n, ok := npcArg.Value.(*Npc); ok && n.wz != nil && n.wz.BaseSpawn != nil {
					npc = int(n.wz.BaseSpawn.ID)
				}
			}
			message := L.CheckString(3)
			def := int32(L.CheckInt(4))
			min := int32(L.CheckInt(5))
			max := int32(L.CheckInt(6))
			if min > max {
				L.ArgError(5, "min must not exceed max")
				return 0
			}

			ch.listener.OnDialogNumber(ch, uint32(npc), message, def, min, max)
			ch.dialog.AskNumber(L, min, max)
			return L.Yield(lua.LNumber(0))
		},
		"cheat_report": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LString(ch.cheat.Report()))
			return 1
		},
		"open_storage": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			npcID := uint32(L.CheckInt(2))
			storeFee := int32(L.OptInt(3, 100))
			takeOutFee := int32(L.OptInt(4, 0))
			ch.account.storage.Open(npcID, storeFee, takeOutFee)
			return 0
		},
		"open_store_bank": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			npcID := uint32(L.CheckInt(2))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.storeBank.Open(cfg.ActorContext, npcID)
			return 0
		},
		"open_duey": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			ch.duey.Open(cfg.ActorContext, false)
			return 0
		},
		"open_quick_parcel": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.duey.OpenQuick() == nil))
			return 1
		},
		"send_parcel": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			recipient := L.CheckString(2)
			senderName := L.CheckString(3)
			itemID := uint32(L.OptInt(4, 0))
			count := uint16(L.OptInt(5, 1))
			meso := int32(L.OptInt(6, 0))
			message := L.OptString(7, "")
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.ArgError(1, "actor context not available")
				return 0
			}
			L.Push(lua.LBool(ch.duey.SendFromSystem(cfg.ActorContext, recipient, senderName, itemID, count, meso, message) == nil))
			return 1
		},
		"lie_detector_answer": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			m := ch.GetMap()
			if m == nil {
				L.Push(lua.LNil)
				return 1
			}
			target := m.FindPlayerByName(L.CheckString(2))
			if target == nil || target.lieDetector.testerID == 0 {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(lua.LString(target.lieDetector.answer))
			return 1
		},
		"open_npc": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 2 {
				L.ArgError(2, "open_npc(npc) requires npc")
				return 0
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.RaiseError("open_npc: thread has no actor context")
				return 0
			}
			var npc *Npc
			if ud, ok := L.Get(2).(*lua.LUserData); ok {
				npc, _ = ud.Value.(*Npc)
			}
			if npc == nil {
				npcID, ok := LuaCheckNpcID(L, 2)
				if !ok {
					return 0
				}
				npc = ch.GetMap().NpcByTemplate(npcID)
				if npc == nil {
					npc = &Npc{wz: &wz.NpcSpawn{BaseSpawn: &wz.BaseSpawn{ID: npcID}}}
				}
			}
			ch.dialog.Close()
			if err := ch.dialog.Open(cfg.ActorContext, npc); err != nil {
				L.RaiseError("open_npc: %v", err)
				return 0
			}
			return 0
		},
		"script": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			if argc < 3 {
				L.ArgError(2, "script(scriptPath, funcName, ...) requires at least 2 arguments")
				return 0
			}
			scriptPath := L.CheckString(2)
			funcName := L.CheckString(3)
			args := make([]interface{}, 0, argc-2)
			args = append(args, ch)
			for i := 4; i <= argc; i++ {
				if v, ok := luax.LValueToInterface(L.Get(i)); ok {
					args = append(args, v)
				}
			}
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.RaiseError("script: thread has no actor PID (call from command context)")
				return 0
			}
			pid := cfg.ActorContext.Self()
			if pid == nil {
				L.RaiseError("script: thread has no actor PID (call from command context)")
				return 0
			}
			mapInstance := ch.GetMap()
			if mapInstance == nil {
				L.RaiseError("script: character map not found")
				return 0
			}
			root := mapInstance.GetLuaRoot()
			if root == nil {
				L.RaiseError("script: root lua state not found")
				return 0
			}
			thread, err := luax.NewThread(root, scriptPath)
			if err != nil {
				L.RaiseError("script: %v", err)
				return 0
			}
			luax.SetConfiguration(thread, luax.Configuration{
				ActorContext: cfg.ActorContext,
				ActorPID:     mapInstance.LogicActorPID(),
			})
			p := luax.CallAsync(cfg.ActorContext, root, thread, funcName, args...)
			if !p.Completed() {
				return 0
			}
			vals, err := p.Result()
			if err != nil {
				L.RaiseError("script: %v", err)
				return 0
			}
			luax.Close(thread)
			if len(vals) == 0 {
				return 0
			}
			L.Push(vals[0])
			return 1
		},
		"run_quest_hook": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 3 {
				L.ArgError(2, "run_quest_hook(quest_id, hook) requires quest id and hook")
				return 0
			}
			questID := uint32(L.CheckInt(2))
			hook := L.CheckString(3)
			ch.quests.RunHook(questID, hook)
			return 0
		},
		"mktimer": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			if argc < 5 {
				L.ArgError(2, "mktimer(key, intervalMs, repeat, function [, ...args])")
				return 0
			}
			key := L.CheckString(2)
			intervalMs := L.CheckInt(3)
			if intervalMs < 1 {
				L.ArgError(3, "intervalMs must be positive")
				return 0
			}
			repeat := L.CheckBool(4)
			detached, err := luax.Detach(L, L.CheckFunction(5))
			if err != nil {
				L.RaiseError("mktimer: %v", err)
				return 0
			}
			args := luaValuesToInterfaces(L, 6, argc)
			// The timer follows the character across maps; each map runs its own copy of fn.
			var fnRoot *lua.LState
			var fn *lua.LFunction
			callback := func() {
				m := ch.GetMap()
				if m == nil {
					return
				}
				root := m.GetLuaRoot()
				if root == nil {
					return
				}
				if root != fnRoot {
					attached, err := detached.Attach(root)
					if err != nil {
						log.Printf("RunObjectTimer %s: %v", key, err)
						return
					}
					fnRoot = root
					fn = attached
				}
				fullArgs := make([]interface{}, 0, len(args)+1)
				fullArgs = append(fullArgs, ch)
				fullArgs = append(fullArgs, args...)
				if _, err := luax.CallFunction(root, fn, fullArgs...); err != nil {
					log.Printf("RunObjectTimer %s: %v", key, err)
				}
			}
			added := ch.AddTimer(key, time.Duration(intervalMs)*time.Millisecond, repeat, callback)
			L.Push(lua.LBool(added))
			return 1
		},
		"rmtimer": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			key := L.CheckString(2)
			removed := ch.RemoveTimer(key)
			L.Push(lua.LBool(removed))
			return 1
		},
	}
}
