-- Skill name (String.wz/Skill.img.xml): 미스틱 도어

function on_activated_2311002(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local duration_ms = effect.time or 0
	if duration_ms <= 0 then
		return
	end
	local wz = skill:wz()
	if wz == nil then
		return
	end
	me:create_door(wz.id, duration_ms)
end
