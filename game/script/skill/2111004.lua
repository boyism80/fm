-- Skill name (String.wz/Skill.img.xml): 씰

function on_activated_2111004(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	for_each_mob_in_skill_area(me, skill, function(mob)
		local wz = mob:wz()
		if wz ~= nil and wz.boss then
			return false
		end
		mob:buff(MobBuff.Seal, 1, duration_ms, skill, me)
		return true
	end)
end
