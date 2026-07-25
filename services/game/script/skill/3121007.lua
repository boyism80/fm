-- Skill name (String.wz/Skill.img.xml): 햄스트링

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end

		if effect.x == 0 then
			return
		end

		me:buff(skill, {[BuffFlag.Hamstring] = effect.x})
	end
}
