-- Skill name (String.wz/Skill.img.xml): 실버 호크

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

function on_attack_3111005(me, skill, damages)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local prop = effect.prop
	if prop <= 0 then
		return
	end
	if prop > 100 then
		prop = 100
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

function on_activated_3111005(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	me:buff(skill, BuffFlag.Summon, 1)
end

function on_buff_3111005(me, skill)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	local wz = skill:wz()
	local level = skill:level()
	if level == nil or level <= 0 then
		return
	end
	me:create_summon(wz.id, level, effect.time, SummonMovementType.CircleFollow, SummonType.Normal)
end

function on_unbuff_3111005(me, skill)
	local wz = skill:wz()
	me:remove_summon(wz.id)
end
