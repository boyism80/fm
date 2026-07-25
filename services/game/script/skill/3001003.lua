-- Skill name (String.wz/Skill.img.xml): 포커스

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		if effect.acc <= 0 and effect.eva <= 0 then
			return
		end
		me:buff(skill, {
			[BuffFlag.Acc] = effect.acc,
			[BuffFlag.Avoid] = effect.eva,
		})
	end
}
