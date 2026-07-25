-- Skill name (String.wz/Skill.img.xml): 픽파킷

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
		    return
		end

		me:buff(skill, BuffFlag.Pickpocket, effect.x)
	end
}
