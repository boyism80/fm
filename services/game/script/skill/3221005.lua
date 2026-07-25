-- Skill name (String.wz/Skill.img.xml): 프리져

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

return {
	on_attack = function(me, skill, damages)
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
					mob:buff(MobBuff.Freeze, 1, duration_ms, skill, me)
				end
			end
		end
	end,

	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.time <= 0 then
			return
		end
		me:buff(skill, BuffFlag.Summon, 1)
	end,

	on_buff = function(me, skill)
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
	end,

	on_unbuff = function(me, skill)
		local wz = skill:wz()
		me:remove_summon(wz.id)
	end
}
