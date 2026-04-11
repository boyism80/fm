-- Skill name (String.wz/Skill.img.xml): 미스틱 도어

function on_activated_2311002(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.time <= 0 then
		return
	end
	local wz = skill:wz()
	me:create_door(wz.id, effect.time)
end
