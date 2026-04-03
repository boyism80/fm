-- Skill name (String.wz/Skill.img.xml): 실버 호크

local function summon_attack_hits_positive(hits)
	if hits == nil then
		return false
	end
	for i = 1, #hits do
		local v = hits[i]
		if v ~= nil and v > 0 then
			return true
		end
	end
	return false
end

function on_attack_3111005(me, skill, damages)
	if me == nil or skill == nil or damages == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local prop = math.floor(tonumber(effect.prop) or 0)
	if prop <= 0 then
		return
	end
	if prop > 100 then
		prop = 100
	end
	local duration_ms = 4000
	for mob, hits in pairs(damages) do
		if mob ~= nil and summon_attack_hits_positive(hits) then
			if math.random(1, 100) <= prop then
				mob:set_status(MobStatus.Stun, 1, duration_ms, skill, me)
			end
		end
	end
end

function on_activated_3111005(me, skill, params)
	if me == nil or skill == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local duration_ms = tonumber(effect.time) or 0
	if duration_ms <= 0 then
		return
	end
	me:buff(skill, BuffFlag.Summon, 1)
end

function on_buff_3111005(me, skill)
	if me == nil or skill == nil then
		return
	end
	local wz = skill:wz()
	if wz == nil or wz.id == nil then
		return
	end

	local effect = wz.effects[skill:level()]
	if effect == nil then
		return
	end
	local duration_ms = tonumber(effect.time) or 0
	if duration_ms <= 0 then
		return
	end
	local level = skill:level()
	if level == nil or level <= 0 then
		return
	end
	me:create_summon(wz.id, level, duration_ms, SummonMovementType.CircleFollow, SummonType.Normal)
end

function on_unbuff_3111005(me, skill)
	if me == nil or skill == nil then
		return
	end
	local wz = skill:wz()
	if wz == nil or wz.id == nil then
		return
	end
	me:remove_summon(wz.id)
end
