-- Shared skill logic. Loaded once per root Lua state.
-- Individual skill scripts call Skill.apply_*(me, skill) for common buff patterns.

Skill = {}

-- Applies a buff using value from effect[value_key] (e.g. "x" or "prop"). Default 0 if missing.
function Skill.apply_buff_from_effect(me, skill, flag, value_key)
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
function Skill.apply_buff_fixed(me, skill, flag, value)
	value = value or 1
	me:buff(skill, flag, value)
end

function Skill.apply_booster(me, skill)
	Skill.apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function Skill.apply_stance(me, skill)
	Skill.apply_buff_from_effect(me, skill, BuffFlag.Stance, "prop")
end

function Skill.apply_invincible(me, skill)
	Skill.apply_buff_from_effect(me, skill, BuffFlag.Invincible, "x")
end

function Skill.apply_maple_warrior(me, skill)
	Skill.apply_buff_from_effect(me, skill, BuffFlag.MapleWarrior, "x")
end

function Skill.apply_powerguard(me, skill)
	Skill.apply_buff_from_effect(me, skill, BuffFlag.Powerguard, "x")
end

function Skill.apply_darksight(me, skill)
	Skill.apply_buff_fixed(me, skill, BuffFlag.Darksight, 1)
end

function Skill.apply_soul_arrow(me, skill)
	Skill.apply_buff_fixed(me, skill, BuffFlag.SoulArrow, 1)
end

function Skill.apply_wk_charge(me, skill)
	Skill.apply_buff_fixed(me, skill, BuffFlag.WkCharge, 1)
end

function Skill.apply_combo(me, skill)
	Skill.apply_buff_fixed(me, skill, BuffFlag.Combo, 1)
end

function Skill.apply_holy_symbol(me, skill)
	Skill.apply_buff_from_effect(me, skill, BuffFlag.HolySymbol, "x")
end

function Skill.add_holy_symbol_bonus(me, skill)
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

function Skill.remove_holy_symbol_bonus(me, skill)
	local wz = skill:wz()
	if wz == nil or wz.effects == nil then return end
	local effect = wz.effects[skill:level()]
	if effect == nil then return end
	local bonus = effect.x or 0
	local current = me:bonus_exp_rate()
	me:bonus_exp_rate(current - bonus)
end

function Skill.apply_hyper_body(me, skill)
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

function Skill.add_hyper_body_bonus(me, skill)
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

function Skill.remove_hyper_body_bonus(me, skill)
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
