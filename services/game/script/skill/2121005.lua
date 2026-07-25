-- Skill name (String.wz/Skill.img.xml): 엘퀴네스

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.time <= 0 then
			return
		end
		me:buff(skill, BuffFlag.Summon, 1)
	end,

	on_buff = function(me, skill)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.time <= 0 then
			return
		end
		local wz = skill:wz()
		local level = skill:level()
		if level == nil or level <= 0 then
			return
		end
		me:create_summon(wz.id, level, effect.time, SummonMovementType.Follow, SummonType.Normal)
	end,

	on_unbuff = function(me, skill)
		local wz = skill:wz()
		me:remove_summon(wz.id)
	end
}
