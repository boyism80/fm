-- Skill name (String.wz/Skill.img.xml): 지화천폭

return {
	dojo = true,

	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		me:dojo_energy(0)
		me:buff(skill, { [BuffFlag.BerserkFury] = effect.x })
	end
}
