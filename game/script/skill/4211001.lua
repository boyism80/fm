-- Skill name (String.wz/Skill.img.xml): 챠크라

local function make_heal_hp(rate, stat, lower_f, upper_f)
	local low = math.floor(stat * lower_f * rate)
	local hi = math.floor(stat * upper_f * rate)
	local span = hi - low + 1
	if span < 1 then
		span = 1
	end
	return math.floor(math.random() * span) + low
end

function on_activated_4211001(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local y = effect.y
	if y == nil then
		y = 0
	end
	local luk = me:base_luk() + me:bonus_luk()
	local dex = me:base_dex() + me:bonus_dex()
	local v42 = y + 100
	local v38 = math.random(1, 100) + 100
	local hpchange = math.floor((v38 * luk * 0.033 + dex) * v42 * 0.002)
	hpchange = hpchange + make_heal_hp(y / 100.0, luk, 2.3, 3.5)
	if hpchange <= 0 then
		return
	end
	me:add_hp(hpchange)
	me:update_stats({ STAT.Hp })
end
