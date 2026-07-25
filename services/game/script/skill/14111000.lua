-- Skill name (String.wz/Skill.img.xml): 쉐도우 파트너

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
		    return
		end
		me:buff(skill, BuffFlag.ShadowPartner, effect.x)
	end
}
