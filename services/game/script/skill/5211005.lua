-- Skill name (String.wz/Skill.img.xml): 쿨링 이펙트

local function has_positive_hit(hits)
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
		if skill == nil or damages == nil then
			return
		end
		local effect = skill:effect()
		if effect == nil then
			return
		end
		local duration_ms = effect.time or 0
		if duration_ms <= 0 then
			return
		end
		local prop = effect.prop or 0
		if prop <= 0 then
			prop = 100
		end
		local boost = me:skill(Skill.ElementBoost)
		if boost ~= nil and boost:level() > 0 then
			local boost_effect = boost:effect()
			if boost_effect ~= nil and boost_effect.y ~= nil and boost_effect.y > 0 then
				duration_ms = duration_ms + boost_effect.y * 1000
			end
		end
		for mob, hits in pairs(damages) do
			if has_positive_hit(hits) and math.random(1, 100) <= prop then
				mob:buff(MobBuff.Freeze, 1, duration_ms, skill, me)
			end
		end
	end
}
