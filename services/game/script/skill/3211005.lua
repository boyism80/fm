-- Skill name (String.wz/Skill.img.xml): 골든 이글

function on_activated_3211005(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end

	if effect.time <= 0 then
		return
	end

	me:buff(skill, BuffFlag.Summon, 1)
end

function on_buff_3211005(me, skill)
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

	local wz = skill:wz()
	me:create_summon(wz.id, level, effect.time, SummonMovementType.CircleFollow, SummonType.Normal)
end

function on_unbuff_3211005(me, skill)
	local wz = skill:wz()
	me:remove_summon(wz.id)
end

local function summon_attack_hits_positive(hits)
	if hits == nil then
		return false
	end
	for i = 1, #hits do
		if (hits[i] or 0) > 0 then
			return true
		end
	end
	return false
end

function on_attack_3211005(me, skill, damages)
	local effect = skill:effect()
	local prop = 100
	if effect ~= nil and effect.prop > 0 then
		prop = effect.prop
	end
	if prop <= 0 then
		return
	end

	local duration_ms = 4000
	for mob, hits in pairs(damages) do
		if summon_attack_hits_positive(hits) then
			if math.random(1, 100) <= prop then
				mob:buff(MobBuff.Stun, 1, duration_ms, skill, me)
			end
		end
	end
end
