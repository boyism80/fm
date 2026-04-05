-- Skill name (String.wz/Skill.img.xml): 쉐도우 메소

local function shadow_meso_roll_cost(skill)
	local effect = get_skill_effect(skill)
	if effect == nil then
		return 0
	end
	local lv = skill:level()
	local delta = 170 + (lv - 1) * 10
	local half = math.floor(delta / 2)
	local min_m
	local max_m
	if effect.money_con > 0 then
		min_m = effect.money_con - half
		max_m = effect.money_con + half
	else
		min_m = 50 + (lv - 1) * 10
		max_m = 220 + (lv - 1) * 20
	end
	if min_m < 0 then
		min_m = 0
	end
	if max_m < min_m then
		return 0
	end
	return math.random(min_m, max_m)
end

local function apply_shadow_meso_cost(me, skill)
	if me == nil or skill == nil then
		return true
	end
	local cost = shadow_meso_roll_cost(skill)
	local have = math.floor(tonumber(me:meso()) or 0)
	if cost > have then
		cost = have
	end
	if cost > 0 then
		me:meso(have - cost)
		me:update_stats({ STAT.Meso })
	end
	return true
end

function on_activating_4111004(me, skill, params)
	return apply_shadow_meso_cost(me, skill)
end

function on_activated_4111004(me, skill, params)
end
