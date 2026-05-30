-- Skill name (String.wz/Skill.img.xml): 슬로우

function on_activated_2201003(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	for_each_mob_in_skill_area(me, skill, function(mob)
		mob:buff(MobBuff.Speed, effect.x, effect.time, skill, me)
	end)
end
