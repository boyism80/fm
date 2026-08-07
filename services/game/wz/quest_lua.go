package wz

import (
	"sort"

	lua "github.com/yuin/gopher-lua"
)

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

func questRequirementsToLuaTable(L *lua.LState, reqs QuestRequirements) *lua.LTable {
	out := L.NewTable()
	i := 1
	appendReq := func(tbl *lua.LTable) {
		out.RawSetInt(i, tbl)
		i++
	}
	if reqs.NPC != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("npc"))
		tbl.RawSetString("value", lua.LNumber(reqs.NPC))
		appendReq(tbl)
	}
	if reqs.LevelMin != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("lvmin"))
		tbl.RawSetString("value", lua.LNumber(reqs.LevelMin))
		appendReq(tbl)
	}
	if reqs.LevelMax != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("lvmax"))
		tbl.RawSetString("value", lua.LNumber(reqs.LevelMax))
		appendReq(tbl)
	}
	if reqs.Level != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("level"))
		tbl.RawSetString("value", lua.LNumber(reqs.Level))
		appendReq(tbl)
	}
	if len(reqs.Job) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("class"))
		tbl.RawSetString("classes", intSliceToLuaTable(L, reqs.Job))
		appendReq(tbl)
	}
	if len(reqs.Item) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("item"))
		tbl.RawSetString("items", questItemCountsToLuaTable(L, reqs.Item))
		appendReq(tbl)
	}
	if len(reqs.Mob) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("mob"))
		tbl.RawSetString("mobs", questMobCountsToLuaTable(L, reqs.Mob))
		appendReq(tbl)
	}
	if len(reqs.Quest) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("quest"))
		tbl.RawSetString("quests", questStatesToLuaTable(L, reqs.Quest))
		appendReq(tbl)
	}
	if len(reqs.Skill) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("skill"))
		tbl.RawSetString("skills", questSkillAcquiresToLuaTable(L, reqs.Skill))
		appendReq(tbl)
	}
	if len(reqs.Pet) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("pet"))
		tbl.RawSetString("pet_ids", uint32SliceToLuaTable(L, reqs.Pet))
		appendReq(tbl)
	}
	if reqs.Pop != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("pop"))
		tbl.RawSetString("value", lua.LNumber(reqs.Pop))
		appendReq(tbl)
	}
	if reqs.HasInterval {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("interval"))
		tbl.RawSetString("value", lua.LNumber(reqs.Interval))
		appendReq(tbl)
	}
	if reqs.FieldEnter != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("fieldEnter"))
		tbl.RawSetString("value", lua.LNumber(reqs.FieldEnter))
		appendReq(tbl)
	}
	if reqs.QuestComplete != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("questComplete"))
		tbl.RawSetString("value", lua.LNumber(reqs.QuestComplete))
		appendReq(tbl)
	}
	if reqs.StartScript != "" {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("startscript"))
		tbl.RawSetString("value", lua.LString(reqs.StartScript))
		appendReq(tbl)
	}
	if reqs.EndScript != "" {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("endscript"))
		tbl.RawSetString("value", lua.LString(reqs.EndScript))
		appendReq(tbl)
	}
	if reqs.Start != "" {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("start"))
		tbl.RawSetString("value", lua.LString(reqs.Start))
		appendReq(tbl)
	}
	if reqs.End != "" {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("end"))
		tbl.RawSetString("value", lua.LString(reqs.End))
		appendReq(tbl)
	}
	if len(reqs.Info) > 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("info"))
		tbl.RawSetString("value", lua.LNumber(0))
		appendReq(tbl)
	}
	if reqs.InfoNumber != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("infoNumber"))
		tbl.RawSetString("value", lua.LNumber(reqs.InfoNumber))
		appendReq(tbl)
	}
	if reqs.DayByDay {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("dayByDay"))
		tbl.RawSetString("value", lua.LNumber(1))
		appendReq(tbl)
	}
	if reqs.NormalAutoStart {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("normalAutoStart"))
		tbl.RawSetString("value", lua.LNumber(1))
		appendReq(tbl)
	}
	if reqs.PartyQuestS != 0 {
		tbl := L.NewTable()
		tbl.RawSetString("kind", lua.LString("partyQuest_S"))
		tbl.RawSetString("value", lua.LNumber(reqs.PartyQuestS))
		appendReq(tbl)
	}
	return out
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
