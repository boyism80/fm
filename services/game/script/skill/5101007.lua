-- Skill name (String.wz/Skill.img.xml): 오크통

return {
	on_activated = function(me, skill, params)
		me:buff(skill, BuffFlag.Morph, Morph.Barrel)
	end
}
