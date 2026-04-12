-- Skill name (String.wz/Skill.img.xml): 미스틱 도어

function on_activated_2311002(me, skill, params)
	local effect = skill:effect()
	if effect == nil or effect.time <= 0 then
		return
	end
	me:buff(skill, BuffFlag.SoulArrow, 1)
end

function on_buff_2311002(me, skill)
	local effect = skill:effect()
	if effect == nil or effect.time <= 0 then
		return
	end
	local wz = skill:wz()
	me:create_door(wz.id)
end

function on_unbuff_2311002(me, skill)
	local wz = skill:wz()
	me:remove_door(wz.id)
end
