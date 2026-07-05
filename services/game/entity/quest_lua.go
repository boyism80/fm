package entity

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) LuaQuest(questID uint32) *Quest {
	if ch == nil || ch.Quests == nil || ch.GameWorld == nil {
		return nil
	}
	if qp := ch.Quests.Get(questID); qp != nil {
		return qp
	}
	def := ch.GameWorld.GetResources().GetQuest(questID)
	if def == nil {
		return nil
	}
	return &Quest{
		owner:   ch,
		Wz:      def,
		QuestID: questID,
		Status:  QuestStatusNotStarted,
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
			if q.owner == nil || q.owner.Quests.Get(q.QuestID) == nil {
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
			if L.GetTop() >= 2 {
				npcID = uint32(L.CheckInt(2))
			}
			if q.owner == nil || q.Wz == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			L.Push(lua.LBool(q.owner.Quests.IsStartable(q.Wz, npcID, QuestPrepareOpts{IgnoreScriptRequirement: true})))
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
			if q.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			L.Push(lua.LBool(q.IsCompletable(q.owner)))
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
			if q.owner == nil || q.Wz == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			_, err := q.owner.Quests.Start(q.Wz, npcID, QuestPrepareOpts{IgnoreScriptRequirement: true})
			L.Push(lua.LBool(err == nil))
			return 1
		},
		"force_start": func(L *lua.LState) int {
			q, ok := LuaCheckQuest(L, 1)
			if !ok {
				return 0
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
			if q.owner == nil || q.Wz == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			qp, err := q.owner.Quests.Start(q.Wz, npcID, QuestPrepareOpts{Force: true})
			if err == nil && record != "" && qp != nil {
				qp.StatusRecord = record
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
			var selection *uint32
			if L.GetTop() >= 3 {
				sel := uint32(L.CheckInt(3))
				selection = &sel
			}
			if q.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			err := q.Complete(q.owner, npcID, selection, QuestPrepareOpts{IgnoreScriptRequirement: true})
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
			if q.owner == nil {
				L.Push(lua.LBool(false))
				return 1
			}
			err := q.Complete(q.owner, npcID, nil, QuestPrepareOpts{Force: true})
			L.Push(lua.LBool(err == nil))
			return 1
		},
	}
}
