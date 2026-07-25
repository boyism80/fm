-- Skill name (String.wz/Skill.img.xml): 블라인드

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.x == 0 then
			return
		end
		me:buff(skill, {[BuffFlag.Blind] = effect.x})
	end
}
