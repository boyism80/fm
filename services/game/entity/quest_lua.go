package entity

import (
	"fmt"
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/luax"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) LuaQuest(questID uint32) *Quest {
	if ch == nil || ch.Quests == nil {
		return nil
	}
	if qp := ch.Quests.Get(questID); qp != nil {
		return qp
	}
	def := ch.Quests.wzDef(questID)
	return &Quest{
		container: ch.Quests,
		Wz:        def,
		QuestID:   questID,
		Status:    QuestStatusNotStarted,
	}
}

func (qp *Quest) LuaTypeName() string {
	return "LuaQuest"
}

func (qp *Quest) String() string {
	if qp == nil {
		return "LuaQuest(nil)"
	}
	return fmt.Sprintf("LuaQuest(%d)", qp.QuestID)
}

func (qp *Quest) Type() lua.LValueType {
	return lua.LTUserData
}

func LuaCheckQuest(L *lua.LState, idx int) (*Quest, bool) {
	ud := L.CheckUserData(idx)
	qp, ok := ud.Value.(*Quest)
	if !ok || qp == nil {
		L.ArgError(idx, "LuaQuest expected")
		return nil, false
	}
	return qp, true
}

func (qp *Quest) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "id() is read-only")
				return 0
			}
			L.Push(lua.LNumber(q.QuestID))
			return 1
		},
		"status": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "status() is read-only")
				return 0
			}
			L.Push(lua.LNumber(q.Status))
			return 1
		},
		"record": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() == 1 {
				L.Push(lua.LString(q.StatusRecord.AsString()))
				return 1
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "record() takes 0 or 1 arguments")
				return 0
			}
			if q.container == nil || q.container.Get(q.QuestID) == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			q.StatusRecord.WriteString(L.CheckString(2))
			L.Push(lua.LBool(true))
			return 1
		},
		"deadline": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() == 1 {
				if q.Deadline.IsZero() {
					L.Push(lua.LNil)
				} else {
					L.Push(lua.LNumber(q.Deadline.UnixMilli()))
				}
				return 1
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "deadline() takes 0 or 1 arguments")
				return 0
			}
			if q.container == nil || q.container.Get(q.QuestID) == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if L.Get(2) == lua.LNil {
				q.ResetDeadline()
			} else {
				q.SetDeadline(time.UnixMilli(L.CheckInt64(2)))
			}
			L.Push(lua.LBool(true))
			return 1
		},
		"can_start": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			npcID := uint32(0)
			opts := QuestPhaseOpts{}
			if L.GetTop() >= 2 {
				npcID = uint32(L.CheckInt(2))
				opts.NpcID = &npcID
			}
			if q.container == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			L.Push(lua.LBool(q.container.CanStart(q.QuestID, opts) == nil))
			return 1
		},
		"can_complete": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "can_complete() takes no arguments")
				return 0
			}
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			L.Push(lua.LBool(q.CanComplete(q.container.owner, QuestPhaseOpts{}) == nil))
			return 1
		},
		"started": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "started() takes no arguments")
				return 0
			}
			L.Push(lua.LBool(q.IsStarted()))
			return 1
		},
		"completed": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "completed() takes no arguments")
				return 0
			}
			L.Push(lua.LBool(q != nil && q.Status == QuestStatusCompleted))
			return 1
		},
		"wz": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "wz() takes no arguments")
				return 0
			}
			if q.Wz == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(q.Wz.ToLuaTable(L))
			return 1
		},
		"mob_kills": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "mob_kills(mob_id) requires mob id")
				return 0
			}
			mobID := uint32(L.CheckInt(2))
			kills := 0
			if q.MobKills != nil {
				kills = q.MobKills[mobID]
			}
			L.Push(lua.LNumber(kills))
			return 1
		},
		"set_mob_kills": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 3 {
				L.ArgError(2, "set_mob_kills(mob_id, count) requires mob id and count")
				return 0
			}
			if q.container == nil || q.container.owner == nil || q.container.Get(q.QuestID) == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if !q.IsStarted() {
				L.Push(lua.LBool(false))
				return 1
			}
			mobID := uint32(L.CheckInt(2))
			count := int(L.CheckInt(3))
			if count < 0 {
				count = 0
			}
			if q.MobKills == nil {
				q.MobKills = make(map[uint32]int)
			}
			q.MobKills[mobID] = count
			L.Push(lua.LBool(true))
			return 1
		},
		"sync_progress": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "sync_progress() takes no arguments")
				return 0
			}
			if q.container == nil || q.container.owner == nil || q.container.Get(q.QuestID) == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if q.container.owner.Listener != nil {
				q.container.owner.Listener.OnQuestProgress(q.container.owner, q)
			}
			L.Push(lua.LBool(true))
			return 1
		},
		"record_ex": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 2:
				key := L.CheckString(2)
				if q.container == nil {
					L.Push(lua.LNil)
					return 1
				}
				value, found := q.RecordExField(key)
				if !found {
					L.Push(lua.LNil)
				} else {
					L.Push(lua.LString(value))
				}
				return 1
			case 3:
				key := L.CheckString(2)
				value := L.CheckString(3)
				if q.container == nil {
					L.Push(lua.LBool(false))
					return 1
				}
				L.Push(lua.LBool(q.SetRecordExField(key, value)))
				return 1
			default:
				L.ArgError(2, "record_ex(key[, value]) requires key")
				return 0
			}
		},
		"start_time": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			argc := L.GetTop()
			switch argc {
			case 1:
				if q.StartTime.IsZero() {
					L.Push(lua.LNumber(0))
				} else {
					L.Push(lua.LNumber(q.StartTime.Unix()))
				}
				return 1
			case 2:
				unix := L.CheckInt64(2)
				if unix <= 0 {
					q.ResetStartTime()
				} else {
					q.SetStartTime(time.Unix(unix, 0))
				}
				return 0
			default:
				L.ArgError(2, "start_time() takes 0 or 1 arguments")
				return 0
			}
		},
		"start": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			opts := QuestPhaseOpts{}
			if q.Wz == nil {
				if L.GetTop() != 2 {
					L.ArgError(2, "start(record) requires record when quest has no WZ definition")
					return 0
				}
				record := L.CheckString(2)
				opts.Force = true
				opts.Record = &record
			} else {
				if L.GetTop() < 2 || L.GetTop() > 3 {
					L.ArgError(2, "start(npc[, force|record]) requires npc id")
					return 0
				}
				npcID := uint32(L.CheckInt(2))
				npc := npcID
				opts.NpcID = &npc
				if L.GetTop() >= 3 {
					switch L.Get(3).Type() {
					case lua.LTBool:
						opts.Force = L.CheckBool(3)
					case lua.LTString:
						record := L.CheckString(3)
						opts.Force = true
						opts.Record = &record
					default:
						L.ArgError(3, "start(npc[, force|record]) second argument must be boolean or string")
						return 0
					}
				}
			}
			qp, err := q.container.Start(q.QuestID, opts)
			if err == nil && qp != nil {
				if ud, ok := L.Get(1).(*lua.LUserData); ok {
					ud.Value = qp
				}
			}
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"complete": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() < 2 || L.GetTop() > 3 {
				L.ArgError(2, "complete(npc[, selection]) requires npc id")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			opts := QuestPhaseOpts{NpcID: &npcID}
			if L.GetTop() >= 3 {
				sel := uint32(L.CheckInt(3))
				opts.Selection = &sel
			}
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if stored := q.container.Get(q.QuestID); stored != nil {
				q = stored
			}
			err := q.Complete(q.container.owner, opts)
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"force_complete": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "force_complete(npc) requires npc id")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if stored := q.container.Get(q.QuestID); stored != nil {
				q = stored
			}
			npc := npcID
			err := q.Complete(q.container.owner, QuestPhaseOpts{NpcID: &npc, Force: true})
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"forfeit": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "forfeit() takes no arguments")
				return 0
			}
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if stored := q.container.Get(q.QuestID); stored != nil {
				q = stored
			}
			err := q.Forfeit(q.container.owner)
			L.Push(lua.LBool(err == nil))
			return 1
		},
	}
}

