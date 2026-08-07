local combat = require("script/lib/combat")
local M = {}

M.apply_buff_from_effect = function(me, skill, flag, value_key)
	value_key = value_key or "x"
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local val = effect[value_key]
	if val == nil then
		return
	end
	me:buff(skill, flag, val)
end

M.apply_buff_fixed = function(me, skill, flag, value)
	value = value or 1
	me:buff(skill, flag, value)
end

M.add_hyper_body_bonus = function(me, skill)
	local effect = skill:effect()
	if effect == nil then return end
	local hp_percent = me:bonus_max_hp_ratio()
	local mp_percent = me:bonus_max_mp_ratio()

	-- Apply ratio delta without intermediate notifications.
	me:bonus_max_hp_ratio(hp_percent + effect.x, false)
	me:bonus_max_mp_ratio(mp_percent + effect.x, false)

	-- Notify once after clamping to the new max values.
	me:update_stats({
		STAT.MaxHp,
		STAT.Hp,
		STAT.MaxMp,
		STAT.Mp,
	})
end

M.remove_hyper_body_bonus = function(me, skill)
	local effect = skill:effect()
	if effect == nil then return end
	local hp_percent = me:bonus_max_hp_ratio()
	local mp_percent = me:bonus_max_mp_ratio()

	me:bonus_max_hp_ratio(hp_percent - effect.x, false)
	me:bonus_max_mp_ratio(mp_percent - effect.x, false)

	me:update_stats({
		STAT.MaxHp,
		STAT.Hp,
		STAT.MaxMp,
		STAT.Mp,
	})
end

M.apply_iron_body = function(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	if effect.pdd <= 0 then
		return
	end

	me:buff(skill, BuffFlag.WeaponDef, effect.pdd)
end

M.apply_monster_magnet = function(me, skill, params)
	if params == nil or params.magnet == nil then
		return
	end
	for _, entry in ipairs(params.magnet.mobs) do
		local mob = entry.mob
		local success = entry.success
		if mob then
			me:show_magnet(mob, success)
			mob:controller(me)
		end
	end
end

M.consume_combo_orbs = function(me, howmany)
	local current = me:buff_value(BuffFlag.Combo)
	if current == nil or current <= 1 then
		return
	end

	if howmany == nil then
		howmany = current - 1
	end
	if howmany <= 0 then
		return
	end

	local new_orbs = current - howmany
	if new_orbs < 1 then
		new_orbs = 1
	end

	me:buff_value(BuffFlag.Combo, new_orbs)
end

M.apply_hero_will = function(me)
	me:remove_debuff(DebuffFlag.Seduce)
end

M.apply_dispel = function(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.prop <= 0 then
		return
	end
	if effect.prop < 100 and math.random(1, 100) > effect.prop then
		return
	end
	me:remove_debuff(
		DebuffFlag.Curse,
		DebuffFlag.Darkness,
		DebuffFlag.Poison,
		DebuffFlag.Seal,
		DebuffFlag.Slow,
		DebuffFlag.Weaken
	)
end

M.get_heal_recovery_amount = function(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return 0
	end
	local hp_rate = effect.hp / 100.0
	if hp_rate <= 0 then
		return 0
	end

	local total_magic = me:base_int() + me:bonus_int()
	if total_magic <= 0 then
		total_magic = 1
	end

	local min_heal = math.floor(total_magic * 3.0 * hp_rate)
	local max_heal = math.floor(total_magic * 5.0 * hp_rate)
	if min_heal < 1 then
		min_heal = 1
	end
	if max_heal < min_heal then
		max_heal = min_heal
	end

	local heal = math.random(min_heal, max_heal)
	if me:has_debuff(DebuffFlag.Zombify) then
		heal = -heal
	end
	return heal
end

M.apply_mana_reflection = function(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	local val = effect.prop - 30
	if val < 0 then
		val = 0
	end

	me:buff(skill, BuffFlag.ManaReflection, val)
end

M.apply_resurrection = function(me, skill)
	combat.for_each_near_party_member(me, skill, function(ch)
		ch:stance(0)
		ch:hp(ch:max_hp())
		ch:mp(ch:max_mp())
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end, { allow_dead = true })
end

local archer_puppet_offset_x = 180

local function archer_puppet_facing_left(me)
	local s = me:stance()
	return s % 2 ~= 0
end

M.mark_archer_puppet = function(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	me:buff(skill, BuffFlag.Puppet, 1)
end

M.apply_archer_puppet_buff = function(me, skill)
	local wz = skill:wz()
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	local level = skill:level()
	if level == nil or level <= 0 then
		return
	end
	if effect.x <= 0 then
		return
	end
	local px, py = me:position()
	local ox = px + (archer_puppet_facing_left(me) and -archer_puppet_offset_x or archer_puppet_offset_x)
	local map = me:map()
	if map == nil then
		return
	end
	local grounded = map:point_below({ x = ox, y = py })
	local spawn_pos
	if grounded == nil then
		spawn_pos = { x = ox, y = py }
	else
		spawn_pos = { x = grounded.x, y = grounded.y }
	end
	local s = me:create_summon(wz:id(), level, effect.time, SummonMovementType.Stationary, SummonType.Puppet, spawn_pos)
	if s ~= nil then
		s:max_hp(effect.x, false)
		s:hp(effect.x, false)
	end
end

M.apply_archer_puppet_unbuff = function(me, skill)
	local wz = skill:wz()
	me:remove_summon(wz:id())
end

M.absorb_magic_guard = function(me, attacker, skill, damage)
	if damage == nil or damage <= 0 then
		return damage
	end

	local guard_percent = me:buff_value(BuffFlag.MagicGuard)
	if guard_percent == nil or guard_percent <= 0 then
		return damage
	end

	local current_mp = me:mp()
	if current_mp == nil or current_mp <= 0 then
		return damage
	end

	local mp_loss = math.floor(damage * (guard_percent / 100.0))
	if mp_loss < 0 then
		mp_loss = 0
	end
	if mp_loss > current_mp then
		mp_loss = current_mp
	end

	local hp_loss = damage - mp_loss
	if hp_loss < 0 then
		hp_loss = 0
	end

	if mp_loss > 0 then
		me:add_mp(-mp_loss)
	end

	return hp_loss
end

M.absorb_meso_guard = function(me, attacker, skill, damage)
	if damage == nil or damage <= 0 then
		return damage
	end

	local guard_percent = me:buff_value(BuffFlag.MesoGuard)
	if guard_percent == nil or guard_percent <= 0 then
		return damage
	end
	local current_meso = me:meso()
	if current_meso == nil or current_meso <= 0 then
		return damage
	end
	local meso_loss = math.floor(damage * (guard_percent / 100.0))
	if meso_loss < 0 then
		meso_loss = 0
	end
	if meso_loss > current_meso then
		meso_loss = current_meso
	end
	local hp_loss = damage - meso_loss
	if hp_loss < 0 then
		hp_loss = 0
	end
	if meso_loss > 0 then
		me:meso(current_meso - meso_loss)
	end
	return hp_loss
end

return M
