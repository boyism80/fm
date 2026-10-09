-- Skill name (String.wz/Skill.img.xml): 금강불괴

return {
	dojo = true,

	on_activated = function(me, skill, params)
		me:dojo_energy(0)
		me:buff(skill, { [BuffFlag.DivineBody] = 1 })
	end
}