func (ch *Character) RunQuestScript(actx actor.Context, questID uint32, npcID uint32, entry string) error {
	if ch == nil {
		return fmt.Errorf("character is nil")
	}
	if ch.GetDialog() != nil {
		return fmt.Errorf("dialog already active")
	}
	mapInstance := ch.GetMap()
	if mapInstance == nil {
		return fmt.Errorf("map not found")
	}
	root := mapInstance.GetLuaRoot()
	if root == nil {
		return fmt.Errorf("lua state not available")
	}
	scriptPath := fmt.Sprintf("script/quest/%d.lua", questID)
	luaThread, err := luax.NewThread(root, scriptPath)
	if err != nil {
		return err
	}
	luax.SetConfiguration(luaThread, luax.Configuration{
		ActorContext: actx,
		MapActorPID:  mapInstance.GetActorPID(),
		KeepAlive:    true,
	})
	luax.CallAsync(root, luaThread, entry, ch, npcID).Then(func(_ interface{}) (interface{}, error) {
		if ch.GetDialog() == nil {
			ch.ResetDialog()
			if ch.Listener != nil {
				ch.Listener.OnUnlockAction(ch)
			}
		}
		return nil, nil
	}).OnError(func(err error) {
		log.Printf("quest script %s quest=%d npc=%d: %v", entry, questID, npcID, err)
		if ch.GetDialog() == nil {
			ch.ResetDialog()
			if ch.Listener != nil {
				ch.Listener.OnUnlockAction(ch)
			}
		}
	})
	return nil
}
