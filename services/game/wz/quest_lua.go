package wz

import (
	"sort"

	lua "github.com/yuin/gopher-lua"
)

func (q *Quest) ToLuaTable(L *lua.LState) *lua.LTable {
	tbl := L.NewTable()
	if q == nil {
		return tbl
	}
	tbl.RawSetString("id", lua.LNumber(q.ID))
	tbl.RawSetString("name", lua.LString(q.Meta.Name))
	tbl.RawSetString("complete_requirements", questRequirementsToLuaTable(L, q.Complete.Requirements))
	tbl.RawSetString("start_requirements", questRequirementsToLuaTable(L, q.Start.Requirements))
	if len(q.PartyRanks) > 0 {
		tbl.RawSetString("party_ranks", partyRanksToLuaTable(L, q.PartyRanks))
	}
	return tbl
}

func partyRanksToLuaTable(L *lua.LState, ranks map[string][]PartyQuestRankCheck) *lua.LTable {
	out := L.NewTable()
	if len(ranks) == 0 {
		return out
	}
	keys := make([]string, 0, len(ranks))
	for rank := range ranks {
		keys = append(keys, rank)
	}
	sort.Strings(keys)
	for _, rank := range keys {
		checks := ranks[rank]
		arr := L.NewTable()
		for i, check := range checks {
			entry := L.NewTable()
			entry.RawSetString("mode", lua.LString(check.Mode))
			entry.RawSetString("property", lua.LString(check.Property))
			entry.RawSetString("value", lua.LNumber(check.Value))
			arr.RawSetInt(i+1, entry)
		}
		out.RawSetString(rank, arr)
	}
	return out
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
		tbl.RawSetString("quests", questStatesToLuaTable(L, req.Quests))
	case QuestReqClass:
		tbl.RawSetString("classes", intSliceToLuaTable(L, req.Classes))
	case QuestReqSkill:
		tbl.RawSetString("skills", questSkillAcquiresToLuaTable(L, req.Skills))
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

func questItemCountsToLuaTable(L *lua.LState, items map[uint32]int) *lua.LTable {
	out := L.NewTable()
	if len(items) == 0 {
		return out
	}
	ids := make([]uint32, 0, len(items))
	for id := range items {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for i, id := range ids {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(id))
		entry.RawSetString("count", lua.LNumber(items[id]))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func questMobCountsToLuaTable(L *lua.LState, mobs map[uint32]int) *lua.LTable {
	out := L.NewTable()
	if len(mobs) == 0 {
		return out
	}
	ids := make([]uint32, 0, len(mobs))
	for id := range mobs {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for i, id := range ids {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(id))
		entry.RawSetString("count", lua.LNumber(mobs[id]))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func questStatesToLuaTable(L *lua.LState, quests map[uint32]QuestStatus) *lua.LTable {
	out := L.NewTable()
	if len(quests) == 0 {
		return out
	}
	ids := make([]uint32, 0, len(quests))
	for id := range quests {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for i, id := range ids {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(id))
		entry.RawSetString("state", lua.LNumber(quests[id]))
		out.RawSetInt(i+1, entry)
	}
	return out
}

func questSkillAcquiresToLuaTable(L *lua.LState, skills map[uint32]int) *lua.LTable {
	out := L.NewTable()
	if len(skills) == 0 {
		return out
	}
	ids := make([]uint32, 0, len(skills))
	for id := range skills {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	for i, id := range ids {
		entry := L.NewTable()
		entry.RawSetString("id", lua.LNumber(id))
		entry.RawSetString("acquire", lua.LNumber(skills[id]))
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
