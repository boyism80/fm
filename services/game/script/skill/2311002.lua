-- Skill name (String.wz/Skill.img.xml): 미스틱 도어

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil or effect.time <= 0 then
			return
		end
		me:buff(skill, BuffFlag.SoulArrow, 1)
	end,

	on_buff = function(me, skill)
		local effect = skill:effect()
		if effect == nil or effect.time <= 0 then
			return
		end
		me:create_door(skill)
	end,

	on_unbuff = function(me, skill)
		me:remove_door(skill)
	end
}
