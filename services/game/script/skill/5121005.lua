-- Skill name (String.wz/Skill.img.xml): 스내치

return {
	on_activating = function(me, skill, params)
		local morph = me:buff_value(BuffFlag.Morph)
		if morph == Morph.SuperTransform or morph == Morph.SuperTransformFemale then
			return true
		end
		return false
	end
}
