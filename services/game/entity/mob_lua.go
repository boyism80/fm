package entity

import (
	"fmt"
	"time"

	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/services/game/constant"
	"github.com/boyism80/fm/services/game/wz"
	lua "github.com/yuin/gopher-lua"
)

func (m *Mob) LuaTypeName() string {
	return "LuaMob"
}

func (m *Mob) LuaBuiltinFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {

				L.Push(lua.LNumber(mob.Wz.ID))
				return 1
			} else {
				L.ArgError(2, "id() is read-only")
				return 0
			}
		},
		"name": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {

				L.Push(lua.LString(fmt.Sprintf("Mob_%d", mob.Wz.ID)))
				return 1
			} else {
				L.ArgError(2, "name() is read-only")
				return 0
			}
		},
		"exp": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc == 1 {
				L.Push(lua.LNumber(mob.Wz.EXP))
				return 1
			} else {
				L.ArgError(2, "exp() is read-only")
				return 0
			}
		},
		"foothold": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(mob.Foothold))
				return 1
			case 2:
				foothold := L.CheckInt(2)
				mob.Foothold = int16(foothold)
				return 0
			default:
				L.ArgError(2, "foothold() requires 0 or 1 arguments")
				return 0
			}
		},
		"map_id": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			if L.GetTop() != 1 {
				L.ArgError(2, "map_id() is read-only (map is fixed at spawn)")
				return 0
			}
			mapInstance := mob.GetMap()
			if mapInstance == nil {
				L.Push(lua.LNumber(0))
				return 1
			}
			L.Push(lua.LNumber(mapInstance.id))
			return 1
		},
		"buff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			buffTable := L.CheckTable(2)
			maskLV := buffTable.RawGetString("mask")
			parseMobBuffs := func(v lua.LValue) (constant.MobBuffFlag, bool) {
				t, ok := v.(*lua.LTable)
				if !ok {
					return 0, false
				}
				m := t.RawGetString("mask")
				if m.Type() != lua.LTNumber {
					return 0, false
				}
				return constant.MobBuffFlag(uint32(lua.LVAsNumber(m))), true
			}
			if maskLV.Type() == lua.LTNumber {
				buff := constant.MobBuffFlag(uint32(lua.LVAsNumber(maskLV)))
				value := int32(L.CheckInt(3))
				durationMs := L.CheckInt64(4)
				if durationMs < 0 {
					durationMs = 0
				}
				var skillWz *wz.Skill
				var skillLevel uint8
				if L.GetTop() >= 5 {
					if lv, ok := L.Get(5).(*lua.LUserData); ok {
						if se, ok := lv.Value.(*SkillEntry); ok && se != nil && se.Wz != nil {
							skillWz = se.Wz
							skillLevel = uint8(se.Level())
						} else if m, ok := lv.Value.(*Mist); ok && m != nil && m.SkillWz != nil {
							skillWz = m.SkillWz
							skillLevel = m.SkillLevel
						} else if sb, ok := lv.Value.(*SkillBuff); ok && sb != nil && sb.Wz != nil {
							skillWz = sb.Wz
							skillLevel = sb.SkillLevel
						} else if mb, ok := lv.Value.(*MobSkillBuff); ok && mb != nil && mb.Wz != nil {
							skillWz = mb.Wz
							skillLevel = mb.SkillLevel
						}
					}
				}
				var causerOID uint32
				if L.GetTop() >= 6 && L.Get(6) != lua.LNil {
					switch cv := L.Get(6).(type) {
					case *lua.LUserData:
						cud := cv
						if ch, ok := cud.Value.(*Character); ok && ch != nil {
							causerOID = ch.GetID()
						} else {
							L.ArgError(6, "causer must be Character, OID(number), or nil")
							return 0
						}
					case lua.LNumber:
						oid := int64(lua.LVAsNumber(cv))
						if oid < 0 {
							L.ArgError(6, "causer OID must be >= 0")
							return 0
						}
						causerOID = uint32(oid)
					default:
						L.ArgError(6, "causer must be Character, OID(number), or nil")
						return 0
					}
				}
				stack := uint8(1)
				if L.GetTop() >= 7 && L.Get(7) != lua.LNil {
					s := L.CheckInt(7)
					if s < 1 {
						s = 1
					}
					if s > 255 {
						s = 255
					}
					stack = uint8(s)
				}
				now := time.Now()
				mob.ensureMobBuffs().AddSkillBuff(now, durationMs, skillWz, skillLevel, causerOID,
					map[constant.MobBuffFlag]int32{buff: value},
					map[constant.MobBuffFlag]uint8{buff: stack})
				return 0
			}
			durationMs := L.CheckInt64(3)
			if durationMs < 0 {
				durationMs = 0
			}
			var skillWz *wz.Skill
			var skillLevel uint8
			if L.GetTop() >= 4 {
				if lv, ok := L.Get(4).(*lua.LUserData); ok {
					if se, ok := lv.Value.(*SkillEntry); ok && se != nil && se.Wz != nil {
						skillWz = se.Wz
						skillLevel = uint8(se.Level())
					} else if m, ok := lv.Value.(*Mist); ok && m != nil && m.SkillWz != nil {
						skillWz = m.SkillWz
						skillLevel = m.SkillLevel
					} else if sb, ok := lv.Value.(*SkillBuff); ok && sb != nil && sb.Wz != nil {
						skillWz = sb.Wz
						skillLevel = sb.SkillLevel
					} else if mb, ok := lv.Value.(*MobSkillBuff); ok && mb != nil && mb.Wz != nil {
						skillWz = mb.Wz
						skillLevel = mb.SkillLevel
					}
				}
			}
			var causerOID uint32
			if L.GetTop() >= 5 && L.Get(5) != lua.LNil {
				switch cv := L.Get(5).(type) {
				case *lua.LUserData:
					cud := cv
					if ch, ok := cud.Value.(*Character); ok && ch != nil {
						causerOID = ch.GetID()
					} else {
						L.ArgError(5, "causer must be Character, OID(number), or nil")
						return 0
					}
				case lua.LNumber:
					oid := int64(lua.LVAsNumber(cv))
					if oid < 0 {
						L.ArgError(5, "causer OID must be >= 0")
						return 0
					}
					causerOID = uint32(oid)
				default:
					L.ArgError(5, "causer must be Character, OID(number), or nil")
					return 0
				}
			}
			values := make(map[constant.MobBuffFlag]int32)
			buffTable.ForEach(func(key lua.LValue, value lua.LValue) {
				buff, ok := parseMobBuffs(key)
				if !ok {
					return
				}
				if value.Type() != lua.LTNumber {
					return
				}
				values[buff] = int32(lua.LVAsNumber(value))
			})
			if len(values) > 0 {
				now := time.Now()
				mob.ensureMobBuffs().AddSkillBuff(now, durationMs, skillWz, skillLevel, causerOID, values, nil)
			}
			return 0
		},
		"exp_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			if L.GetTop() == 1 {
				v := mob.ExpRate
				if v <= 0 {
					v = 100
				}
				L.Push(lua.LNumber(v))
				return 1
			}
			if L.GetTop() == 2 {
				n := int32(L.CheckInt(2))
				if n < 100 {
					n = 100
				}
				mob.ExpRate = n
				return 0
			}
			L.ArgError(2, "exp_rate() get or set one value")
			return 0
		},
		"drop_rate": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			if L.GetTop() == 1 {
				v := mob.DropRate
				if v <= 0 {
					v = 100
				}
				L.Push(lua.LNumber(v))
				return 1
			}
			if L.GetTop() == 2 {
				n := int32(L.CheckInt(2))
				if n < 100 {
					n = 100
				}
				mob.DropRate = n
				return 0
			}
			L.ArgError(2, "drop_rate() get or set one value")
			return 0
		},
		"drops": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			t := L.NewTable()
			if mob.Context == nil || mob.Wz == nil {
				L.Push(t)
				return 1
			}
			resources := mob.Context.GetResources()
			if resources == nil {
				L.Push(t)
				return 1
			}
			mobDrops, ok := resources.Drops[mob.Wz.ID]
			if !ok {
				L.Push(t)
				return 1
			}
			idx := 1
			for _, d := range mobDrops {
				row := L.NewTable()
				row.RawSetString("item", lua.LNumber(d.Item))
				row.RawSetString("money", lua.LNumber(d.Money))
				row.RawSetString("prob", lua.LNumber(d.Prob))
				row.RawSetString("min", lua.LNumber(d.Min))
				row.RawSetString("max", lua.LNumber(d.Max))
				row.RawSetString("quest", lua.LNumber(d.QuestID))
				t.RawSetInt(idx, row)
				idx++
			}
			L.Push(t)
			return 1
		},
		"has_stolen": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			L.Push(lua.LBool(mob.stealOutcome != nil))
			return 1
		},
		"homing": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			if L.GetTop() != 3 {
				L.ArgError(2, "homing(causer, skill_or_nil) requires causer and skill (nil to clear)")
				return 0
			}
			causerUd := L.CheckUserData(2)
			causer, causerOk := causerUd.Value.(*Character)
			if !causerOk || causer == nil {
				L.ArgError(2, "Character expected")
				return 0
			}
			if L.Get(3) == lua.LNil {
				mob.SetHoming(causer.GetID(), nil)
				return 0
			}
			skillUd := L.CheckUserData(3)
			skill, skillOk := skillUd.Value.(*SkillEntry)
			if !skillOk || skill == nil {
				L.ArgError(3, "Skill expected or nil to clear")
				return 0
			}
			if skill.Wz == nil {
				L.ArgError(3, "Skill has no wz data")
				return 0
			}
			lv := skill.Level()
			if lv < 0 || lv > 255 {
				L.ArgError(3, "skill level out of range for homing")
				return 0
			}
			h := &Homing{SkillWz: skill.Wz, SkillLevel: uint8(lv)}
			mob.SetHoming(causer.GetID(), h)
			return 0
		},
		"record_stolen_item": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			if L.GetTop() < 2 {
				return 0
			}
			lv := L.Get(2)
			if lv == lua.LNil {
				return 0
			}
			n, ok := lv.(lua.LNumber)
			if !ok {
				return 0
			}
			v := uint32(n)
			if v == 0 {
				return 0
			}
			mob.stealOutcome = &v
			return 0
		},
		"clear_buffs": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			buffTable := L.CheckTable(2)
			maskLV := buffTable.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobBuff table with numeric mask expected")
				return 0
			}
			buff := constant.MobBuffFlag(uint32(lua.LVAsNumber(maskLV)))
			mob.CancelMobBuff(buff)
			return 0
		},
		"has_buff": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			buffTable := L.CheckTable(2)
			maskLV := buffTable.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobBuff table with numeric mask expected")
				return 0
			}
			buff := constant.MobBuffFlag(uint32(lua.LVAsNumber(maskLV)))
			L.Push(lua.LBool(mob.HasBuff(buff)))
			return 1
		},
		"buff_value": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			buffTable := L.CheckTable(2)
			maskLV := buffTable.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobBuff table with numeric mask expected")
				return 0
			}
			buff := constant.MobBuffFlag(uint32(lua.LVAsNumber(maskLV)))
			L.Push(lua.LNumber(mob.GetMobBuffValue(buff)))
			return 1
		},

		"buff_stack": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok || mob == nil {
				L.ArgError(1, "Mob expected")
				return 0
			}
			buffTable := L.CheckTable(2)
			maskLV := buffTable.RawGetString("mask")
			if maskLV.Type() != lua.LTNumber {
				L.ArgError(2, "MobBuff table with numeric mask expected")
				return 0
			}
			buff := constant.MobBuffFlag(uint32(lua.LVAsNumber(maskLV)))
			if L.GetTop() >= 3 {
				stack := uint8(L.CheckInt(3))
				if stack < 1 {
					stack = 1
				}
				L.Push(lua.LBool(mob.SetMobBuffStack(buff, stack)))
				return 1
			}
			L.Push(lua.LNumber(mob.GetMobBuffStack(buff)))
			return 1
		},
		"wz": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			tbl := L.NewTable()
			if mob.Wz != nil {
				tbl.RawSetString("id", lua.LNumber(mob.Wz.ID))
				tbl.RawSetString("level", lua.LNumber(mob.Wz.Level))
				tbl.RawSetString("max_hp", lua.LNumber(mob.Wz.MaxHP))
				tbl.RawSetString("max_mp", lua.LNumber(mob.Wz.MaxMP))
				tbl.RawSetString("exp", lua.LNumber(mob.Wz.EXP))
				tbl.RawSetString("boss", lua.LBool(mob.Wz.Boss))
				if len(mob.Wz.ElemResist) > 0 {
					er := L.NewTable()
					for k, v := range mob.Wz.ElemResist {
						er.RawSetString(k, lua.LNumber(v))
					}
					tbl.RawSetString("elem_resist", er)
				}
			}
			L.Push(tbl)
			return 1
		},
		"damage": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}

			argc := L.GetTop()
			if argc != 3 {
				L.ArgError(2, "damage(attacker, amount) requires attacker (Character or nil) and amount")
				return 0
			}

			var attacker *Character
			if L.Get(2) != lua.LNil {
				if aud, ok := L.Get(2).(*lua.LUserData); ok {
					if ch, ok := aud.Value.(*Character); ok {
						attacker = ch
					} else {
						L.ArgError(2, "attacker must be Character or nil")
						return 0
					}
				} else {
					L.ArgError(2, "attacker must be Character or nil")
					return 0
				}
			}

			amount := L.CheckInt(3)
			if amount <= 0 {
				L.Push(lua.LBool(false))
				return 1
			}

			killed := mob.ApplyDamage(attacker, uint32(amount))
			L.Push(lua.LBool(killed))
			return 1
		},
		"controller": func(L *lua.LState) int {
			ud := L.CheckUserData(1)
			mob, ok := ud.Value.(*Mob)
			if !ok {
				L.ArgError(1, "Mob expected")
				return 0
			}
			argc := L.GetTop()
			if argc == 1 {
				mapInstance := mob.GetMap()
				if mapInstance == nil {
					L.Push(lua.LNil)
					return 1
				}
				ct := mapInstance.GetControllerTable()
				controller, ok := ct.GetController(mob)
				if !ok {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(luax.NewLuable(L, controller))
				return 1
			}
			if argc == 2 {
				var ch *Character
				if L.Get(2) != lua.LNil {
					cud, ok := L.Get(2).(*lua.LUserData)
					if !ok {
						L.ArgError(2, "controller must be Character or nil")
						return 0
					}
					var isChar bool
					ch, isChar = cud.Value.(*Character)
					if !isChar {
						L.ArgError(2, "controller must be Character or nil")
						return 0
					}
				}
				mapInstance := mob.GetMap()
				if mapInstance == nil {
					return 0
				}
				ct := mapInstance.GetControllerTable()
				ct.SwitchController(mob, ch)
				return 0
			}
			L.ArgError(2, "controller() requires 0 or 1 arguments")
			return 0
		},
	}
}

func (m *Mob) String() string {
	return m.LuaTypeName()
}

func (m *Mob) Type() lua.LValueType {
	return lua.LTUserData
}
