package entity

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) LuaQuest(questID uint32) *Quest {
	if ch == nil || ch.Quests == nil {
		return nil
	}
	if qp := ch.Quests.Get(questID); qp != nil {
		return qp
	}
	def := ch.Quests.questDef(questID)
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
				L.Push(lua.LString(q.StatusRecord))
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
			q.StatusRecord = L.CheckString(2)
			L.Push(lua.LBool(true))
			return 1
		},
		"is_startable": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			npcID := uint32(0)
			opts := QuestPrepareOpts{}
			if L.GetTop() >= 2 {
				npcID = uint32(L.CheckInt(2))
				opts.NpcID = &npcID
			}
			if q.container == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			L.Push(lua.LBool(q.container.IsStartable(q.QuestID, opts)))
			return 1
		},
		"is_completable": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "is_completable() takes no arguments")
				return 0
			}
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			L.Push(lua.LBool(q.IsCompletable(q.container.owner)))
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
		"start": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if L.GetTop() != 2 {
				L.ArgError(2, "start(npc) requires npc id")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			npc := npcID
			qp, err := q.container.Start(q.QuestID, QuestPrepareOpts{NpcID: &npc})
			if err == nil && qp != nil {
				if ud, ok := L.Get(1).(*lua.LUserData); ok {
					ud.Value = qp
				}
			}
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"force_start": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
			}
			if q.container == nil || q.container.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			if q.Wz == nil {
				if L.GetTop() != 2 {
					L.ArgError(2, "force_start(record) requires record when quest has no WZ definition")
					return 0
				}
				record := L.CheckString(2)
				qp, err := q.container.ForceStart(q.QuestID, record)
				if err == nil && qp != nil {
					if ud, ok := L.Get(1).(*lua.LUserData); ok {
						ud.Value = qp
					}
				}
				L.Push(lua.LBool(err == nil))
				return 1
			}
			if L.GetTop() < 2 || L.GetTop() > 3 {
				L.ArgError(2, "force_start(npc[, record]) requires npc id")
				return 0
			}
			npcID := uint32(L.CheckInt(2))
			record := ""
			if L.GetTop() >= 3 {
				record = L.CheckString(3)
			}
			npc := npcID
			qp, err := q.container.Start(q.QuestID, QuestPrepareOpts{NpcID: &npc, Force: true})
			if err == nil && qp != nil {
				if record != "" {
					qp.StatusRecord = record
				}
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
			opts := QuestPrepareOpts{NpcID: &npcID}
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
			err := q.Complete(q.container.owner, QuestPrepareOpts{NpcID: &npc, Force: true})
			L.Push(lua.LBool(err == nil))
			return 1
		},
	}
}
