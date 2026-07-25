-- Skill name (String.wz/Skill.img.xml): 트랜스폼

return {
	on_activated = function(me, skill, params)
		local morph = Morph.Transform
		if me:gender() == Gender.Female then
			morph = Morph.TransformFemale
		end
		me:buff(skill, BuffFlag.Morph, morph)
	end
}
