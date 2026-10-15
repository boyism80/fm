package wz

import lua "github.com/yuin/gopher-lua"

func (d *MobSkillLevelData) ToLuaTable(L *lua.LState) *lua.LTable {
	tbl := L.NewTable()
	if d == nil {
		return tbl
	}
	tbl.RawSetString("skill_id", lua.LNumber(d.SkillID))
	tbl.RawSetString("level", lua.LNumber(d.Level))
	tbl.RawSetString("hp_percent", lua.LNumber(d.HpPercent))
	tbl.RawSetString("mp_con", lua.LNumber(d.MpCon))
	tbl.RawSetString("x", lua.LNumber(d.X))
	tbl.RawSetString("y", lua.LNumber(d.Y))
	tbl.RawSetString("time", lua.LNumber(d.DurationMs))
	tbl.RawSetString("cooltime_ms", lua.LNumber(d.CooltimeMs))
	propPct := int(d.Prop * 100)
	if propPct <= 0 {
		propPct = 100
	}
	tbl.RawSetString("prop", lua.LNumber(propPct))
	tbl.RawSetString("limit", lua.LNumber(d.Limit))
	tbl.RawSetString("spawn_effect", lua.LNumber(d.SpawnEffect))
	tbl.RawSetString("summon_once", lua.LBool(d.SummonOnce))
	summons := L.NewTable()
	for i, id := range d.Summons {
		summons.RawSetInt(i+1, lua.LNumber(id))
	}
	tbl.RawSetString("summons", summons)
	if d.ElemAttr != "" {
		tbl.RawSetString("elem_attr", lua.LString(d.ElemAttr))
	}
	bounds := L.NewTable()
	bounds.RawSetString("left", lua.LNumber(d.Bounds.Left))
	bounds.RawSetString("top", lua.LNumber(d.Bounds.Top))
	bounds.RawSetString("right", lua.LNumber(d.Bounds.Right))
	bounds.RawSetString("bottom", lua.LNumber(d.Bounds.Bottom))
	tbl.RawSetString("bounds", bounds)
	return tbl
}
