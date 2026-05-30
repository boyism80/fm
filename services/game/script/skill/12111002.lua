-- Skill name (String.wz/Skill.img.xml): 씰

function on_activated_12111002(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	for_each_mob_in_skill_area(me, skill, function(mob)
		local wz = mob:wz()
		if wz ~= nil and wz.boss then
			return false
		end
		mob:buff(MobBuff.Seal, 1, effect.time, skill, me)
		return true
	end)
end
