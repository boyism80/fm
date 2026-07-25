-- Skill name (String.wz/Skill.img.xml): 매직 가드

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
		    return
		end

		me:buff(skill, BuffFlag.MagicGuard, effect.x)
	end
}
