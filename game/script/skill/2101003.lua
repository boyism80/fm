-- Skill name (String.wz/Skill.img.xml): 슬로우

function on_activated_2101003(me, skill, params)
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	local speed_x = effect.x
	if speed_x == nil then
		return
	end
	for_each_mob_in_skill_area(me, skill, function(mob)
		mob:set_status(MobStatus.Speed, speed_x, duration_ms, skill, me)
	end)
end
