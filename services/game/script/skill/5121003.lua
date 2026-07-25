-- Skill name (String.wz/Skill.img.xml): 슈퍼트랜스폼

return {
	on_activated = function(me, skill, params)
		local morph = Morph.SuperTransform
		if me:gender() == Gender.Female then
			morph = Morph.SuperTransformFemale
		end
		me:buff(skill, BuffFlag.Morph, morph)
	end
}
