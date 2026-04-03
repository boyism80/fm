-- Shared skill logic. Loaded once per root Lua state.
-- Individual skill scripts call apply_*(me, skill) etc. as global functions.
-- Skill (Endure, ImprovingHpRecovery, etc.) is injected by Go.

-- Applies a buff using value from effect[value_key] (e.g. "x" or "prop"). Default 0 if missing.
function apply_buff_from_effect(me, skill, flag, value_key)
	value_key = value_key or "x"
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end
	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end
	local val = effect[value_key]
	if val == nil then
		val = 0
	end
	if val <= 0 then
		return
	end
	me:buff(skill, flag, val)
end

-- Buff with fixed value (e.g. 1 for Darksight, Soul Arrow, WkCharge, Combo).
function apply_buff_fixed(me, skill, flag, value)
	value = value or 1
	me:buff(skill, flag, value)
end

function apply_sharp_eyes(me, skill)
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local x = math.floor(tonumber(effect.x) or 0)
	local y = math.floor(tonumber(effect.y) or 0)
	if x <= 0 or y <= 0 then
		return
	end
	local packed = x * 256 + (y % 256)
	for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or not ch:is_alive() then
			return
		end
		ch:buff(skill, BuffFlag.SharpEyes, packed)
	end)
end

function add_holy_symbol_bonus(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local bonus = effect.x or 0
	local current = me:bonus_exp_rate()
	if current <= 0 then
		current = 100
	end
	me:bonus_exp_rate(current + bonus)
end

function remove_holy_symbol_bonus(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local bonus = effect.x or 0
	local current = me:bonus_exp_rate()
	me:bonus_exp_rate(current - bonus)
end

function apply_hyper_body(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 0
	me:buff(skill, {
		[BuffFlag.MaxHp] = percent,
		[BuffFlag.MaxMp] = percent,
	})
end

function add_hyper_body_bonus(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 0
	local hp_percent = me:bonus_max_hp_ratio()
	local mp_percent = me:bonus_max_mp_ratio()

	-- Apply ratio delta without intermediate notifications.
	me:bonus_max_hp_ratio(hp_percent + percent, false)
	me:bonus_max_mp_ratio(mp_percent + percent, false)

	-- Notify once after clamping to the new max values.
	me:update_stats({
		STAT.MaxHp,
		STAT.Hp,
		STAT.MaxMp,
		STAT.Mp,
	})
end

function remove_hyper_body_bonus(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 0
	local hp_percent = me:bonus_max_hp_ratio()
	local mp_percent = me:bonus_max_mp_ratio()

	me:bonus_max_hp_ratio(hp_percent - percent, false)
	me:bonus_max_mp_ratio(mp_percent - percent, false)

	me:update_stats({
		STAT.MaxHp,
		STAT.Hp,
		STAT.MaxMp,
		STAT.Mp,
	})
end

function apply_iron_body(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end

	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end

	local pdd = effect.pdd or 0
	if pdd <= 0 then
		return
	end

	me:buff(skill, BuffFlag.WeaponDef, pdd)
end

function apply_monster_magnet(me, skill, params)
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

-- Consumes combo orbs from the ComboAttack buff.
-- If howmany is nil, consumes all orbs except leaves at least 1 (matching old server behavior).
function consume_combo_orbs(me, howmany)
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

function apply_hero_will(me)
	me:remove_debuff(DebuffFlag.Seduce)
end

function apply_dispel(me, skill)
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local prop = tonumber(effect.prop) or 0
	if prop <= 0 then
		return
	end
	if prop < 100 and math.random(1, 100) > prop then
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

function get_heal_recovery_amount(me, skill)
	local effect = get_skill_effect(skill)
	if effect == nil then
		return 0
	end
	local hp_rate = (effect.hp or 0) / 100.0
	if hp_rate <= 0 then
		return 0
	end

	local total_magic = (me:base_int() or 0) + (me:bonus_int() or 0)
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

function apply_mana_reflection(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then
		return
	end

	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end

	local prop = effect.prop
	if prop == nil then
		prop = 0
	end

	local val = prop - 30
	if val < 0 then
		val = 0
	end

	me:buff(skill, BuffFlag.ManaReflection, val)
end

function apply_resurrection_in_range(me, skill)
	if me == nil or skill == nil then
		return
	end
	for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or ch:is_alive() then
			return
		end
		ch:stance(0)
		ch:hp(ch:max_hp())
		ch:mp(ch:max_mp())
	end)
end
