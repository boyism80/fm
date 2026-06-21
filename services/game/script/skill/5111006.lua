-- Skill name (String.wz/Skill.img.xml): 쇼크웨이브

function on_activating_5111006(me, skill, params)
	local morph = me:buff_value(BuffFlag.Morph)
	if morph == Morph.Transform or morph == Morph.TransformFemale then
		return true
	end
	if morph == Morph.SuperTransform or morph == Morph.SuperTransformFemale then
		return true
	end
	return false
end
