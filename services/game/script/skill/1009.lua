-- Skill name (String.wz/Skill.img.xml): 죽간천격

return {
	dojo = true,

	on_activated = function(me, skill, params)
		me:dojo_energy(0)
	end
}
