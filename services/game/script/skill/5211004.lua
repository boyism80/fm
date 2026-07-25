-- Skill name (String.wz/Skill.img.xml): 파이어 버너

local combat = require("script/lib/combat")

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
		local multiplier = combat.compute_poison_tick_multiplier(me, skill)
		local boost = me:skill(Skill.ElementBoost)
		if boost ~= nil and boost:level() > 0 then
			local boost_effect = boost:effect()
			if boost_effect ~= nil and boost_effect.x ~= nil and boost_effect.x > 0 then
				multiplier = multiplier * (boost_effect.x / 100.0 + 1.0)
			end
		end
		for mob, hits in pairs(damages) do
			if has_positive_hit(hits) and math.random(1, 100) <= prop then
				local value = combat.compute_poison_tick_damage(skill, mob, multiplier)
				if value > 0 then
					mob:buff(MobBuff.Poison, value, duration_ms, skill, me)
				end
			end
		end
	end
}
