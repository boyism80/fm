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
	me:buff(skill, flag, val)
end

-- Buff with fixed value (e.g. 1 for Darksight, Soul Arrow, WkCharge, Combo).
function apply_buff_fixed(me, skill, flag, value)
	value = value or 1
	me:buff(skill, flag, value)
end

function apply_booster(me, skill)
	apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function apply_stance(me, skill)
	apply_buff_from_effect(me, skill, BuffFlag.Stance, "prop")
end

function apply_invincible(me, skill)
	apply_buff_from_effect(me, skill, BuffFlag.Invincible, "x")
end

function apply_maple_warrior(me, skill)
	apply_buff_from_effect(me, skill, BuffFlag.MapleWarrior, "x")
end

function apply_powerguard(me, skill)
	apply_buff_from_effect(me, skill, BuffFlag.Powerguard, "x")
end

function apply_darksight(me, skill)
	apply_buff_fixed(me, skill, BuffFlag.Darksight, 1)
end

function apply_soul_arrow(me, skill)
	apply_buff_fixed(me, skill, BuffFlag.SoulArrow, 1)
end

function apply_wk_charge(me, skill)
	apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

function apply_combo(me, skill)
	apply_buff_fixed(me, skill, BuffFlag.Combo, 1)
end

function apply_holy_symbol(me, skill)
	apply_buff_from_effect(me, skill, BuffFlag.HolySymbol, "x")
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
	local hp_fixed, hp_percent = me:bonus_max_hp()
	local mp_fixed, mp_percent = me:bonus_max_mp()
	me:bonus_max_hp(hp_fixed, hp_percent + percent)
	me:bonus_max_mp(mp_fixed, mp_percent + percent)
end

function remove_hyper_body_bonus(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local percent = effect.x or 0
	local hp_fixed, hp_percent = me:bonus_max_hp()
	local mp_fixed, mp_percent = me:bonus_max_mp()
	me:bonus_max_hp(hp_fixed, hp_percent - percent)
	me:bonus_max_mp(mp_fixed, mp_percent - percent)
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
	local magnet_direction = params.magnet.direction or 3
	me:show_buffeffect(skill, 1, magnet_direction)
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
