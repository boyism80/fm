-- Skill name (String.wz/Skill.img.xml): 알케미스트

return {
	on_passive = function(me, skill)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.x <= 0 then
			return
		end
	end,

	on_unpassive = function(me, skill)
		me:potion_heal_rate(0)
		me:potion_duration_rate(0)
	end
}
