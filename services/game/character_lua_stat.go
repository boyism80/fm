package game

import (
	"log"
	"time"

	"github.com/asynkron/protoactor-go/actor"
	"github.com/boyism80/fm/core/clock"
	"github.com/boyism80/fm/core/luax"
	"github.com/boyism80/fm/maple/constant"
	internal "github.com/boyism80/fm/protocol/protobuf/gengo/fminternal"
	"github.com/boyism80/fm/protocol/response"
	lua "github.com/yuin/gopher-lua"
)

func (ch *Character) luaStatFuncs() map[string]lua.LGFunction {
	return map[string]lua.LGFunction{
		"__eq": func(L *lua.LState) int {
			leftUD := L.CheckUserData(1)
			rightUD := L.CheckUserData(2)
			leftCh, leftOk := leftUD.Value.(*Character)
			rightCh, rightOk := rightUD.Value.(*Character)
			if !leftOk || !rightOk || leftCh == nil || rightCh == nil {
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(leftCh.GetID() == rightCh.GetID()))
			return 1
		},
		"id": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetID()))
				return 1
			default:
				L.ArgError(2, "id() is read-only")
				return 0
			}
		},
		"name": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LString(ch.name))
				return 1
			case 2:
				name := L.CheckString(2)
				ch.name = name
				return 0
			default:
				L.ArgError(2, "name() requires 0 or 1 arguments")
				return 0
			}
		},
		"gender": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetGender()))
				return 1
			}
			ch.SetGender(uint8(L.CheckInt(2)))
			return 0
		},
		"hair": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetHair()))
				return 1
			}
			ch.SetHair(uint32(L.CheckInt(2)))
			return 0
		},
		"face": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetFace()))
				return 1
			}
			ch.SetFace(uint32(L.CheckInt(2)))
			return 0
		},
		"skin": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetSkinColor()))
				return 1
			}
			ch.SetSkin(uint8(L.CheckInt(2)))
			return 0
		},
		"level": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.level))
				return 1
			case 2:
				level := L.CheckInt(2)
				if level < 1 {
					level = 1
				}
				if level > 200 {
					level = 200
				}
				oldLevel := ch.level
				ch.SetLevel(uint8(level))
				if ch.level > oldLevel && ch.quests != nil {
					ch.quests.RunAutoTriggers(nil, AutoQuestTriggerLevelUp, 0)
				}
				return 0
			default:
				L.ArgError(2, "level() requires 0 or 1 arguments")
				return 0
			}
		},
		"ability_point": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.points.ap))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.points.SetAP(uint16(v), notify)
				return 0
			default:
				L.ArgError(2, "ability_point() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"skill_point": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.points.sp))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.points.SetSP(uint16(v), notify)
				return 0
			default:
				L.ArgError(2, "skill_point() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"base_hp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.baseHp))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.SetBaseHp(uint32(v), notify)
				return 0
			default:
				L.ArgError(2, "base_hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"base_mp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.baseMp))
				return 1
			case 2, 3:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				notify := argc == 2 || L.ToBool(3)
				ch.SetBaseMp(uint32(v), notify)
				return 0
			default:
				L.ArgError(2, "base_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"update_stats": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			tbl := L.CheckTable(2)
			stats := make(map[constant.Stat]int32)
			tbl.ForEach(func(k lua.LValue, v lua.LValue) {
				if v.Type() != lua.LTNumber {
					return
				}
				stat := constant.Stat(lua.LVAsNumber(v))
				if val, ok := ch.GetStatValue(stat); ok {
					stats[stat] = val
				}
			})
			if len(stats) > 0 {
				ch.listener.OnUpdateStats(ch, stats, true)
			}
			return 0
		},
		"exp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.exp))
				return 1
			case 2:
				value := float64(L.CheckNumber(2))
				if value < 0 {
					L.ArgError(2, "exp() setter requires a non-negative value; compute new total in script")
					return 0
				}
				maxU := float64(^uint32(0))
				var newExp uint32
				if value > maxU {
					newExp = ^uint32(0)
				} else {
					newExp = uint32(value)
				}
				cur := ch.exp
				if newExp == cur {
					return 0
				}
				if newExp > cur {
					ch.AddExp(newExp - cur)
					return 0
				}
				ch.exp = newExp
				ch.expReachedAt = clock.Now()
				ch.listener.OnUpdateStats(ch, map[constant.Stat]int32{
					constant.StatEXP: int32(ch.exp),
				}, false)
				return 0
			default:
				L.ArgError(2, "exp() getter: exp(); setter: exp(value)")
				return 0
			}
		},
		"meso": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.inventory.meso))
				return 1
			case 2:
				value := L.CheckNumber(2)
				if value < 0 {
					L.ArgError(2, "meso() setter requires a non-negative value; compute new total in script")
					return 0
				}
				if value <= 2147483647 {
					ch.inventory.SetMeso(int32(value))
				} else {
					ch.inventory.SetMeso(2147483647)
				}
				ch.listener.OnMesoChanged(ch, ch.inventory.meso)
				return 0
			default:
				L.ArgError(2, "meso() getter: meso(); setter: meso(value)")
				return 0
			}
		},
		"population": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)

			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.population))
				return 1
			case 2:
				ch.stats.setPopulation(int32(L.CheckInt(2)))
				return 0
			default:
				L.ArgError(2, "population() getter: population(); setter: population(value)")
				return 0
			}
		},
		"clear_population_history": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			ch.stats.populationHistory = PopulationHistory{given: map[uint32]time.Time{}}
			return 0
		},
		"reset_ap": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.points.ResetAP(constant.StatType(L.CheckInt(2)), constant.StatType(L.CheckInt(3)))))
			return 1
		},
		"reset_sp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.skills.MoveSP(uint32(L.CheckInt(2)), uint32(L.CheckInt(3)), L.CheckInt(4))))
			return 1
		},
		"forbidden": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.Forbidden(L.CheckString(2))))
			return 1
		},
		"role": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.role))
				return 1
			case 2:
				ch.role = constant.CharacterRole(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "role() requires 0 or 1 arguments")
				return 0
			}
		},
		"account_role": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			role := constant.CharacterRole(L.CheckInt(2))
			cfg, ok := luax.GetConfiguration(L)
			if !ok || cfg.ActorContext == nil {
				L.Push(lua.LFalse)
				return 1
			}
			return LuaYieldPromise(L, ch.gameWorld, ch.gameWorld.Systems().Get[*InternalSystem]().SetAccountRoleAsync(cfg.ActorContext, ch, role), func(_ *internal.SetAccountRoleReply, err error) []lua.LValue {
				if err != nil {
					return []lua.LValue{lua.LFalse}
				}
				ch.role = role
				return []lua.LValue{lua.LTrue}
			})
		},
		"class": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.class))
				return 1
			case 2, 3:
				classValue := L.CheckInt(2)
				ch.ChangeClass(constant.ClassType(classValue), argc == 2 || L.ToBool(3))
				return 0
			default:
				L.ArgError(2, "class() requires 0, 1, or 2 arguments")
				return 0
			}
		},
		"chair": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.chair))
			return 1
		},
		"stance": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stance))
				return 1
			case 2:
				ch.stance = uint8(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "stance() requires 0 or 1 argument")
				return 0
			}
		},
		"stance_of": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			group := constant.StanceKind(L.CheckInt(2))
			members, ok := constant.StanceGroupByKind[group]
			if ok {
				for _, v := range members {
					if ch.stance == v {
						L.Push(lua.LTrue)
						return 1
					}
				}
				L.Push(lua.LFalse)
				return 1
			}
			L.Push(lua.LBool(ch.stance == uint8(group)))
			return 1
		},
		"class_of": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			classCode := constant.ClassType(L.CheckInt(2))
			L.Push(lua.LBool(ch.class.Of(classCode)))
			return 1
		},
		"homing": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() != 1 {
				L.ArgError(2, "homing() takes no arguments")
				return 0
			}
			if ch.skills.homingTargetOID == nil {
				L.Push(lua.LNil)
				return 1
			}
			mapInstance := ch.GetMap()
			if mapInstance == nil {
				L.Push(lua.LNil)
				return 1
			}
			mob := mapInstance.GetMob(*ch.skills.homingTargetOID)
			if mob == nil {
				L.Push(lua.LNil)
				return 1
			}
			L.Push(luax.NewLuable(L, mob))
			return 1
		},
		"max_hp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetMaxHp()))
				return 1
			}
			ch.SetBaseHp(uint32(L.CheckInt(2)), true)
			return 0
		},
		"max_mp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if L.GetTop() == 1 {
				L.Push(lua.LNumber(ch.GetMaxMp()))
				return 1
			}
			ch.SetBaseMp(uint32(L.CheckInt(2)), true)
			return 0
		},
		"hidden": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LBool(ch.IsHidden()))
				return 1
			case 2:
				hidden := L.CheckBool(2)
				ch.SetHidden(hidden)
				return 0
			default:
				L.ArgError(2, "hidden() requires 0 or 1 arguments")
				return 0
			}
		},
		"in_dialog": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LBool(ch.dialog.Thread() != nil))
			return 1
		},
		"base_str": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.base.str))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.str = uint16(v)
				ch.notifyStatChange(constant.StatStr)
				return 0
			default:
				L.ArgError(2, "base_str() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_dex": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.base.dex))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.dex = uint16(v)
				ch.notifyStatChange(constant.StatDex)
				return 0
			default:
				L.ArgError(2, "base_dex() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_int": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.base.int))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.int = uint16(v)
				ch.notifyStatChange(constant.StatInt)
				return 0
			default:
				L.ArgError(2, "base_int() requires 0 or 1 arguments")
				return 0
			}
		},
		"base_luk": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.base.luk))
				return 1
			case 2:
				v := L.CheckInt(2)
				if v < 0 {
					v = 0
				}
				if v > int(constant.StatMaxStrDexIntLuk) {
					v = int(constant.StatMaxStrDexIntLuk)
				}
				ch.stats.base.luk = uint16(v)
				ch.notifyStatChange(constant.StatLuk)
				return 0
			default:
				L.ArgError(2, "base_luk() requires 0 or 1 arguments")
				return 0
			}
		},
		"total_str": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.totalStats().STR))
			return 1
		},
		"total_dex": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.totalStats().DEX))
			return 1
		},
		"total_int": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.totalStats().INT))
			return 1
		},
		"total_luk": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			L.Push(lua.LNumber(ch.totalStats().LUK))
			return 1
		},
		"total_magic": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			total := ch.totalStats()
			L.Push(lua.LNumber(total.MAD + total.INT + ch.buffs.value(constant.BuffFlagMagicAtk)))
			return 1
		},
		"bonus_watk": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.watk))
				return 1
			case 2:
				ch.stats.bonus.watk = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_watk() requires 0 or 1 arguments")
				return 0
			}
		},
		"bonus_max_hp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetBonusHp()))
				L.Push(lua.LNumber(ch.stats.bonus.maxHpPercent))
				return 2
			case 2:
				ch.SetBonusHp(int32(L.CheckInt(2)), true)
				return 0
			case 3:
				ch.SetBonusHp(int32(L.CheckInt(2)), false)
				ch.SetMaxHpPercent(int16(L.CheckInt(3)), true)
				return 0
			default:
				L.ArgError(2, "bonus_max_hp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_mp": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetBonusMp()))
				L.Push(lua.LNumber(ch.stats.bonus.maxMpPercent))
				return 2
			case 2:
				ch.SetBonusMp(int32(L.CheckInt(2)), true)
				return 0
			case 3:
				ch.SetBonusMp(int32(L.CheckInt(2)), false)
				ch.SetMaxMpPercent(int16(L.CheckInt(3)), true)
				return 0
			default:
				L.ArgError(2, "bonus_max_mp() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_hp_fixed": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetBonusHp()))
				return 1
			case 2, 3:
				v := int32(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetBonusHp(v, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_hp_fixed() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_hp_ratio": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.maxHpPercent))
				return 1
			case 2, 3:
				p := int16(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetMaxHpPercent(p, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_hp_ratio() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_mp_fixed": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.GetBonusMp()))
				return 1
			case 2, 3:
				v := int32(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetBonusMp(v, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_mp_fixed() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_max_mp_ratio": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.maxMpPercent))
				return 1
			case 2, 3:
				p := int16(L.CheckInt(2))
				notify := argc == 2 || L.ToBool(3)
				ch.SetMaxMpPercent(p, notify)
				return 0
			default:
				L.ArgError(2, "bonus_max_mp_ratio() requires 0, 1 or 2 arguments")
				return 0
			}
		},
		"bonus_meso_multiplier": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.mesoMultiplier))
				return 1
			case 2:
				ch.stats.bonus.mesoMultiplier = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_meso_multiplier() requires 1 or 2 arguments")
				return 0
			}
		},
		"bonus_drop_rate": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.dropRate))
				return 1
			case 2:
				ch.stats.bonus.dropRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_drop_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"channel_drop_rate": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			if ch.gameWorld == nil {
				L.Push(lua.LNumber(1))
				return 1
			}
			L.Push(lua.LNumber(ch.gameWorld.GetDropRate()))
			return 1
		},
		"bonus_exp_rate": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.expRate))
				return 1
			case 2:
				ch.stats.bonus.expRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "bonus_exp_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"potion_heal_rate": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.potionHealRate))
				return 1
			case 2:
				ch.stats.bonus.potionHealRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "potion_heal_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"potion_duration_rate": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				L.Push(lua.LNumber(ch.stats.bonus.potionDurationRate))
				return 1
			case 2:
				ch.stats.bonus.potionDurationRate = int16(L.CheckInt(2))
				return 0
			default:
				L.ArgError(2, "potion_duration_rate() requires 1 or 2 arguments")
				return 0
			}
		},
		"map": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			argc := L.GetTop()
			switch argc {
			case 1:
				m := ch.GetMap()
				if m == nil {
					L.Push(lua.LNil)
					return 1
				}
				L.Push(luax.NewLuable(L, m))
				return 1
			case 2, 3, 4:
				var targetMap *Map
				arg2 := L.Get(2)
				switch v := arg2.(type) {
				case *lua.LUserData:
					var ok bool
					targetMap, ok = v.Value.(*Map)
					if !ok || targetMap == nil {
						L.ArgError(2, "Map, map name (string), or map id (number) expected")
						return 0
					}
				case lua.LString:
					if ch.gameWorld == nil {
						L.RaiseError("map: no context to resolve map name")
						return 0
					}
					resources := ch.gameWorld.GetResources()
					if resources == nil {
						L.RaiseError("map: no resources to resolve map name")
						return 0
					}
					mapId, ok := resources.NameToMap(string(v))
					if !ok {
						L.RaiseError("map: unknown map name %q", string(v))
						return 0
					}
					targetMap = ch.gameWorld.Systems().Get[*MapSystem]().Find(ch.stateMachine, mapId)
					if targetMap == nil {
						L.RaiseError("map: map %q (id %d) not found", string(v), mapId)
						return 0
					}
				case lua.LNumber:
					if ch.gameWorld == nil {
						L.RaiseError("map: no context to resolve map id")
						return 0
					}
					mapId := uint32(v)
					targetMap = ch.gameWorld.Systems().Get[*MapSystem]().Find(ch.stateMachine, mapId)
					if targetMap == nil {
						L.RaiseError("map: map id %d not found", mapId)
						return 0
					}
				default:
					L.ArgError(2, "Map, map name (string), or map id (number) expected")
					return 0
				}
				spawnPoint := uint8(0)
				optsIndex := 4
				switch v := L.Get(3).(type) {
				case *lua.LNilType:
				case *lua.LTable:
					optsIndex = 3
				case lua.LNumber:
					spawnPoint = uint8(v)
				case lua.LString:
					portal := targetMap.FindPortalByName(string(v))
					if portal == nil {
						L.RaiseError("map: portal %q not found in map %d", string(v), targetMap.GetMapID())
						return 0
					}
					spawnPoint = portal.wz.ID
				default:
					L.ArgError(3, "spawn point id (number), portal name (string), or options (table) expected")
					return 0
				}

				relocate := false
				var callback *lua.LFunction
				if opts, ok := L.Get(optsIndex).(*lua.LTable); ok {
					relocate = lua.LVAsBool(opts.RawGetString("relocate"))
					callback, _ = opts.RawGetString("callback").(*lua.LFunction)
				}

				currentMap := ch.GetMap()
				if relocate && currentMap != nil && currentMap.GetMapID() == targetMap.GetMapID() {
					err := ch.Relocate(spawnPoint)
					if err != nil {
						L.RaiseError("warp: %v", err)
						return 0
					}
					if callback != nil {
						L.CallByParam(lua.P{Fn: callback, NRet: 0, Protect: false}, luax.NewLuable(L, ch))
					}
					return 0
				}

				var onEnter func(actor.Context)
				if callback != nil {
					detached, err := luax.Detach(L, callback)
					if err != nil {
						L.RaiseError("map: callback %v", err)
						return 0
					}
					onEnter = func(ctx actor.Context) {
						m := ch.GetMap()
						if m != targetMap {
							return
						}
						root := m.LuaRoot(ctx)
						fn, err := detached.Attach(root)
						if err == nil {
							err = luax.Spawn(root, fn, luax.Configuration{ActorContext: ctx, ActorPID: m.LogicActorPID()}, ch)
						}
						if err != nil {
							log.Printf("map callback: %v", err)
						}
					}
				}
				cfg, _ := luax.GetConfiguration(L)
				err := ch.gameWorld.Systems().Get[*MapSystem]().Warp(cfg.ActorContext, ch, targetMap, spawnPoint, onEnter)
				if err != nil {
					L.RaiseError("warp: %v", err)
				}
				return 0
			default:
				L.ArgError(2, "map() getter: 0 args; setter: map, name (string), or id (number), optional spawnPoint or portal name, optional { relocate, callback }")
				return 0
			}
		},
		"spawn_point": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			spawnID := ch.GetSpawnPoint()
			L.Push(lua.LNumber(spawnID))
			m := ch.GetMap()
			if m == nil {
				return 1
			}
			portal := m.FindPortal(spawnID)
			if portal == nil || portal.wz == nil {
				return 1
			}
			L.Push(lua.LString(portal.wz.Name))
			L.Push(lua.LNumber(portal.wz.Position.X))
			L.Push(lua.LNumber(portal.wz.Position.Y))
			return 4
		},
		"show_magnet": func(L *lua.LState) int {
			ch := luax.Check[Character](L, 1)
			mob := luax.Check[Mob](L, 2)
			success := L.CheckBool(3)
			var successByte uint8
			if success {
				successByte = 1
			}
			m := ch.GetMap()
			if m != nil {
				ch.Broadcast(&response.ShowMagnet{MobID: mob.oid, Success: successByte}, &ObjectBroadcastOption{})
			}
			return 0
		},
	}
}
