package wz

import lua "github.com/yuin/gopher-lua"

func (d *SkillLevelData) ToLuaTable(L *lua.LState) *lua.LTable {
	tbl := L.NewTable()
	if d == nil {
		return tbl
	}
	tbl.RawSetString("mp_con", lua.LNumber(d.MPCon))
	tbl.RawSetString("hp_con", lua.LNumber(d.HPCon))
	tbl.RawSetString("money_con", lua.LNumber(d.MoneyCon))
	tbl.RawSetString("item_con", lua.LNumber(d.ItemCon))
	tbl.RawSetString("item_con_no", lua.LNumber(d.ItemConNo))
	tbl.RawSetString("item_consume", lua.LNumber(d.ItemConsume))
	tbl.RawSetString("bullet_consume", lua.LNumber(d.BulletConsume))
	tbl.RawSetString("bullet_count", lua.LNumber(d.BulletCount))
	tbl.RawSetString("damage", lua.LNumber(d.Damage))
	tbl.RawSetString("damage_pc", lua.LNumber(d.DamagePC))
	tbl.RawSetString("fix_damage", lua.LNumber(d.FixDamage))
	tbl.RawSetString("critical_damage", lua.LNumber(d.CriticalDamage))
	tbl.RawSetString("attack_count", lua.LNumber(d.AttackCount))
	tbl.RawSetString("mob_count", lua.LNumber(d.MobCount))
	tbl.RawSetString("pad", lua.LNumber(d.PAD))
	tbl.RawSetString("mad", lua.LNumber(d.MAD))
	tbl.RawSetString("pdd", lua.LNumber(d.PDD))
	tbl.RawSetString("mdd", lua.LNumber(d.MDD))
	tbl.RawSetString("eva", lua.LNumber(d.EVA))
	tbl.RawSetString("acc", lua.LNumber(d.ACC))
	tbl.RawSetString("str", lua.LNumber(d.STR))
	tbl.RawSetString("hp", lua.LNumber(d.HP))
	tbl.RawSetString("mp", lua.LNumber(d.MP))
	tbl.RawSetString("jump", lua.LNumber(d.Jump))
	tbl.RawSetString("speed", lua.LNumber(d.Speed))
	tbl.RawSetString("mastery", lua.LNumber(d.Mastery))
	tbl.RawSetString("prop", lua.LNumber(d.Prop))
	tbl.RawSetString("range", lua.LNumber(d.Range))
	tbl.RawSetString("time", lua.LNumber(d.Time.Milliseconds()))
	tbl.RawSetString("cooldown", lua.LNumber(d.Cooldown.Milliseconds()))
	tbl.RawSetString("morph", lua.LNumber(d.Morph))
	tbl.RawSetString("x", lua.LNumber(d.X))
	tbl.RawSetString("y", lua.LNumber(d.Y))
	tbl.RawSetString("z", lua.LNumber(d.Z))
	ltTable := L.NewTable()
	ltTable.RawSetString("x", lua.LNumber(d.LT.X))
	ltTable.RawSetString("y", lua.LNumber(d.LT.Y))
	tbl.RawSetString("lt", ltTable)
	rbTable := L.NewTable()
	rbTable.RawSetString("x", lua.LNumber(d.RB.X))
	rbTable.RawSetString("y", lua.LNumber(d.RB.Y))
	tbl.RawSetString("rb", rbTable)
	if d.HS != "" {
		tbl.RawSetString("hs", lua.LString(d.HS))
	}
	if d.Action != "" {
		tbl.RawSetString("action", lua.LString(d.Action))
	}
	return tbl
}
