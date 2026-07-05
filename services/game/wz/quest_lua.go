package wz

import lua "github.com/yuin/gopher-lua"

func (q *Quest) ToLuaTable(L *lua.LState) *lua.LTable {
	tbl := L.NewTable()
	if q == nil {
		return tbl
	}
	tbl.RawSetString("id", lua.LNumber(q.ID))
	tbl.RawSetString("name", lua.LString(q.Meta.Name))
	tbl.RawSetString("complete_requirements", questRequirementsToLuaTable(L, q.Complete.Requirements))
	tbl.RawSetString("start_requirements", questRequirementsToLuaTable(L, q.Start.Requirements))
	return tbl
}

func questRequirementsToLuaTable(L *lua.LState, reqs []QuestRequirement) *lua.LTable {
	out := L.NewTable()
	if len(reqs) == 0 {
		return out
	}
	for i, req := range reqs {
		out.RawSetInt(i+1, questRequirementToLuaTable(L, req))
	}
	return out
}

func questRequirementToLuaTable(L *lua.LState, req QuestRequirement) *lua.LTable {
	tbl := L.NewTable()
	tbl.RawSetString("kind", lua.LString(questRequirementKindForLua(req.Kind)))

	switch req.Kind {
	case QuestReqItem:
		tbl.RawSetString("items", questItemCountsToLuaTable(L, req.Items))
	case QuestReqMob:
		tbl.RawSetString("mobs", questMobCountsToLuaTable(L, req.Mobs))
	case QuestReqQuest:
		tbl.RawSetString("quests", questStateRefsToLuaTable(L, req.Quests))
	case QuestReqClass:
		tbl.RawSetString("classes", intSliceToLuaTable(L, req.Classes))
	case QuestReqSkill:
		tbl.RawSetString("skills", questSkillRefsToLuaTable(L, req.Skills))
	case QuestReqPet:
		tbl.RawSetString("pet_ids", uint32SliceToLuaTable(L, req.PetIDs))
	case QuestReqTimeStart, QuestReqTimeEnd, QuestReqStartScript, QuestReqEndScript:
		tbl.RawSetString("value", lua.LString(req.StrValue))
	default:
		tbl.RawSetString("value", lua.LNumber(req.IntValue))
	}

	return tbl
}

func questRequirementKindForLua(kind QuestRequirementKind) string {
	if kind == QuestReqClass {
		return "class"
	}
	return string(kind)
}

func questItemCountsToLuaTable(L *lua.LState, items []QuestItemCount) *lua.LTable {
	out := L.NewTable()
	for i, item := range items {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(item.ItemID))
		entry.RawSetString("count", lua.LNumber(item.Count))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func questMobCountsToLuaTable(L *lua.LState, mobs []QuestMobCount) *lua.LTable {
	out := L.NewTable()
	for i, mob := range mobs {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(mob.MobID))
		entry.RawSetString("count", lua.LNumber(mob.Count))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func questStateRefsToLuaTable(L *lua.LState, quests []QuestStateRef) *lua.LTable {
	out := L.NewTable()
	for i, ref := range quests {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(ref.QuestID))
		entry.RawSetString("state", lua.LNumber(ref.State))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func questSkillRefsToLuaTable(L *lua.LState, skills []QuestSkillRef) *lua.LTable {
	out := L.NewTable()
	for i, skill := range skills {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(skill.SkillID))
		entry.RawSetString("acquire", lua.LNumber(skill.Acquire))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func intSliceToLuaTable(L *lua.LState, values []int) *lua.LTable {
	out := L.NewTable()
	for i, v := range values {
		out.RawSetInt(i+1, lua.LNumber(v))
	}
	return out
}

func uint32SliceToLuaTable(L *lua.LState, values []uint32) *lua.LTable {
	out := L.NewTable()
	for i, v := range values {
		out.RawSetInt(i+1, lua.LNumber(v))
	}
	return out
}
